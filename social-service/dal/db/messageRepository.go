package db

import (
	"douyin/social-service/dal/config"
	"douyin/social-service/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
	"strconv"
	"time"
)

type IChatRepository interface {
	AddMessage(userId int64, toUserId string, content string, actionType string) error
	ListMessage(userId int64, toUserID int64, preMsgTime int64) ([]proto.Message, error)
}

type ChatRepository struct {
}

func (c ChatRepository) ListMessage(userId int64, toUserID int64, preMsgTime int64) ([]proto.Message, error) {

	var messages []proto.Message
	timeObj := time.Unix(preMsgTime, 0)
	formattedTime := timeObj.Format("2002-11-02 11:01:12")
	err := config.DB.Table("tb_message").Where("user_id = ? and to_user_id = ? and create_time > ? ", userId, toUserID, formattedTime).Or("user_id = ? and to_user_id = ? and create_time > ?", toUserID, userId, formattedTime).Order("create_time").Find(&messages).Error
	return messages, err
}

func (c ChatRepository) AddMessage(userId int64, toUserId string, content string, actionType string) error {

	var message proto.Message
	message.UserId = userId
	message.ToUserId, _ = strconv.ParseInt(toUserId, 10, 64)
	message.Content = content
	message.ActionType = actionType

	// 转换一下时间类型
	currentTime := time.Now()
	message.CreateTime = timestamppb.New(currentTime)
	message.UpdateTime = timestamppb.New(currentTime)

	err := config.DB.Table("tb_message").Create(&message).Error
	return err
}

func NewChatRepository() IChatRepository {
	return ChatRepository{}
}
