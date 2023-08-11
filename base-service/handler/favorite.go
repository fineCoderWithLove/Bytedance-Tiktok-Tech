package handler

import (
	"context"
	"douyin/base-service/dao"
	"douyin/base-service/model"
	"douyin/base-service/proto/favorite"
	"errors"
	"google.golang.org/protobuf/proto"
	"gorm.io/gorm"
)

// TODO 视频流微服务
type FavoriteService struct{}

func (s *FavoriteService) TotalFavorite(ctx context.Context, request *favorite.TotalFavoriteRequest) (resp *favorite.TotalFavoriteResponse, err error) {
	//TODO 通过userid获取所有videoId
	videos := []int64{0, 1}
	videoIds := make([]int64, len("Videos"))
	for _, videoId := range videos {
		videoIds = append(videoIds, videoId)
	}

	var count int64
	count, err = dao.Q.Favorite.
		Where(dao.Favorite.VideoId.In(videoIds...)).
		Count()
	if err != nil {
		//TODO
		resp = &favorite.TotalFavoriteResponse{
			Total:      count,
			StatusCode: 50001,
			StatusMsg:  proto.String("数据库错误"),
		}
	}

	resp = &favorite.TotalFavoriteResponse{
		Total:      count,
		StatusCode: 200,
		StatusMsg:  proto.String("ok"),
	}
	return
}

// UserFavoriteCount 用户喜欢视频总数
func (s *FavoriteService) UserFavoriteCount(ctx context.Context, request *favorite.UserFavoriteCountRequest) (resp *favorite.UserFavoriteCountResponse, err error) {
	var count int64
	count, err = dao.Q.Favorite.
		WithContext(ctx).
		Where(dao.Favorite.UserId.Eq(request.UserId)).
		Count()
	if err != nil {
		//TODO
		resp = &favorite.UserFavoriteCountResponse{
			Count:      0,
			StatusCode: 50001,
			StatusMsg:  proto.String("数据库错误"),
		}
	}
	resp = &favorite.UserFavoriteCountResponse{
		Count:      count,
		StatusCode: 200,
		StatusMsg:  proto.String("ok"),
	}
	return
}

// VideoFavoriteCount 视频获赞总数
func (s *FavoriteService) VideoFavoriteCount(ctx context.Context, request *favorite.VideoFavoriteCountRequest) (resp *favorite.VideoFavoriteCountResponse, err error) {
	//TODO implement me
	var count int64
	count, err = dao.Q.Favorite.
		WithContext(ctx).
		Where(dao.Favorite.VideoId.Eq(request.VideoId)).
		Count()
	if err != nil {
		//TODO
		resp = &favorite.VideoFavoriteCountResponse{
			Count:      0,
			StatusCode: 50001,
			StatusMsg:  proto.String("数据库错误"),
		}
	}
	resp = &favorite.VideoFavoriteCountResponse{
		Count:      count,
		StatusCode: 200,
		StatusMsg:  proto.String("ok"),
	}
	return
}

func (s *FavoriteService) IsFavorite(ctx context.Context, req *favorite.IsFavoriteRequest) (resp *favorite.IsFavoriteResponse, err error) {
	//TODO 视频是否存在

	_, err = dao.Q.WithContext(ctx).Favorite.
		Where(dao.Favorite.UserId.Eq(req.UserId), dao.Favorite.VideoId.Eq(req.VideoId)).First()

	resp = &favorite.IsFavoriteResponse{
		StatusCode: 200,
		IsFavorite: err == nil,
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		resp.StatusCode = 500002
	}
	return
}

/**
TODO 视频流客户端
*/

func (s *FavoriteService) FavoriteAction(ctx context.Context, req *favorite.FavoriteActionRequest) (resp *favorite.FavoriteActionResponse, err error) {
	//TODO 视频是否存在
	if req.ActionType == 1 {
		resp, err = like(ctx, req.UserId, req.VideoId)
	} else {
		resp, err = cancelLike(ctx, req.UserId, req.VideoId)
	}
	return
}

func like(ctx context.Context, uid int64, vid int64) (resp *favorite.FavoriteActionResponse, err error) {

	err = dao.Q.WithContext(ctx).Favorite.Create(&model.Favorite{
		UserId:  uid,
		VideoId: vid,
	})
	if err != nil {
		resp = &favorite.FavoriteActionResponse{
			StatusCode: 500001,
			StatusMsg:  proto.String("点赞失败"),
		}
	}
	resp = &favorite.FavoriteActionResponse{
		StatusCode: 200,
		StatusMsg:  proto.String("ok"),
	}
	return
}

func cancelLike(ctx context.Context, uid int64, vid int64) (resp *favorite.FavoriteActionResponse, err error) {
	_, err = dao.Q.Favorite.WithContext(ctx).Unscoped().
		Where(dao.Favorite.UserId.Eq(uid), dao.Favorite.VideoId.Eq(vid)).
		Delete()
	if err != nil {
		resp = &favorite.FavoriteActionResponse{
			StatusCode: 5000002,
			StatusMsg:  proto.String("点赞失败"),
		}
		return
	}

	resp = &favorite.FavoriteActionResponse{
		StatusCode: 200,
		StatusMsg:  proto.String("ok"),
	}
	return
}
