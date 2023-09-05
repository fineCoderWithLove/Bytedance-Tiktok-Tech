package main

import (
	"context"
	"douyin/douyin-api/util"

	"fmt"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"interaction-service/proto/favorite"
)

var favoriteSrvClient favorite.FavoriteServiceClient

func init() {
	conn, err := grpc.Dial("localhost:8881", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		zap.S().Errorw("[FavoriteAction]连接【点赞服务失败】", "msg", err.Error())
	}
	favoriteSrvClient = favorite.NewFavoriteServiceClient(conn)
}

func main() {
	//TotalFavorite()
	//UserFavoriteCount()
	//VideoFavoriteCount()
	//IsFavorite()
	//FavoriteAction()
	redisTest()
}

/*
	func TotalFavorite() {
		favoriteResponse, err := favoriteSrvClient.TotalFavorite(context.Background(), &favorite.TotalFavoriteRequest{
			UserId: 2})
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(favoriteResponse.Total)
	}

	func UserFavoriteCount() {
		resp, err := favoriteSrvClient.UserFavoriteCount(context.Background(), &favorite.UserFavoriteCountRequest{
			UserId: 1})
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(resp.Count)
	}

	func VideoFavoriteCount() {
		response, err := favoriteSrvClient.VideoFavoriteCount(context.Background(), &favorite.VideoFavoriteCountRequest{
			VideoId: 100})
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(response.Count)
	}

	func IsFavorite() {
		resp, err := favoriteSrvClient.IsFavorite(context.Background(), &favorite.IsFavoriteRequest{
			UserId:  10,
			VideoId: 100,
		})
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(resp.IsFavorite)
	}
*/
func FavoriteAction() {
	resp, err := favoriteSrvClient.FavoriteAction(context.Background(), &favorite.FavoriteActionRequest{
		UserId:     1,
		VideoId:    22,
		ActionType: 1,
	})
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(resp.StatusCode)
}
func redisTest() {
	/*//util.InitializeRedisData(global.DB, global.RS)
	var count int64
	count, _ = util.GetFavoriteCount(3)
	fmt.Println("初始favoriteCount:", count)

	util.LikeVideo(3)
	count, _ = util.GetFavoriteCount(3)
	fmt.Println("点赞后favoriteCount:", count)

	main2.UnlikeVideo(3)
	count, _ = util.GetFavoriteCount(3)
	fmt.Println("取消点赞后favoriteCount:", count)

	count, _ = util.GetCommentCount(3)
	fmt.Println("初始视频CommentCount:", count)

	main2.AddComment(3)
	count, _ = util.GetCommentCount(3)
	fmt.Println("添加评论后CommentCount:", count)

	//fmt.Println("更新mysql")
	main2.LikeVideo(3)
	main2.LikeVideo(3)
	main2.LikeVideo(3)
	main2.LikeVideo(3)*/
	//util.WriteRedisVideoToMySQL(global.DB)
	util.GetVideoFavoriteAndCommentCount(-100)
}
