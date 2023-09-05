package handler

import (
	"context"
	"demotest/social-service/global"
	"demotest/social-service/util"
	"strconv"

	"demotest/social-service/proto/relation"
	"fmt"
)

// TODO 用来做用户关系的接口

type RelationServer struct {
	rpb.UnimplementedUserServiceServer
}

/*
	1.用户关注操作接口 Post
*/
type Attention struct {
	UserID   int64
	ToUserID int64
	// 其他列...
}
type User struct {
	ID              int64
	Name            string
	FollowCount     int64
	FollowerCount   int64
	IsFollow        bool
	Avatar          string
	BackgroundImage string
	Signature       string
	TotalFavorited  int64
	WorkCount       int64
	FavoriteCount   int64
}

func (s *RelationServer) UserRelationAction(ctx context.Context, req *rpb.UserRelationActionReq) (*rpb.UserRelationActionResp, error) {
	/*
		TODO 解析出来自己的uid然后，拿到获取的id然后添加记录进去
	*/
	resp := rpb.UserRelationActionResp{}
	uid := global.RS.Get(req.Token).Val()
	Intuid, _ := strconv.Atoi(uid)
	//如果在redis中查询不出来就直接返回
	if Intuid == 0 {
		resp = rpb.UserRelationActionResp{
			StatusCode: 500,
			StatusMsg:  "请先登录",
		}
		return &resp, nil
	}
	//检测是否对自己操作
	if req.ToUserId == int64(Intuid) {
		resp = rpb.UserRelationActionResp{
			StatusCode: 500,
			StatusMsg:  "不能对自己操作",
		}
		return &resp, nil
	}

	//添加一条记录到attention表中 TODO 关注成功后还要往数据库或者redis中的数据+1
	attention := Attention{}
	if Intuid != 0 && req.ToUserId != 0 {
		attention = Attention{
			UserID:   int64(Intuid),
			ToUserID: req.ToUserId,
		}
	}
	if req.ActionType == 1 {
		//表示为关注操作
		if err := global.DB.Table("attention").Create(&attention).Error; err != nil {
			fmt.Println(err)
			resp = rpb.UserRelationActionResp{
				StatusCode: 500,
				StatusMsg:  "fail to insert data",
			}
			return &resp, err
		} else {
			//否则就为关注成功,不仅要把关注数目+1，还要给被关注者粉丝数目+1
			util.IncreaseFollowCount(int64(Intuid))
			util.IncreaseFollowerCount(req.ToUserId)
		}
	} else {
		//表示为取消关注操作
		result := global.DB.Table("attention").Delete(&Attention{}, "user_id = ? AND to_user_id = ?", Intuid, req.ToUserId)
		if result.RowsAffected != 0 {
			//说明成功取消关注,不仅要把关注数目-1，还要给被取消关注者粉丝数目-1
			util.DecreaseFollowCount(global.RS, int64(Intuid))
			util.DecreaseFollowerCount(req.ToUserId)
			resp = rpb.UserRelationActionResp{
				StatusCode: 0,
				StatusMsg:  "success",
			}
			return &resp, nil
		} else {
			resp = rpb.UserRelationActionResp{
				StatusCode: 500,
				StatusMsg:  "faili to delete data",
			}
			return &resp, nil
		}
	}

	resp = rpb.UserRelationActionResp{
		StatusCode: 0,
		StatusMsg:  "success",
	}
	return &resp, nil
}

/*
	2.查询用户的关注列表接口  Get
*/
func (s *RelationServer) UserRelationFollowList(ctx context.Context, req *rpb.UserRelationFollowListReq) (*rpb.UserRelationFollowListResp, error) {
	resp := rpb.UserRelationFollowListResp{}
	//如果在redis中查询不出来就直接返回
	uid := global.RS.Get(req.Token).Val()
	Intuid, _ := strconv.Atoi(uid)

	if Intuid == 0 {
		resp = rpb.UserRelationFollowListResp{
			StatusCode: 500,
			StatusMsg:  "请先登录",
			UserList:   nil,
		}
		return &resp, nil
	}
	var users []User
	global.DB.Select("user.*").Table("user").
		Joins("JOIN attention ON user.id = attention.to_user_id").
		Where("attention.user_id = ?", req.UserId).
		Find(&users)
	var userList []*rpb.User
	// TODO 记得处理IsFollow的信息
	for _, v := range users {
		//查询视频流的时候要从redis中查询出用户的所有信息
		TotalFavorited, FavoriteCount, FollowerCount, FollowCount, _ := util.GetUserAllUserData(int64(v.ID))
		user := &rpb.User{
			Id:            v.ID,
			Name:          v.Name,
			FollowCount:   FollowCount,
			FollowerCount: FollowerCount,
			//TODO isFollow要进行处理
			IsFollow: true,
			Avatar:   v.Avatar,
			//
			BackgroundImage: v.BackgroundImage,
			Signature:       v.Signature,
			TotalFavorited:  TotalFavorited,
			WorkCount:       v.WorkCount,
			FavoriteCount:   FavoriteCount,
		}
		userList = append(userList, user)
	}
	resp = rpb.UserRelationFollowListResp{
		StatusCode: 0,
		StatusMsg:  "success",
		UserList:   userList,
	}

	return &resp, nil
}

