package handler

import (
	"context"
	pb "demotest/douyin-api/proto/user"
	"demotest/douyin-api/util"
	"demotest/interaction-service/dao"
	"demotest/interaction-service/global"
	"demotest/douyin-api/globalinit/constant"
	"demotest/interaction-service/model"
	"demotest/interaction-service/proto/comment"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/proto"
	"strconv"
	"sync"
	"fmt"
)

var userServiceClient pb.UserServiceClient

func init() {
	conn, err := grpc.Dial(constant.UserServiceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {

		panic(err)
	}
	global.InfoLogger.Info(constant.UserServiceClientName,
		zap.String("Addr: ", constant.UserServiceAddr),
	)
	userServiceClient = pb.NewUserServiceClient(conn)
}

type CommentService struct {
}

func (c *CommentService) CommentList(ctx context.Context, request *comment.CommentListRequest) (resp *comment.CommentListResponse, err error) {

	global.InfoLogger.Info(constant.CommentServiceName,
		zap.String("method", "CommentList"),
	)

	//视频是否存在
	_, _, err = util.GetVideoFavoriteAndCommentCount(request.VideoId)
	if err != nil {
		global.InfoLogger.Error(constant.VideoNotExist,
			zap.Int64("videoId", request.VideoId),
			zap.Error(err),
		)

		resp = &comment.CommentListResponse{
			CommentList: nil,
			StatusCode:  constant.VideoNotExistCode,
			StatusMsg:   proto.String(constant.ErrorMsg),
		}
		return
	}

	var comments []*model.Comment

	comments, err = dao.Q.Comment.
		WithContext(ctx).
		Where(dao.Q.Comment.VideoId.Eq(request.VideoId)).
		Find()

	if err != nil {

		global.InfoLogger.Error(constant.FavoriteServiceName,
			zap.Int64("vid", request.VideoId),
			zap.Error(err),
		)

		return &comment.CommentListResponse{
			StatusCode:  constant.CommentListErrCode,
			StatusMsg:   proto.String(constant.ErrorMsg),
			CommentList: nil,
		}, nil
	}
	uids := getUids(comments)
	users := GetUsers(uids)
	userMap := sync.Map{}

	// 设置最大并发数
	maxConcurrency := constant.MaxConcurrency
	concurrencyCh := make(chan struct{}, maxConcurrency)
	var wg sync.WaitGroup
	for _, u := range users {
		u := u
		concurrencyCh <- struct{}{} // 占用一个并发槽
		wg.Add(1)
		go func() {
			defer wg.Done()
			defer func() {
				<-concurrencyCh // 释放一个并发槽
			}()
			userMap.Store(u.ID, u)
		}()
	}
	wg.Wait()

	commentList := make([]*comment.Comment, len(comments))

	for i, cmt := range comments {
		i := i
		cmt := cmt
		concurrencyCh <- struct{}{} // 占用一个并发槽
		wg.Add(1)
		go func() {
			defer wg.Done()
			defer func() {
				<-concurrencyCh // 释放一个并发槽
			}()
			user := &pb.User{}
			if u,ok := userMap.Load(cmt.UserId); ok!=false {
				us:=u.(*model.User)
				user = &pb.User{
					Id:              us.ID,
					Name:            us.Name,
					FollowCount:     us.FollowCount,
					FollowerCount:   us.FollowerCount,
					IsFollow:        us.IsFollow,
					Avatar:          us.Avatar,
					BackgroundImage: us.BackgroundImage,
					Signature:       us.Signature,
					TotalFavorited:  us.TotalFavorited,
					WorkCount:       us.WorkCount,
					FavoriteCount:   us.FavoriteCount,
				}
			}
			commentList[i] = &comment.Comment{
				Id:         cmt.CommentId,
				User:       user,
				Content:    cmt.Content,
				CreateDate: cmt.CreateDate.Format("01-02"),
			}
		}()
	}
	// 等待所有协程完成
	wg.Wait()

	resp = &comment.CommentListResponse{
		StatusCode:  0,
		StatusMsg:   proto.String(constant.SuccessMsg),
		CommentList: commentList,
	}
	return
}

func (c *CommentService) CommentAction(ctx context.Context, request *comment.CommentActionRequest) (resp *comment.CommentActionResponse, err error) {

	global.InfoLogger.Info(constant.CommentServiceName,
		zap.String("method", "CommentAction"),
	)

	//视频是否存在
	_, _, err = util.GetVideoFavoriteAndCommentCount(request.VideoId)
	if err != nil {
		global.InfoLogger.Error(constant.VideoExistErrMsg,
			zap.Int64("videoId", request.VideoId),
			zap.Error(err),
		)

		resp = &comment.CommentActionResponse{
			StatusCode: constant.VideoFavoriteCountErrCode,
			StatusMsg:  proto.String(constant.ErrorMsg),
		}
		err = nil
		return
	}

	//用户不存在打印和保存日志并返回
	u, err := GetUser(request.UserId)

	if err != nil {
		global.InfoLogger.Error(constant.FavoriteServiceName,
			zap.Int64("uid", request.UserId),
			zap.String("token", request.Token),
			zap.Error(err),
		)

		return &comment.CommentActionResponse{
			StatusCode: constant.CommentActionErrCode,
			StatusMsg:  proto.String(constant.ErrorMsg),
		}, err
	}

	actionType := request.ActionType

	user := &pb.User{}

	user.Id = u.ID
	user.Avatar = u.Avatar
	user.FavoriteCount = u.FavoriteCount
	user.BackgroundImage = u.BackgroundImage
	user.Signature = u.Signature
	user.IsFollow = u.IsFollow
	user.FollowCount = u.FollowerCount
	user.Name = u.Name
	user.WorkCount = u.WorkCount
	user.TotalFavorited = u.TotalFavorited
	user.FollowerCount = u.FollowCount

	if actionType == 1 {
		resp, err = addComment(ctx, request.VideoId, user, *request.CommentText)
	} else if actionType == 2 {
		commentIdStr := request.CommentId
		commentId, err := strconv.ParseInt(*commentIdStr, 10, 64)
		if err != nil {
			global.InfoLogger.Error(constant.FavoriteServiceName,
				zap.String("commentId", *commentIdStr),
				zap.Error(err),
			)
			resp = &comment.CommentActionResponse{
				StatusCode: constant.CommentActionErrCode,
				StatusMsg:  proto.String(constant.ErrorMsg),
			}
		}

		resp, err = deleteComment(ctx, request.VideoId, commentId, user.Id)
	} else {
		global.InfoLogger.Error(constant.FavoriteServiceName,
			zap.Int64("actionType", actionType),
			zap.Error(err),
		)
		resp = &comment.CommentActionResponse{
			StatusCode: constant.CommentActionErrCode,
			StatusMsg:  proto.String(constant.ErrorMsg),
		}
	}
	return
}

func deleteComment(ctx context.Context, videoId int64, commentId int64, userId int64) (resp *comment.CommentActionResponse, err error) {

	cmt := &model.Comment{}
	cmt, err = dao.Q.WithContext(ctx).Comment.
		Where(dao.Comment.VideoId.Eq(videoId), dao.Comment.CommentId.Eq(commentId)).
		First()
	if err != nil {
		global.InfoLogger.Error(constant.FavoriteServiceName,
			zap.Int64("videoId", videoId),
			zap.Int64("commentId", commentId),
			zap.Error(err),
		)

		resp = &comment.CommentActionResponse{
			StatusCode: constant.CommentActionErrCode,
			StatusMsg:  proto.String(constant.ErrorMsg),
		}
		err = nil
		return
	}
	if cmt.UserId != userId {
		global.InfoLogger.Error(constant.FavoriteServiceName,
			zap.Int64("comment.UserId", cmt.UserId),
			zap.Int64("userId", userId),
			zap.Int64("videoId", videoId),
			zap.Int64("commentId", commentId),
			zap.Error(err),
		)
		resp = &comment.CommentActionResponse{
			StatusCode: constant.CommentActionErrCode,
			StatusMsg:  proto.String(constant.ErrorMsg),
		}
		err = nil
		return
	}

	_, err = dao.Q.WithContext(ctx).Comment.Where(dao.Comment.CommentId.Eq(commentId)).Delete()
	err = util.DelComment(videoId)

	if err != nil {
		global.InfoLogger.Error(constant.FavoriteServiceName,
			zap.Int64("commentId", commentId),
			zap.Error(err),
		)

		resp = &comment.CommentActionResponse{
			StatusCode: constant.CommentActionErrCode,
			StatusMsg:  proto.String(constant.ErrorMsg),
		}
		err = nil
		return
	}

	resp = &comment.CommentActionResponse{
		StatusCode: 0,
		StatusMsg:  proto.String("ok"),
	}
	return
}

func addComment(ctx context.Context, videoId int64, user *pb.User, commentText string) (resp *comment.CommentActionResponse, err error) {

	cmt := &model.Comment{
		VideoId: videoId,
		UserId:  user.Id,
		Content: commentText,
	}

	if err := dao.Q.WithContext(ctx).Comment.Create(cmt); err != nil {
		addCommentErr(cmt, err, resp)
		return resp, err
	}
	util.AddComment(videoId)
	
	resp = &comment.CommentActionResponse{
		StatusCode: 0,
		StatusMsg:  proto.String("ok"),
		Comment: &comment.Comment{
			Id:         cmt.CommentId,
			User:       user,
			Content:    cmt.Content,
			CreateDate: cmt.CreateDate.Format("01-02"),
		},
	}
	return
}

func getUids(comments []*model.Comment) (uids []int64) {
	idMap := sync.Map{}
	ids := make([]int64, 0)
	// 设置最大并发数
	maxConcurrency := constant.MaxConcurrency
	concurrencyCh := make(chan struct{}, maxConcurrency)
	var wg sync.WaitGroup
	for _, cmt := range comments {
		uid := cmt.UserId
		concurrencyCh <- struct{}{} // 占用一个并发槽
		wg.Add(1)
		go func() {
			defer wg.Done()
			defer func() {
				<-concurrencyCh // 释放一个并发槽
			}()
			if _, ok := idMap.Load(uid); ok == false {
				ids = append(ids, uid)
				idMap.Store(uid, true)
			}
		}()
	}
	wg.Wait()
	return ids
}
func GetUsers(ids []int64) []*model.User {
	var users []*model.User
	if err := global.DB.Table("user").Where("id In ?", ids).Find(&users).Error; err != nil {
		global.InfoLogger.Error(constant.FavoriteServiceName,
			zap.String("UserIds", fmt.Sprintf("%+v", ids)),
			zap.Error(err),
		)
		return []*model.User{}
	}
	return users
}
func addCommentErr(cmt *model.Comment, err error, resp *comment.CommentActionResponse) {
	global.InfoLogger.Error(constant.FavoriteServiceName,
		zap.String("Content", cmt.Content),
		zap.Int64("videoId", cmt.VideoId),
		zap.Int64("UserId", cmt.UserId),
		zap.Int64("CommentId", cmt.CommentId),
		zap.Error(err),
	)
	resp = &comment.CommentActionResponse{
		StatusCode: constant.CommentActionErrCode,
		StatusMsg:  proto.String(constant.ErrorMsg),
	}
	err = nil
	return
}