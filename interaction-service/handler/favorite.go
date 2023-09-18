package handler

import (
	"context"
	vp "demotest/douyin-api/proto/video"
	"demotest/douyin-api/util"
	"demotest/interaction-service/dao"
	"demotest/interaction-service/global"
	"demotest/douyin-api/globalinit/constant"
	"demotest/interaction-service/model"
	"demotest/interaction-service/proto/favorite"
	"errors"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
	"gorm.io/gorm"
	"sync"
	"fmt"
)

type FavoriteService struct{}
/*
喜欢列表(❤ ω ❤)
 */

func (s *FavoriteService) FavoriteList(ctx context.Context, request *favorite.FavoriteListRequest) (resp *favorite.FavoriteListResponse, err error) {

	global.InfoLogger.Info(constant.FavoriteServiceName,
		zap.String("method", "FavoriteList"),
	)

	favorites, err := dao.Q.Favorite.
		Where(dao.Favorite.UserId.Eq(request.UserId)).
		Select(dao.Q.Favorite.VideoId).
		Find()

	if err != nil {
		global.InfoLogger.Error(constant.FavoriteServiceName,
			zap.Int64("userId", request.UserId),
			zap.Error(err),
		)

		resp = &favorite.FavoriteListResponse{
			StatusCode: constant.TotalFavoriteErrCode,
			StatusMsg:  proto.String(constant.ErrorMsg),
		}
		err = nil
		return
	}
	vids := getVids(favorites)        //点赞视频ids
	videos := GetVideos(vids)         //所有点赞视频信息
	videoMap := sync.Map{}            //点赞视频信息映射
	authorIds := getAuthorIds(videos) //作者ids
	authors := GetUsers(authorIds)    //作者信息
	authorsMap := sync.Map{}          //作者信息映射
	// 设置最大并发数
	maxConcurrency := constant.MaxConcurrency
	concurrencyCh := make(chan struct{}, maxConcurrency)
	var wg sync.WaitGroup
	for _, v := range videos {
		v := v
		concurrencyCh <- struct{}{} // 占用一个并发槽
		wg.Add(1)
		go func() {
			defer wg.Done()
			defer func() {
				<-concurrencyCh // 释放一个并发槽
			}()
			videoMap.Store(v.VideoID, v)
		}()
	}
	wg.Wait()

	for _, user := range authors {
		user := user
		concurrencyCh <- struct{}{} // 占用一个并发槽
		wg.Add(1)
		go func() {
			defer wg.Done()
			defer func() {
				<-concurrencyCh // 释放一个并发槽
			}()
			u := vp.User{}
			totalFavorited, favoriteCount, followerCount, followCount, _ := util.GetUserAllUserData(user.ID)
			u.TotalFavorited = totalFavorited
			u.FavoriteCount = favoriteCount
			u.FollowCount = followerCount
			u.FollowCount = followCount
			u.Id = user.ID
			u.Name = user.Name
			u.IsFollow = user.IsFollow
			u.Avatar = user.Avatar
			u.BackgroundImage = user.BackgroundImage
			u.Signature = user.Signature
			u.WorkCount = user.WorkCount
			authorsMap.Store(u.Id, &u)
		}()
	}
	wg.Wait()

	videoList := make([]*vp.Video, len(favorites))
	for i, f := range favorites {
		i := i
		f := f
		concurrencyCh <- struct{}{} // 占用一个并发槽
		wg.Add(1)
		go func() {
			defer wg.Done()
			defer func() {
				<-concurrencyCh // 释放一个并发槽
			}()
			if value, ok := videoMap.Load(f.VideoId); ok != false {
				video := value.(*model.Video)
				if author, ok := authorsMap.Load(video.UserID); ok != false {
					v := &vp.Video{
						Id:            video.VideoID,
						Author:        author.(*vp.User),
						PlayUrl:       video.PlayURL,
						CoverUrl:      video.CoverURL,
						FavoriteCount: video.FavoriteCount,
						CommentCount:  video.CommentCount,
						IsFavorite:    true,
						Title:         video.Title,
					}
					videoList[i] = v
				}
			}
		}()

	}
	wg.Wait()

	resp = &favorite.FavoriteListResponse{
		VideoList:  videoList,
		StatusCode: 0,
		StatusMsg:  proto.String(constant.SuccessMsg),
	}

	return
}

