package handler

import (
	"context"
	"demotest/douyin-api/globalinit/constant"
	vp "demotest/douyin-api/proto/video"
	"demotest/interaction-service/dao"
	"demotest/interaction-service/global"
	"demotest/interaction-service/model"
	"demotest/interaction-service/proto/favorite"
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
					VideoId:    3,
					ActionType: 1,
				},
			},
			wantResp: &favorite.FavoriteActionResponse{
				StatusCode: 200,
				StatusMsg:  proto.String("ok"),
			},
		},
		{
			name: "FavoriteAction视频不存在测试",
			args: args{
				ctx: context.Background(),
				req: &favorite.FavoriteActionRequest{
					UserId:     1,
					VideoId:    -19,
					ActionType: 1,
				},
			},
			wantResp: &favorite.FavoriteActionResponse{
				StatusCode: constant.VideoNotExistCode,
				StatusMsg:  proto.String(constant.ErrorMsg),
			},
			wantErr: true,
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

func TestFavoriteService_FavoriteList(t *testing.T) {
	type args struct {
		ctx     context.Context
		request *favorite.FavoriteListRequest
	}
	tests := []struct {
		name     string
		args     args
		wantResp *favorite.FavoriteListResponse
		wantErr  bool
	}{
		// TODO: Add test cases.
		{
			name: "FavoriteList测试",
			args: args{
				ctx:     context.Background(),
				request: &favorite.FavoriteListRequest{UserId: 1},
			},
			wantResp: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &FavoriteService{}
			gotResp, err := s.FavoriteList(tt.args.ctx, tt.args.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("FavoriteList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResp, tt.wantResp) {
				t.Errorf("FavoriteList() gotResp = %v, want %v", gotResp, tt.wantResp)
			}
		})
	}
}

func TestGetVideo(t *testing.T) {
	type args struct {
		vid   int64
		token string
	}
	tests := []struct {
		name string
		args args
		want *vp.Video
	}{
		// TODO: Add test cases.
		{
			name: "TestGetVideo",
			args: args{
				vid:   3,
				token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxLCJleHAiOjE2OTMyMjIwMTl9.K5cSCypJT7LdRiRojd9ihdsZbfY6OkiRcIfzAQHd71Q",
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetVideo(tt.args.vid); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetVideo() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetUser(t *testing.T) {
	type args struct {
		uid int64
	}
	tests := []struct {
		name    string
		args    args
		want    *model.User
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "TestGetUser",
			args: args{uid: 1},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetUser(tt.args.uid)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetUser() got = %v, want %v", got, tt.want)
			}
		})
	}
}
