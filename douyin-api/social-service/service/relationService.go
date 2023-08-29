package service

import (
	"douyin/douyin-api/global"
	"douyin/douyin-api/social-service/dal/db"

	"douyin/douyin-api/social-service/proto"
	"github.com/jinzhu/copier"
	"log"
	"strconv"
)

type IRelationService interface {
	FollowAction(userId int64, request proto.SocialRelationActionRequest) proto.SocialRelationActionResponse          // 关注操作
	GetFollow(followerRequest proto.SocialRelationFollowListRequest) proto.SocialRelationFollowListResponse           // 获取关注者
	GetFollower(followerRequest proto.SocialRelationFollowerListRequest) proto.SocialRelationFollowerListResponse     // 获取粉丝
	GetFriendList(userId int64, request proto.SocialRelationFriendListRequest) proto.SocialRelationFriendListResponse // 获取朋友列表

}

type RelationService struct {
	RelationRepository db.IRelationRepository
	UserRepository     db.IUserRepository
}

// GetFollower 获取粉丝
func (r RelationService) GetFollower(SocialRelationFollowerListRequest proto.SocialRelationFollowerListRequest) proto.SocialRelationFollowerListResponse {
	userId, err := strconv.ParseInt(strconv.FormatInt(SocialRelationFollowerListRequest.UserId, 10), 10, 64)
	if err != nil {
		log.Printf("GetFollow|格式转换失败|%v", err)
		return proto.SocialRelationFollowerListResponse{}
	}

	relations, err := r.RelationRepository.GetFollower(userId)
	if err != nil {
		log.Printf("GetFollow|数据库错误|%v", err)
		return proto.SocialRelationFollowerListResponse{}
	}

	var SocialRelationFollowerListResponse proto.SocialRelationFollowerListResponse
	for _, relation := range relations {
		toUserId := relation.UserId
		user, err := r.UserRepository.GetUserById(toUserId)
		if err != nil {
			log.Printf("GetFollow|数据库错误|%v", err)
			return proto.SocialRelationFollowerListResponse{}
		}

		var responseUser proto.User
		_ = copier.Copy(&responseUser, &user)
		responseUser.TotalFavorited = user.TotalFavorited
		follow := r.RelationRepository.CheckIsFollow(userId, toUserId)
		responseUser.IsFollow = follow
		myUser := &responseUser
		SocialRelationFollowerListResponse.UserList = append(SocialRelationFollowerListResponse.UserList, myUser)
	}

	SocialRelationFollowerListResponse.Response = proto.Response{
		StatusCode: 0,
		StatusMsg:  "success",
	}

	return SocialRelationFollowerListResponse
}

// GetFollow 获取关注者
func (r RelationService) GetFollow(SocialRelationFollowListRequest proto.SocialRelationFollowListRequest) proto.SocialRelationFollowListResponse {
	userId, err := strconv.ParseInt(strconv.FormatInt(SocialRelationFollowListRequest.UserId, 10), 10, 64)
	if err != nil {
		log.Printf("GetFollow|格式转换失败|%v", err)
		return proto.SocialRelationFollowListResponse{}
	}

	relations, err := r.RelationRepository.GetFollow(userId)
	if err != nil {
		log.Printf("GetFollow|数据库错误|%v", err)
		return proto.SocialRelationFollowListResponse{}
	}

	var SocialRelationFollowListResponse proto.SocialRelationFollowListResponse
	for _, relation := range relations {
		toUserId := relation.ToUserId
		user, err := r.UserRepository.GetUserById(toUserId)
		if err != nil {
			log.Printf("GetFollow|数据库错误|%v", err)
			return proto.SocialRelationFollowListResponse{}
		}
		var responseUser proto.User
		_ = copier.Copy(&responseUser, &user)
		responseUser.TotalFavorited = user.TotalFavorited
		responseUser.IsFollow = r.RelationRepository.CheckIsFollow(userId, toUserId)
		myUser := &responseUser
		SocialRelationFollowListResponse.UserList = append(SocialRelationFollowListResponse.UserList, myUser)

	}

	SocialRelationFollowListResponse.Response = proto.Response{
		StatusCode: 0,
		StatusMsg:  "success",
	}

	return SocialRelationFollowListResponse
}