// IsFavorite 是否点赞视频
func (s *FavoriteService) IsFavorite(ctx context.Context, req *favorite.IsFavoriteRequest) (resp *favorite.IsFavoriteResponse, err error) {
	resp, err = isFavorite(ctx, req)
	if err != nil {
		global.InfoLogger.Error(constant.VideoNotExist,
			zap.Int64("videoId", req.VideoId),
			zap.Int64("UserId", req.UserId),
			zap.String("msg", "判断是否点赞错误"),
			zap.Error(err),
		)

		resp = &favorite.IsFavoriteResponse{
			StatusCode: constant.VideoNotExistCode,
			IsFavorite: false,
		}
		err = nil
		return
	}
	return
}

// FavoriteAction 点赞、取消点赞操作
func (s *FavoriteService) FavoriteAction(ctx context.Context, req *favorite.FavoriteActionRequest) (resp *favorite.FavoriteActionResponse, err error) {
	global.InfoLogger.Info(constant.FavoriteServiceName,
		zap.String("method", "FavoriteAction"),
	)

	//视频是否存在
	_, _, err = util.GetVideoFavoriteAndCommentCount(req.VideoId)
	if err != nil {
		global.InfoLogger.Error(constant.VideoExistErrMsg,
			zap.String("Msg", "视频不存在"),
			zap.Int64("videoId", req.VideoId),
			zap.Error(err),
		)
		resp = &favorite.FavoriteActionResponse{
			StatusCode: constant.VideoFavoriteCountErrCode,
			StatusMsg:  proto.String(constant.ErrorMsg),
		}
		err = nil
		return
	}
	//用户是否存在
	_, _, _, _, err = util.GetUserAllUserData(req.UserId)
	if err != nil {
		global.InfoLogger.Error(constant.UserExistErrMsg,
			zap.String("Msg", "用户不存在"),
			zap.Int64("userId", req.UserId),
			zap.Error(err),
		)
		resp = &favorite.FavoriteActionResponse{
			StatusCode: constant.VideoFavoriteCountErrCode,
			StatusMsg:  proto.String(constant.ErrorMsg),
		}
		err = nil
		return
	}

	if req.ActionType == 1 {
		resp, err = like(ctx, req.UserId, req.VideoId)
	} else if req.ActionType == 2 {
		resp, err = cancelLike(ctx, req.UserId, req.VideoId)
	} else {
		global.InfoLogger.Error(constant.FavoriteServiceName,
			zap.Int64("videoId", req.VideoId),
			zap.Int64("UserId", req.UserId),
			zap.Int32("ActionType", req.ActionType),
			zap.String("Msg", "ActionType非法"),
		)

		resp = &favorite.FavoriteActionResponse{
			StatusCode: constant.ParamErr,
			StatusMsg:  proto.String(constant.ErrorMsg),
		}
		err = nil
		return
	}
	return
}

