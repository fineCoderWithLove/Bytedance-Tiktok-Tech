package handler

import (
	"context"
	pb "demotest/douyin-api/proto/user"
	"demotest/douyin-api/util"
	"demotest/interaction-service/dao"
	"demotest/interaction-service/global"
	"demotest/interaction-service/global/constant"
	"demotest/interaction-service/model"
	"demotest/interaction-service/proto/comment"
	"fmt"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/proto"
	"strconv"
)

var userServiceClient pb.UserServiceClient

func init() {
	conn, err := grpc.Dial(constant.UserServiceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {

		panic(err)
	}
	global.ConsoleLogger.Info(constant.UserServiceClientName,
		zap.String("Addr: ", constant.UserServiceAddr),
	)
	global.InfoLogger.Info(constant.UserServiceClientName,
		zap.String("Addr: ", constant.UserServiceAddr),
	)
	userServiceClient = pb.NewUserServiceClient(conn)
}

type CommentService struct {
}

func (c *CommentService) CommentList(ctx context.Context, request *comment.CommentListRequest) (resp *comment.CommentListResponse, err error) {

	global.ConsoleLogger.Info(constant.CommentServiceName,
		zap.String("method", "CommentList"),
	)

	var exist bool
	exist, err = VideoExist(request.VideoId)
	if !exist || err != nil {
		global.ConsoleLogger.Error(constant.VideoNotExist,
			zap.Int64("videoId", request.VideoId),
			zap.Error(err),
		)
		global.ErrLogger.Error(constant.VideoNotExist,
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

		global.ConsoleLogger.Error(constant.FavoriteServiceName,
			zap.Int64("vid", request.VideoId),
			zap.Error(err),
		)
		global.ErrLogger.Error(constant.FavoriteServiceName,
			zap.Int64("vid", request.VideoId),
			zap.Error(err),
		)

		return &comment.CommentListResponse{
			StatusCode:  constant.CommentListErrCode,
			StatusMsg:   proto.String(constant.ErrorMsg),
			CommentList: nil,
		}, nil
	}

	commentList := make([]*comment.Comment, len(comments))
	for i, cmt := range comments {

		uidStr := strconv.FormatInt(cmt.UserId, 10)
		userDetail, err := userServiceClient.UserDetail(ctx, &pb.DetailRep{
			UserId: uidStr,
			Token:  request.Token,
		})
		if err != nil {
			global.ConsoleLogger.Error(constant.FavoriteServiceName,
				zap.String("uid", uidStr),
				zap.String("token", request.Token),
				zap.Error(err),
			)
			global.ErrLogger.Error(constant.FavoriteServiceName,
				zap.String("uid", uidStr),
				zap.String("token", request.Token),
				zap.Error(err),
			)

			commentList[i] = nil
			continue
		}

		commentList[i] = &comment.Comment{
			Id:         cmt.CommentId,
			User:       userDetail.User,
			Content:    cmt.Content,
			CreateDate: cmt.CreateDate.Format("01-02"),
		}
	}

	resp = &comment.CommentListResponse{
		StatusCode:  0,
		StatusMsg:   proto.String(constant.SuccessMsg),
		CommentList: commentList,
	}
	return
}

func (c *CommentService) CommentAction(ctx context.Context, request *comment.CommentActionRequest) (resp *comment.CommentActionResponse, err error) {

	global.ConsoleLogger.Info(constant.CommentServiceName,
		zap.String("method", "CommentAction"),
	)

	//视频不存在打印和保存日志并返回
	exist, err := VideoExist(request.VideoId)
	if err != nil {
		global.ConsoleLogger.Error(constant.VideoExistErrMsg,
			zap.Int64("videoId", request.VideoId),
			zap.Error(err),
		)
		global.ErrLogger.Error(constant.VideoExistErrMsg,
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
	if !exist {
		global.ConsoleLogger.Error(constant.VideoNotExist,
			zap.Int64("videoId", request.VideoId),
			zap.Error(err),
		)
		global.ErrLogger.Error(constant.VideoNotExist,
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
		global.ConsoleLogger.Error(constant.FavoriteServiceName,
			zap.Int64("uid", request.UserId),
			zap.String("token", request.Token),
			zap.Error(err),
		)
		global.ErrLogger.Error(constant.FavoriteServiceName,
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

	fmt.Println(user.Id,
		user.Avatar,
		user.FavoriteCount,
		user.BackgroundImage,
		user.Signature,
		user.IsFollow,
		user.FollowCount,
		user.Name,
		user.WorkCount,
		user.TotalFavorited,
		user.FollowerCount)

	if actionType == 1 {
		resp, err = addComment(ctx, request.VideoId, user, *request.CommentText)
	} else if actionType == 2 {
		commentIdStr := request.CommentId
		commentId, err := strconv.ParseInt(*commentIdStr, 10, 64)
		if err != nil {
			global.ConsoleLogger.Error(constant.FavoriteServiceName,
				zap.String("commentId", *commentIdStr),
				zap.Error(err),
			)
			global.ErrLogger.Error(constant.FavoriteServiceName,
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
		global.ConsoleLogger.Error(constant.FavoriteServiceName,
			zap.Int64("actionType", actionType),
			zap.Error(err),
		)
		global.ErrLogger.Error(constant.FavoriteServiceName,
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
		global.ConsoleLogger.Error(constant.FavoriteServiceName,
			zap.Int64("videoId", videoId),
			zap.Int64("commentId", commentId),
			zap.Error(err),
		)
		global.ErrLogger.Error(constant.FavoriteServiceName,
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
		global.ConsoleLogger.Error(constant.FavoriteServiceName,
			zap.Int64("comment.UserId", cmt.UserId),
			zap.Int64("userId", userId),
			zap.Int64("videoId", videoId),
			zap.Int64("commentId", commentId),
			zap.Error(err),
		)
		global.ErrLogger.Error(constant.FavoriteServiceName,
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
		global.ConsoleLogger.Error(constant.FavoriteServiceName,
			zap.Int64("commentId", commentId),
			zap.Error(err),
		)
		global.ErrLogger.Error(constant.FavoriteServiceName,
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
	var commentId int64
	commentId, err = dao.Q.Comment.
		Unscoped().
		Count()

	if err != nil {
		global.ConsoleLogger.Error(constant.FavoriteServiceName,
			zap.Int64("videoId", videoId),
			zap.Error(err),
		)
		global.ErrLogger.Error(constant.FavoriteServiceName,
			zap.Int64("videoId", videoId),
			zap.Error(err),
		)
		resp = &comment.CommentActionResponse{
			StatusCode: constant.CommentActionErrCode,
			StatusMsg:  proto.String(constant.ErrorMsg),
		}
		err = nil
		return
	}

	cmt := &model.Comment{
		CommentId: commentId + 1,
		VideoId:   videoId,
		UserId:    user.Id,
		Content:   commentText,
	}

	err = dao.Q.WithContext(ctx).Comment.Create(cmt)
	err = util.AddComment(videoId)

	if err != nil {

		global.ConsoleLogger.Error(constant.FavoriteServiceName,
			zap.String("Content", cmt.Content),
			zap.Int64("videoId", cmt.VideoId),
			zap.Int64("UserId", cmt.UserId),
			zap.Int64("CommentId", cmt.CommentId),
			zap.Error(err),
		)
		global.ErrLogger.Error(constant.FavoriteServiceName,
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