// FollowAction 关注操作
func (r RelationService) FollowAction(userId int64, SocialRelationActionRequest proto.SocialRelationActionRequest) proto.SocialRelationActionResponse {
	toUserId, err := strconv.ParseInt(strconv.FormatInt(SocialRelationActionRequest.ToUserId, 10), 10, 64)
	if err != nil {
		log.Printf("FollowAction|格式转换失败|%v", err)
		return proto.SocialRelationActionResponse{}
	}

	begin := global.DB.Begin()
	// 关注
	if SocialRelationActionRequest.ActionType == 1 {
		err = r.RelationRepository.AddFollow(userId, toUserId)
		if err != nil {
			log.Printf("FollowAction|插入数据错误|%v", err)
			begin.Rollback()
			return proto.SocialRelationActionResponse{
				Response: proto.Response{StatusCode: 1, StatusMsg: "insert fail"},
			}
		}

		// 增加用户的关注数
		r.UserRepository.UpdateFollowCount(userId, true)

		// 增加用户的获赞总数
		r.UserRepository.UpdateFollowerCount(toUserId, true)
	}

	// 取消关注
	if SocialRelationActionRequest.ActionType == 2 {
		err = r.RelationRepository.RemoveFollow(userId, toUserId)
		if err != nil {
			log.Printf("FollowAction|删除数据错误|%v", err)
			begin.Rollback()
			return proto.SocialRelationActionResponse{
				Response: proto.Response{StatusCode: 1, StatusMsg: "delete fail"},
			}
		}

		// 减少用户的关注数
		r.UserRepository.UpdateFollowCount(userId, false)

		// 减少被关注者的粉丝数
		r.UserRepository.UpdateFollowerCount(toUserId, false)
	}

	begin.Commit()

	return proto.SocialRelationActionResponse{
		Response: proto.Response{StatusCode: 0, StatusMsg: "success"},
	}
}

// GetFriendList 获取朋友列表
func (r RelationService) GetFriendList(userId int64, SocialRelationFriendListRequest proto.SocialRelationFriendListRequest) proto.SocialRelationFriendListResponse {
	userId, err := strconv.ParseInt(strconv.FormatInt(SocialRelationFriendListRequest.UserId, 10), 10, 64)
	if err != nil {
		log.Printf("GetFollow|格式转换失败|%v", err)
		return proto.SocialRelationFriendListResponse{}
	}

	relations, err := r.RelationRepository.GetFollower(userId)

	if err != nil {
		log.Printf("GetFollow|数据库错误|%v", err)
		return proto.SocialRelationFriendListResponse{}
	}

	var SocialRelationFriendListResponse proto.SocialRelationFriendListResponse
	for _, relation := range relations {
		userId := relation.UserId
		user, err := r.UserRepository.GetUserById(userId)
		if err != nil {
			log.Printf("GetFollow|数据库错误|%v", err)
			return proto.SocialRelationFriendListResponse{}
		}
		var responseUser proto.User
		_ = copier.Copy(&responseUser, &user)

		myUser := &responseUser
		SocialRelationFriendListResponse.UserList = append(SocialRelationFriendListResponse.UserList, myUser)
	}

	SocialRelationFriendListResponse.Response = proto.Response{
		StatusCode: 0,
		StatusMsg:  "success",
	}

	return SocialRelationFriendListResponse

}
func NewRelationService() IRelationService {
	relationService := RelationService{RelationRepository: db.NewRelationRepository(), UserRepository: db.NewUserRepository()}
	return relationService
}
