package util

import (
	"demotest/douyin-api/global"
	"demotest/douyin-api/globalinit/constant"
	"errors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

// CommentMiddleware 验证 JWT Token 的中间件
func CommentMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从请求头中获取 Authorization 字段
		token := c.Query("token")
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing token"})
			c.Abort()
			return
		}
		uid := global.RS.Get(token).Val()
		if uid == "" {
			zap.S().Errorw("uid错误：", uid)
			c.JSON(constant.CommentActionErrCode, gin.H{
				"error": constant.ErrorMsg,
			})
			c.Abort()
		} else {
			uid, _ := strconv.ParseInt(uid, 10, 64)
			commentText := c.Query("comment_text")
			searchResult := WorldFilter.FindWords(commentText, uid)
			if len(searchResult.Words) != 0 {
				zap.S().Infof("含有敏感词信息:%+v", searchResult)
				c.JSON(constant.CommentActionErrCode, gin.H{
					"status_code":"500",
					"status_msg": errors.New("内容含有敏感信息！").Error(),
				})
				c.Abort()
			}

			c.Next()
		}

	}
}
/*	1.用户注册(不需要token)
 	2.用户发布视频时候的标题(需要token)
	3.用户发送消息的行为(需要token)
*/

// 用户注册的消息中间件,防止用户注册时候使用违法字符
func UserRegisterMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
			commentText := c.Query("username")
			searchResult := WorldFilter.FindWordsNoUserId(commentText)
			if len(searchResult.Words) != 0 {
				zap.S().Infof("含有敏感词信息:%+v", searchResult)
				c.JSON(http.StatusOK, gin.H{
					"status_code":"500",
					"status_msg": errors.New("内容含有敏感信息！").Error(),
				})
				c.Abort()
			}
			c.Next()
		}
}

//用户发布视频的敏感词过滤中间件
func PublishVideoMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
			commentText := c.Request.FormValue("title")
			searchResult := WorldFilter.FindWordsNoUserId(commentText)
			if len(searchResult.Words) != 0 {
				zap.S().Infof("含有敏感词信息:%+v", searchResult)
				c.JSON(http.StatusOK, gin.H{
					"status_code":"500",
					"status_msg": errors.New("内容含有敏感信息！").Error(),
				})
				c.Abort()
			}
			c.Next()
		
	}
}
//用户发送消息的中间件
func MessageActionMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
			commentText := c.Query("content")
			searchResult := WorldFilter.FindWordsNoUserId(commentText)
			if len(searchResult.Words) != 0 {
				zap.S().Infof("含有敏感词信息:%+v", searchResult)
				c.JSON(http.StatusOK, gin.H{
					"status_code":"500",
					"status_msg": errors.New("内容含有敏感信息！").Error(),
				})
				c.Abort()
			}
			c.Next()
		
	}
}
//