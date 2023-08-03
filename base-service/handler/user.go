package handler

import (
	"context"
	"douyin/base-service/global"
	"douyin/base-service/model"
	pb "douyin/base-service/proto"
	"fmt"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserServe struct {
	pb.UnimplementedUserServer
}

func ModelToResponse(videos []model.Video) []*pb.VideoInfo {
	var videoInfoResp []*pb.VideoInfo

	for _, video := range videos {
		videoInfo := &pb.VideoInfo{
			VideoId:       video.VideoId,
			PlayUrl:       video.PlayUrl,
			CoverUrl:      video.CoverUrl,
			FavoriteCount: video.FavoriteCount,
			CommentCount:  video.CommentCount,
			Title:         video.Title,
		}

		videoInfoResp = append(videoInfoResp, videoInfo)
	}

	return videoInfoResp
}

func (s *UserServe) GetUserVideo(ctx context.Context, req *pb.UserPrimary) (*pb.VideoInfo, error) {
	var videoList []model.Video
	zap.S().Info("[GetUserVideo] is running")
	res := global.DB.Where("user_id = ?", req.UserId).Find(&videoList)
	if res.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "用户不存在")
	}
	if res.Error != nil {
		return nil, res.Error
	}
	VideoInfo := ModelToResponse(videoList)
	fmt.Println(VideoInfo)
	return VideoInfo[0], nil
}