/*
	3.查询用户粉丝列表接口 Get
*/
func (s *RelationServer) UserRelationFollowerList(ctx context.Context, req *rpb.UserRelationFollowerListReq) (*rpb.UserRelationFollowerListResp, error) {
	resp := rpb.UserRelationFollowerListResp{}
	//如果在redis中查询不出来就直接返回
	uid := global.RS.Get(req.Token).Val()
	Intuid, _ := strconv.Atoi(uid)
	if Intuid == 0 {
		resp = rpb.UserRelationFollowerListResp{
			StatusCode: 500,
			StatusMsg:  "请先登录",
			UserList:   nil,
		}
		return &resp, nil
	}
	var users []User
	global.DB.Select("user.*").
		Table("user").
		Joins("JOIN attention ON user.id = attention.user_id AND attention.to_user_id = ?", req.UserId).
		Find(&users)
	var userList []*rpb.User
	// TODO 记得处理IsFollow的信息
	for _, v := range users {
		var Isfollow bool
		var attentions []Attention
		result := global.DB.Table("attention").Where("user_id = ? AND to_user_id = ?", Intuid, v.ID).Find(&attentions)
		if result.RowsAffected != 0 {
			Isfollow = true
		}
		//查询视频流的时候要从redis中查询出用户的所有信息
		TotalFavorited, FavoriteCount, FollowerCount, FollowCount, _ := util.GetUserAllUserData(int64(v.ID))
		user := &rpb.User{
			Id:            v.ID,
			Name:          v.Name,
			FollowCount:   FollowCount,
			FollowerCount: FollowerCount,
			//TODO isFollow要进行处理
			IsFollow: Isfollow,
			Avatar:   v.Avatar,
			//
			BackgroundImage: v.BackgroundImage,
			Signature:       v.Signature,
			TotalFavorited:  TotalFavorited,
			WorkCount:       v.WorkCount,
			FavoriteCount:   FavoriteCount,
		}
		userList = append(userList, user)
	}
	resp = rpb.UserRelationFollowerListResp{
		StatusCode: 0,
		StatusMsg:  "success",
		UserList:   userList,
	}

	return &resp, nil
}

/*
	4.查询用户好友列表接口  Get
*/
func (s *RelationServer) UserRelationFriendList(ctx context.Context, req *rpb.UserRelationFriendListReq) (*rpb.UserRelationFriendListResp, error) {
	//查询好友列表的关键就是给定一个userid和touserid要求都有查询到
	resp := rpb.UserRelationFriendListResp{}
	//如果在redis中查询不出来就直接返回
	uid := global.RS.Get(req.Token).Val()
	Intuid, _ := strconv.Atoi(uid)
	if Intuid == 0 {
		resp = rpb.UserRelationFriendListResp{
			StatusCode: 500,
			StatusMsg:  "请先登录",
			UserList:   nil,
		}
		return &resp, nil
	}
	var friends []User
	err := global.DB.Table("user").
		Select("user.*").
		Joins("JOIN attention AS a1 ON user.id = a1.user_id AND a1.to_user_id = ? "+
			"JOIN attention AS a2 ON user.id = a2.to_user_id AND a2.user_id = ?", req.UserId, req.UserId).
		Find(&friends).Error
	if err != nil {
		return nil, err
	}
	var userList []*rpb.FriendUser
	// TODO 记得处理IsFollow的信息
	for _, v := range friends {
		//查询一次最新聊天记录
		fmt.Println(Intuid)
		var message Message
		if err := global.DB.Table("message").
			Select("*").
			Where("user_id = ? AND to_user_id = ?", Intuid, v.ID).
			Order("create_time DESC").
			Limit(1).
			Find(&message).Error; err != nil {
			// 处理错误
		}

		//查询视频流的时候要从redis中查询出用户的所有信息
		TotalFavorited, FavoriteCount, FollowerCount, FollowCount, _ := util.GetUserAllUserData(int64(v.ID))
		user := &rpb.FriendUser{
			Id:            v.ID,
			Name:          v.Name,
			FollowCount:   FollowCount,
			FollowerCount: FollowerCount,
			//TODO isFollow要进行处理
			IsFollow: true,
			Avatar:   v.Avatar,
			//
			BackgroundImage: v.BackgroundImage,
			Signature:       v.Signature,
			TotalFavorited:  TotalFavorited,
			WorkCount:       v.WorkCount,
			FavoriteCount:   FavoriteCount,
			Message:         message.Content,
			MsgType:         int64(message.ActionType),
		}
		userList = append(userList, user)
	}
	resp = rpb.UserRelationFriendListResp{
		StatusCode: 0,
		StatusMsg:  "success",
		UserList:   userList,
	}

	return &resp, nil
}

