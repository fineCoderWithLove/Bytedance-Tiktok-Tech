package handler

import (
	"context"
	"demotest/interaction-service/dao"
	"demotest/interaction-service/global"
	"demotest/interaction-service/proto/comment"
	"google.golang.org/protobuf/proto"
	"reflect"
	"testing"
)

func init() {
	dao.SetDefault(global.DB)
}

func TestCommentService_CommentList(t *testing.T) {
	type args struct {
		ctx     context.Context
		request *comment.CommentListRequest
	}
	tests := []struct {
		name     string
		args     args
		wantResp *comment.CommentListResponse
		wantErr  bool
	}{
		// TODO: Add test cases.
		{
			name: "CommentList测试",
			args: args{
				ctx:     context.Background(),
				request: &comment.CommentListRequest{VideoId: 1},
			},
			wantResp: &comment.CommentListResponse{
				StatusCode:  200,
				StatusMsg:   proto.String("ok"),
				CommentList: nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &CommentService{}
			gotResp, err := c.CommentList(tt.args.ctx, tt.args.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("CommentList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResp, tt.wantResp) {

				t.Errorf("CommentList() gotResp = %v, want %v", gotResp, tt.wantResp)

			}
		})
	}
}

func TestCommentService_CommentAction(t *testing.T) {
	type args struct {
		ctx     context.Context
		request *comment.CommentActionRequest
	}
	tests := []struct {
		name     string
		args     args
		wantResp *comment.CommentActionResponse
		wantErr  bool
	}{
		// TODO: Add test cases.
		{
			name: "CommentAction添加测试",
			args: args{
				ctx: context.Background(),
				request: &comment.CommentActionRequest{
					Token:       "1",
					VideoId:     1,
					ActionType:  1,
					CommentText: proto.String("777"),
				},
			},
			wantResp: &comment.CommentActionResponse{
				StatusCode: 200,
				StatusMsg:  proto.String("ok"),
				Comment:    nil,
			},
		},
		{
			name: "CommentAction删除测试",
			args: args{
				ctx: context.Background(),
				request: &comment.CommentActionRequest{
					Token:       "1",
					VideoId:     1,
					ActionType:  2,
					CommentText: nil,
					CommentId:   proto.String("2"),
				},
			},
			wantResp: &comment.CommentActionResponse{
				StatusCode: 200,
				StatusMsg:  proto.String("ok"),
				Comment:    nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &CommentService{}
			gotResp, err := c.CommentAction(tt.args.ctx, tt.args.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("CommentAction() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResp, tt.wantResp) {
				t.Errorf("CommentAction() gotResp = %v, want %v", gotResp, tt.wantResp)
			}
		})
	}
}