func like(ctx context.Context, uid int64, vid int64) (resp *favorite.FavoriteActionResponse, err error) {
	//是否重复点赞
	_, err = dao.Q.Favorite.
		Where(dao.Favorite.UserId.Eq(uid), dao.Favorite.VideoId.Eq(vid)).First()
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		global.InfoLogger.Error(constant.FavoriteServiceName,
			zap.Int64("videoId", vid),
			zap.Int64("UserId", uid),
			zap.String("msg", "重复点赞"),
			zap.Error(err),
		)
		resp = &favorite.FavoriteActionResponse{
			StatusCode: constant.FavoriteActionLikeErrCode,
			StatusMsg:  proto.String(constant.ErrorMsg),
		}
		err = nil
		return
	}

	err = dao.Q.Favorite.Create(&model.Favorite{
		UserId:  uid,
		VideoId: vid,
	})
	err = util.LikeVideo(vid)
	err = util.IncreaseFavoriteCount(uid)
	err = util.IncreaseTotalFavorited(GetVideo(vid).Author.Id)

	if err != nil {
		global.InfoLogger.Error(constant.FavoriteServiceName,
			zap.Int64("videoId", vid),
			zap.Int64("UserId", uid),
			zap.String("Msg", "点赞失败"),
			zap.Error(err),
		)

		resp = &favorite.FavoriteActionResponse{
			StatusCode: constant.FavoriteActionLikeErrCode,
			StatusMsg:  proto.String(constant.ErrorMsg),
		}
		err = nil
		return
	}

	resp = &favorite.FavoriteActionResponse{
		StatusCode: 0,
		StatusMsg:  proto.String(constant.SuccessMsg),
	}
	return
}

func cancelLike(ctx context.Context, uid int64, vid int64) (resp *favorite.FavoriteActionResponse, err error) {
	//是否未点赞
	_, err = dao.Q.Favorite.
		Where(dao.Favorite.UserId.Eq(uid), dao.Favorite.VideoId.Eq(vid)).First()
	if err != nil {
		global.InfoLogger.Error(constant.FavoriteServiceName,
			zap.Int64("videoId", vid),
			zap.Int64("UserId", uid),
			zap.String("msg", "未点赞"),
			zap.Error(err),
		)
		resp = &favorite.FavoriteActionResponse{
			StatusCode: constant.FavoriteActionLikeErrCode,
			StatusMsg:  proto.String(constant.ErrorMsg),
		}
		err = nil
		return
	}

	_, err = dao.Q.Favorite.Unscoped().
		Where(dao.Favorite.UserId.Eq(uid), dao.Favorite.VideoId.Eq(vid)).
		Delete()
	err = util.UnlikeVideo(vid)
	err = util.DecreaseTotalFavorited(GetVideo(vid).Author.Id)
	err = util.DecreaseFavoriteCount(uid)

	if err != nil {
		global.InfoLogger.Error(constant.FavoriteServiceName,
			zap.Int64("videoId", vid),
			zap.Int64("UserId", uid),
			zap.Error(err),
		)

		resp = &favorite.FavoriteActionResponse{
			StatusCode: constant.FavoriteActionCancelLikeErrCode,
			StatusMsg:  proto.String(constant.ErrorMsg),
		}
		err = nil
		return
	}

	resp = &favorite.FavoriteActionResponse{
		StatusCode: 0,
		StatusMsg:  proto.String(constant.SuccessMsg),
	}
	return
}

func GetVideo(vid int64) *vp.Video {
	var video *model.Video
	videoRet := &vp.Video{}
	if err := global.DB.Table("videos").Where("video_id=?", vid).
		Find(&video).Error; err != nil {
		return &vp.Video{}
	}

	//获取点赞数和评论数
	if commentCount, favoriteCount, err := util.GetVideoFavoriteAndCommentCount(vid); err != nil {
		global.InfoLogger.Error(constant.FavoriteServiceName,
			zap.String("Msg", "获取视频commentCount、favoriteCount失败"),
			zap.Int64("videoId", vid),
			zap.Error(err),
		)
	} else {
		videoRet.CommentCount = commentCount
		videoRet.FavoriteCount = favoriteCount
	}

	user, err := GetUser(video.UserID)
	if err != nil {
		global.InfoLogger.Error(constant.FavoriteServiceName,
			zap.Int64("UserId", video.UserID),
			zap.Error(err),
		)
	} else {
		videoRet.Author = &vp.User{}
		videoRet.Author.WorkCount = user.WorkCount
		videoRet.Author.Name = user.Name
		videoRet.Author.Signature = user.Signature
		videoRet.Author.IsFollow = user.IsFollow
		videoRet.Author.BackgroundImage = user.BackgroundImage
		videoRet.Author.Avatar = user.Avatar
		videoRet.Author.Id = user.ID
	}

	videoRet.Id = video.VideoID
	videoRet.CoverUrl = video.CoverURL
	videoRet.Title = video.Title
	videoRet.PlayUrl = video.PlayURL
	videoRet.IsFavorite = true

	return videoRet
}

