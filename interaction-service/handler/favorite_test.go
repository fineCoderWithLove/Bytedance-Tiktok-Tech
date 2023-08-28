package handler

import (
	"context"
	"douyin/douyin-api/globalinit/constant"
	vp "douyin/douyin-api/proto/video"
	"google.golang.org/protobuf/proto"
	"interaction-service/dao"
	"interaction-service/global"
	"interaction-service/model"
	"interaction-service/proto/favorite"
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
		{
			name: "UserFavoriteCount测试",
			args: args{
				ctx: context.Background(),
				request: &favorite.UserFavoriteCountRequest{
					UserId: 1,
				},
			},
			wantResp: &favorite.UserFavoriteCountResponse{
				Count:      15,
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
		{
			name: "TotalFavorite测试",
			args: args{
				ctx: context.Background(),
				request: &favorite.TotalFavoriteRequest{
					UserId: 20,
				},
			},
			wantResp: &favorite.TotalFavoriteResponse{
				Total:      3,
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

func TestVideoIds(t *testing.T) {
	type args struct {
		uid int64
	}
	tests := []struct {
		name    string
		args    args
		want    *[]int64
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "获取视频ids测试",
			args: args{uid: 20},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := VideoIds(tt.args.uid)
			if (err != nil) != tt.wantErr {
				t.Errorf("VideoIds() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("VideoIds() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVideoExist(t *testing.T) {
	type args struct {
		vid int64
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "TestVideoExist测试",
			args: args{vid: 3},
			want: true,
		},
		{
			name: "TestVideoExist测试",
			args: args{vid: 30},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := VideoExist(tt.args.vid)
			if (err != nil) != tt.wantErr {
				t.Errorf("VideoExist() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("VideoExist() got = %v, want %v", got, tt.want)
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
			if got := GetVideo(tt.args.vid, tt.args.token); !reflect.DeepEqual(got, tt.want) {
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
