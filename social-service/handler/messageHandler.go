package handler

import (
	"douyin/social-service/proto"
	"douyin/social-service/service"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

type IChatController interface {
	ChatPost(c *gin.Context)        // 发送消息
	ListChatMessage(c *gin.Context) // 聊天记录列表
}

type ChatController struct {
	ChatService service.IChatService
}

func (c2 ChatController) ListChatMessage(c *gin.Context) {
	var SocialMessageHistoryRequest proto.SocialMessageHistoryRequest

	err := c.ShouldBindQuery(&SocialMessageHistoryRequest)
	if err != nil {
		log.Printf("ChatPost|参数错误|%v", err.Error())
		c.JSON(http.StatusBadRequest, nil)
		return
	}
	userId, _ := c.Get("userId")
	SocialMessageHistoryResponse := c2.ChatService.ListMessage(SocialMessageHistoryRequest, userId.(int64))
	c.JSON(http.StatusOK, SocialMessageHistoryResponse)
}

func (c2 ChatController) ChatPost(c *gin.Context) {
	var SocialMessageChatRequest proto.SocialMessageChatRequest

	err := c.ShouldBindQuery(&SocialMessageChatRequest)
	if err != nil {
		log.Printf("ChatPost|参数错误|%v", SocialMessageChatRequest)
		c.JSON(http.StatusBadRequest, nil)
		return
	}

	userId, _ := c.Get("userId")
	SocialMessageChatResponse := c2.ChatService.PostMessage(SocialMessageChatRequest, userId.(int64))
	c.JSON(http.StatusOK, SocialMessageChatResponse)

}

func NewChatController() IChatController {
	return ChatController{ChatService: service.NewChatService()}
}