func getVids(favorites []*model.Favorite) (vids []int64) {
	idMap := sync.Map{}
	ids := make([]int64, 0)
	// 设置最大并发数
	maxConcurrency := constant.MaxConcurrency
	concurrencyCh := make(chan struct{}, maxConcurrency)
	var wg sync.WaitGroup
	for _, f := range favorites {
		vid := f.VideoId
		concurrencyCh <- struct{}{} // 占用一个并发槽
		wg.Add(1)
		go func() {
			defer wg.Done()
			defer func() {
				<-concurrencyCh // 释放一个并发槽
			}()
			if _, ok := idMap.Load(vid); ok == false {
				ids = append(ids, vid)
				idMap.Store(vid, true)
			}
		}()
	}
	wg.Wait()
	return ids
}
func GetVideos(ids []int64) []*model.Video {
	var videos []*model.Video
	if err := global.DB.Table("videos").Where("video_id In ?", ids).Find(&videos).Error; err != nil {
		global.InfoLogger.Error(constant.FavoriteServiceName,
			zap.String("Videos", fmt.Sprintf("%+v", ids)),
			zap.Error(err),
		)
		return []*model.Video{}
	}
	return videos
}
func getAuthorIds(videos []*model.Video) (vids []int64) {
	idMap := sync.Map{}
	ids := make([]int64, 0)
	// 设置最大并发数
	maxConcurrency := constant.MaxConcurrency
	concurrencyCh := make(chan struct{}, maxConcurrency)
	var wg sync.WaitGroup
	for _, v := range videos {
		uid := v.UserID
		concurrencyCh <- struct{}{} // 占用一个并发槽
		wg.Add(1)
		go func() {
			defer wg.Done()
			defer func() {
				<-concurrencyCh // 释放一个并发槽
			}()
			if _,ok:=idMap.Load(uid);ok==false {
				ids = append(ids, uid)
				idMap.Store(uid,true)
			}
		}()
	}
	wg.Wait()
	return ids
}
func GetUser(uid int64) (*model.User, error) {
	var user *model.User
	err := global.DB.Table("user").Where("id=?", uid).
		First(&user).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return &model.User{}, errors.New("用户不存在")
	} else {
		totalFavorited, favoriteCount, followerCount, followCount, _ := util.GetUserAllUserData(user.ID)
		user.TotalFavorited = totalFavorited
		user.FavoriteCount = favoriteCount
		user.FollowCount = followerCount
		user.FollowCount = followCount
		return user, nil
	}
}

func isFavorite(ctx context.Context, req *favorite.IsFavoriteRequest) (resp *favorite.IsFavoriteResponse, err error) {

	global.InfoLogger.Info(constant.FavoriteServiceName,
		zap.String("method", "IsFavorite"),
	)

	_, err = dao.Q.Favorite.
		Where(dao.Favorite.UserId.Eq(req.UserId), dao.Favorite.VideoId.Eq(req.VideoId)).First()

	resp = &favorite.IsFavoriteResponse{
		StatusCode: 0,
		IsFavorite: err == nil,
	}

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {

		global.InfoLogger.Error(constant.FavoriteServiceName,
			zap.Int64("videoId", req.VideoId),
			zap.Int64("UserId", req.UserId),
			zap.Error(err),
		)
		resp.StatusCode = constant.IsFavoriteErrCode
	}

	return
}

