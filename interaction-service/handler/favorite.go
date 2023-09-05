package handler

import (
	"context"
	vp "demotest/douyin-api/proto/video"
	"demotest/douyin-api/util"
	"demotest/interaction-service/dao"
	"demotest/interaction-service/global"
	"demotest/interaction-service/global/constant"
	"demotest/interaction-service/model"
	"demotest/interaction-service/proto/favorite"
	"errors"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/proto"
	"gorm.io/gorm"
)

var videoServiceClient vp.VideoServiceClient

func init() {
	conn, err := grpc.Dial(constant.VideoServiceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {

		panic(err)
	}
	global.ConsoleLogger.Info(constant.VideoServiceClientName,
		zap.String("Addr: ", constant.VideoServiceAddr),
	)
	global.InfoLogger.Info(constant.VideoServiceClientName,
		zap.String("Addr: ", constant.VideoServiceAddr),
	)
	videoServiceClient = vp.NewVideoServiceClient(conn)
}

type FavoriteService struct{}
/*
喜欢列表(❤ ω ❤)
 */
func (s *FavoriteService) FavoriteList(ctx context.Context, request *favorite.FavoriteListRequest) (resp *favorite.FavoriteListResponse, err error) {

	global.ConsoleLogger.Info(constant.FavoriteServiceName,
		zap.String("method", "FavoriteList"),
	)

	favorites, err := dao.Q.Favorite.
		Where(dao.Favorite.UserId.Eq(request.UserId)).
		Select(dao.Q.Favorite.VideoId).
		Find()
	if err != nil {
		global.ConsoleLogger.Error(constant.FavoriteServiceName,
			zap.Int64("userId", request.UserId),
			zap.Error(err),
		)
		global.ErrLogger.Error(constant.FavoriteServiceName,
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

	videoList := make([]*vp.Video, len(favorites))
	for i, f := range favorites {
		video := GetVideo(f.VideoId)
		favoriteCount,_, _ :=util.GetVideoFavoriteAndCommentCount(video.Id)
		video.FavoriteCount = favoriteCount
		videoList[i] = video
	}

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
		global.ConsoleLogger.Error(constant.VideoNotExist,
			zap.Int64("videoId", req.VideoId),
			zap.Int64("UserId", req.UserId),
			zap.String("msg", "判断是否点赞错误"),
			zap.Error(err),
		)
		global.ErrLogger.Error(constant.VideoNotExist,
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

	global.ConsoleLogger.Info(constant.FavoriteServiceName,
		zap.String("method", "FavoriteAction"),
	)

	//视频是否存在
	var exist bool
	exist, err = VideoExist(req.VideoId)
	if err != nil || !exist {
		global.ConsoleLogger.Error(constant.VideoExistErrMsg,
			zap.Int64("videoId", req.VideoId),
			zap.Error(err),
		)
		global.ErrLogger.Error(constant.VideoExistErrMsg,
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

	//用户不存在打印和保存日志并返回
	_, err = GetUser(req.UserId)
	if err != nil {
		global.ConsoleLogger.Error(constant.UserExistErrMsg,
			zap.Int64("userId", req.UserId),
			zap.Error(err),
		)
		global.ErrLogger.Error(constant.UserExistErrMsg,
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
	if !exist {
		global.ConsoleLogger.Error(constant.UserNotExist,
			zap.Int64("userId", req.UserId),
			zap.String("msg", "用户不存在"),
			zap.Error(err),
		)
		global.ErrLogger.Error(constant.UserNotExist,
			zap.Int64("userId", req.UserId),
			zap.String("msg", "用户不存在"),
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
		global.ConsoleLogger.Error(constant.FavoriteServiceName,
			zap.Int64("videoId", req.VideoId),
			zap.Int64("UserId", req.UserId),
			zap.Int32("ActionType", req.ActionType),
			zap.Error(errors.New("ActionType非法")),
		)
		global.ErrLogger.Error(constant.FavoriteServiceName,
			zap.Int64("videoId", req.VideoId),
			zap.Int64("UserId", req.UserId),
			zap.Int32("ActionType", req.ActionType),
			zap.Error(errors.New("ActionType非法")),
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
		global.ConsoleLogger.Error(constant.FavoriteServiceName,
			zap.Int64("videoId", vid),
			zap.Int64("UserId", uid),
			zap.String("msg", "重复点赞"),
			zap.Error(err),
		)
		global.ErrLogger.Error(constant.FavoriteServiceName,
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
		global.ConsoleLogger.Error(constant.FavoriteServiceName,
			zap.Int64("videoId", vid),
			zap.Int64("UserId", uid),
			zap.String("Msg", "点赞失败"),
			zap.Error(err),
		)
		global.ErrLogger.Error(constant.FavoriteServiceName,
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
		global.ConsoleLogger.Error(constant.FavoriteServiceName,
			zap.Int64("videoId", vid),
			zap.Int64("UserId", uid),
			zap.String("msg", "未点赞"),
			zap.Error(err),
		)
		global.ErrLogger.Error(constant.FavoriteServiceName,
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
		global.ConsoleLogger.Error(constant.FavoriteServiceName,
			zap.Int64("videoId", vid),
			zap.Int64("UserId", uid),
			zap.Error(err),
		)
		global.ErrLogger.Error(constant.FavoriteServiceName,
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

// VideoIds 根据用户id获取用户创作的所有视频的id
func VideoIds(uid int64) ([]int64, error) {

	var videos []model.Video

	if err := global.DB.Table("videos").
		Select("video_id").
		Where("user_id=?", uid).
		Find(&videos).Error; err != nil {
		return nil, err
	}
	ids := make([]int64, len(videos))
	for i, v := range videos {
		ids[i] = v.VideoID
	}
	return ids, nil
}

func VideoExist(vid int64) (bool, error) {
	var video *model.Video
	err := global.DB.Table("videos").Where("video_id=?", vid).
		First(&video).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return false, err
	} else {
		return err == nil, nil
	}
}
func GetVideo(vid int64) *vp.Video {
	var video *model.Video
	videoRet := &vp.Video{}
	if err := global.DB.Table("videos").Where("video_id=?", vid).
		Find(&video).Error; err != nil {
		return nil
	}

	//获取点赞数和评论数
	if commentCount, favoriteCount, err := util.GetVideoFavoriteAndCommentCount(vid); err != nil {
		global.ConsoleLogger.Error(constant.FavoriteServiceName,
			zap.Int64("videoId", vid),
			zap.Error(err),
		)
		global.ErrLogger.Error(constant.FavoriteServiceName,
			zap.Int64("videoId", vid),
			zap.Error(err),
		)
	} else {
		videoRet.CommentCount = commentCount
		videoRet.FavoriteCount = favoriteCount
	}
	/*if isFavoriteResponse, err := isFavorite(context.Background(), &favorite.IsFavoriteRequest{
		UserId:  video.UserID,
		VideoId: video.VideoID,
	}); err != nil {
		videoRet.IsFavorite = false
	} else {
		videoRet.IsFavorite = isFavoriteResponse.IsFavorite

	}*/

	user, err := GetUser(video.UserID)
	if err != nil {
		global.ConsoleLogger.Error(constant.FavoriteServiceName,
			zap.Int64("UserId", video.UserID),
			zap.Error(err),
		)
		global.ErrLogger.Error(constant.FavoriteServiceName,
			zap.Int64("UserId", video.UserID),
			zap.Error(err),
		)
	} else {
		videoRet.Author = &vp.User{}
		videoRet.Author.TotalFavorited = user.TotalFavorited
		videoRet.Author.FollowerCount = user.FollowCount
		videoRet.Author.WorkCount = user.WorkCount
		videoRet.Author.Name = user.Name
		videoRet.Author.FollowerCount = user.FollowerCount
		videoRet.Author.Signature = user.Signature
		videoRet.Author.IsFollow = user.IsFollow
		videoRet.Author.BackgroundImage = user.BackgroundImage
		videoRet.Author.Avatar = user.Avatar
		videoRet.Author.FavoriteCount = user.FavoriteCount
		videoRet.Author.Id = user.ID
	}

	videoRet.Id = video.VideoID
	videoRet.CoverUrl = video.CoverURL
	videoRet.Title = video.Title
	videoRet.PlayUrl = video.PlayURL
	videoRet.IsFavorite = true

	return videoRet
}
func GetUser(uid int64) (*model.User, error) {
	var user *model.User
	err := global.DB.Table("user").Where("id=?", uid).
		First(&user).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("用户不存在")
	} else {
		return user, nil
	}
}
func isFavorite(ctx context.Context, req *favorite.IsFavoriteRequest) (resp *favorite.IsFavoriteResponse, err error) {

	global.ConsoleLogger.Info(constant.FavoriteServiceName,
		zap.String("method", "IsFavorite"),
	)

	_, err = dao.Q.Favorite.
		Where(dao.Favorite.UserId.Eq(req.UserId), dao.Favorite.VideoId.Eq(req.VideoId)).First()

	resp = &favorite.IsFavoriteResponse{
		StatusCode: 0,
		IsFavorite: err == nil,
	}

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {

		global.ConsoleLogger.Error(constant.FavoriteServiceName,
			zap.Int64("videoId", req.VideoId),
			zap.Int64("UserId", req.UserId),
			zap.Error(err),
		)
		global.ErrLogger.Error(constant.FavoriteServiceName,
			zap.Int64("videoId", req.VideoId),
			zap.Int64("UserId", req.UserId),
			zap.Error(err),
		)
		resp.StatusCode = constant.IsFavoriteErrCode
	}

	return
}
