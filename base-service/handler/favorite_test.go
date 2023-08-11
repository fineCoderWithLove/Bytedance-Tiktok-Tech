package handler

import (
	"context"
	"douyin/base-service/dao"
	"douyin/base-service/global"
	"douyin/base-service/proto/favorite"
	"google.golang.org/protobuf/proto"
	"reflect"
	"testing"
)

func init() {
	dao.SetDefault(global.DB)
}

func TestFavoriteService_IsFavorite(t *testing.T) {

	type args struct {
		ctx context.Context
		req *favorite.IsFavoriteRequest
	}
	tests := []struct {
		name     string
		args     args
		wantResp *favorite.IsFavoriteResponse
		wantErr  bool
	}{
		// TODO: Add test cases.
		{
			name: "IsFavorite_ture测试",
			args: args{
				ctx: context.Background(),
				req: &favorite.IsFavoriteRequest{
					UserId:  1,
					VideoId: 1,
				},
			},
			wantResp: &favorite.IsFavoriteResponse{
				StatusCode: 200,
				IsFavorite: true,
			},
		},
		{
			name: "IsFavorite_false测试",
			args: args{
				ctx: context.Background(),
				req: &favorite.IsFavoriteRequest{
					UserId:  1,
					VideoId: 2,
				},
			},
			wantErr: true,
			wantResp: &favorite.IsFavoriteResponse{
				StatusCode: 200,
				IsFavorite: false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &FavoriteService{}
			gotResp, err := s.IsFavorite(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("IsFavorite() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResp, tt.wantResp) {
				t.Errorf("IsFavorite() gotResp = %v, want %v", gotResp, tt.wantResp)
			}
		})
	}
}

func TestFavoriteService_FavoriteAction(t *testing.T) {
	type args struct {
		ctx context.Context
		req *favorite.FavoriteActionRequest
	}
	tests := []struct {
		name     string
		args     args
		wantResp *favorite.FavoriteActionResponse
		wantErr  bool
	}{
		// TODO: Add test cases.
		{
			name: "canselLike测试",
			args: args{
				ctx: context.Background(),
				req: &favorite.FavoriteActionRequest{
					UserId:     1,
					VideoId:    1,
					ActionType: 0,
				},
			},
			wantResp: &favorite.FavoriteActionResponse{
				StatusCode: 200,
				StatusMsg:  proto.String("ok"),
			},
		},
		{
			name: "Like测试",
			args: args{
				ctx: context.Background(),
				req: &favorite.FavoriteActionRequest{
					UserId:     1,
					VideoId:    1,
					ActionType: 1,
				},
			},
			wantResp: &favorite.FavoriteActionResponse{
				StatusCode: 200,
				StatusMsg:  proto.String("ok"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &FavoriteService{}
			gotResp, err := s.FavoriteAction(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("FavoriteAction() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResp, tt.wantResp) {
				t.Errorf("FavoriteAction() gotResp = %v, want %v", gotResp, tt.wantResp)
			}
		})
	}
}

func TestFavoriteService_UserFavoriteCount(t *testing.T) {
	type args struct {
		ctx     context.Context
		request *favorite.UserFavoriteCountRequest
	}
	tests := []struct {
		name     string
		args     args
		wantResp *favorite.UserFavoriteCountResponse
		wantErr  bool
	}{
		// TODO: Add test cases.
		{
			name: "UserFavoriteCount测试",
			args: args{
				ctx: context.Background(),
				request: &favorite.UserFavoriteCountRequest{
					UserId: 1,
				},
			},
			wantResp: &favorite.UserFavoriteCountResponse{
				Count:      8,
				StatusCode: 200,
				StatusMsg:  proto.String("ok"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &FavoriteService{}
			gotResp, err := s.UserFavoriteCount(tt.args.ctx, tt.args.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("UserFavoriteCount() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResp, tt.wantResp) {
				t.Errorf("UserFavoriteCount() gotResp = %v, want %v", gotResp, tt.wantResp)
			}
		})
	}
}

func TestFavoriteService_VideoFavoriteCount(t *testing.T) {
	type args struct {
		ctx     context.Context
		request *favorite.VideoFavoriteCountRequest
	}
	tests := []struct {
		name     string
		args     args
		wantResp *favorite.VideoFavoriteCountResponse
		wantErr  bool
	}{
		// TODO: Add test cases.
		{
			name: "VideoFavoriteCount测试",
			args: args{
				ctx: context.Background(),
				request: &favorite.VideoFavoriteCountRequest{
					VideoId: 0,
				},
			},
			wantResp: &favorite.VideoFavoriteCountResponse{
				Count:      7,
				StatusCode: 200,
				StatusMsg:  proto.String("ok"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &FavoriteService{}
			gotResp, err := s.VideoFavoriteCount(tt.args.ctx, tt.args.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("VideoFavoriteCount() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResp, tt.wantResp) {
				t.Errorf("VideoFavoriteCount() gotResp = %v, want %v", gotResp, tt.wantResp)
			}
		})
	}
}

func TestFavoriteService_TotalFavorite(t *testing.T) {
	type args struct {
		ctx     context.Context
		request *favorite.TotalFavoriteRequest
	}
	tests := []struct {
		name     string
		args     args
		wantResp *favorite.TotalFavoriteResponse
		wantErr  bool
	}{
		// TODO: Add test cases.
		{
			name: "TotalFavorite测试",
			args: args{
				ctx: context.Background(),
				request: &favorite.TotalFavoriteRequest{
					UserId: 1,
				},
			},
			wantResp: &favorite.TotalFavoriteResponse{
				Total:      8,
				StatusCode: 200,
				StatusMsg:  proto.String("ok"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &FavoriteService{}
			gotResp, err := s.TotalFavorite(tt.args.ctx, tt.args.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("TotalFavorite() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResp, tt.wantResp) {
				t.Errorf("TotalFavorite() gotResp = %v, want %v", gotResp, tt.wantResp)
			}
		})
	}
}
