package handler

import (
	"douyin/douyin-api/social-service/proto"
	"douyin/douyin-api/social-service/service"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

type IRelationController interface {
	FollowAction(c *gin.Context)  // 关注操作
	GetFollow(c *gin.Context)     // 获取关注的人
	GetFollower(c *gin.Context)   // 获取粉丝
	GetFriendList(c *gin.Context) // 获取朋友列表
}

type RelationController struct {
	RelationService service.IRelationService
}

// GetFollower 获取粉丝
func (r RelationController) GetFollower(c *gin.Context) {
	var SocialRelationFollowerListRequest proto.SocialRelationFollowerListRequest
	err := c.ShouldBindQuery(&SocialRelationFollowerListRequest)
	if err != nil {
		log.Printf("FollowAction|参数错误|%v", err)
		return
	}
	actionResponse := r.RelationService.GetFollower(SocialRelationFollowerListRequest)
	c.JSON(http.StatusOK, actionResponse)
}

// GetFollow 获取关注列表
func (r RelationController) GetFollow(c *gin.Context) {
	var SocialRelationFollowListRequest proto.SocialRelationFollowListRequest
	err := c.ShouldBindQuery(&SocialRelationFollowListRequest)
	if err != nil {
		log.Printf("FollowAction|参数错误|%v", err)
		return
	}
	actionResponse := r.RelationService.GetFollow(SocialRelationFollowListRequest)
	c.JSON(http.StatusOK, actionResponse)
}

// FollowAction 关注操作
func (r RelationController) FollowAction(c *gin.Context) {
	var SocialRelationActionRequest proto.SocialRelationActionRequest
	err := c.ShouldBindQuery(&SocialRelationActionRequest)
	if err != nil {
		log.Printf("FollowAction|参数错误|%v", err)
		return
	}
	value, _ := c.Get("userId")

	actionResponse := r.RelationService.FollowAction(value.(int64), SocialRelationActionRequest)
	c.JSON(http.StatusOK, actionResponse)
}

func (r RelationController) GetFriendList(c *gin.Context) {
	var SocialRelationFriendListRequest proto.SocialRelationFriendListRequest
	err := c.ShouldBindQuery(&SocialRelationFriendListRequest)
	if err != nil {
		log.Printf("FollowAction|参数错误|%v", err)
		return
	}
	value, _ := c.Get("userId")

	actionResponse := r.RelationService.GetFriendList(value.(int64), SocialRelationFriendListRequest)
	c.JSON(http.StatusOK, actionResponse)
}

func NewRelationController() IRelationController {
	relationController := RelationController{RelationService: service.NewRelationService()}
	return relationController
}
