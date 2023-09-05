package handler

import (
	"context"
	"demotest/social-service/global"
	mpb "demotest/social-service/proto/message"
	"strconv"
	"time"
)

// TODO 用来用户消息的handler
/*
	MessageChat(context.Context, *MessageChatReq) (*MessageChatResp, error)
	MessageAction(context.Context, *MessageActionReq) (*MessageActionResp, error)
*/

type MessageServer struct {
	mpb.UnimplementedMessageServerServer
}


/*
	1.用户发送消息的接口
*/
type Message struct {
	MsgID      int64 `gorm:"primaryKey;autoIncrement"`
	Content    string
	ToUserID   int64
	UserId     int64
	CreateTime int64
	ActionType int32
}

func (s *MessageServer) MessageAction(ctx context.Context, req *mpb.MessageActionReq) (*mpb.MessageActionResp, error) {
	// TODO 发送消息需要解析token得到自己的userid，然后把自己的userid和传来的touserid一起created到记录中
	uid := global.RS.Get(req.Token).Val()
	resp := mpb.MessageActionResp{}
	Intuid, _ := strconv.Atoi(uid)
	//如果在redis中查询不出来就直接返回
	if Intuid == 0 {
		resp = mpb.MessageActionResp{
			StatusCode: 403,
			StatusMsg:  "请先登录",
		}
		return &resp, nil
	}
	//产生一个当前时间的毫秒级别的时间戳
	currentTime := time.Now()
	timestamp := currentTime.UnixNano() / int64(1e6)
	message := Message{
		Content:    req.Content,
		ToUserID:   req.ToUserId,
		UserId:     int64(Intuid),
		CreateTime: timestamp,
		ActionType: req.ActionType,
	}
	result := global.DB.Table("message").Omit("msg_id").Create(&message)
	if result.Error != nil {
		resp = mpb.MessageActionResp{
			StatusCode: 500,
			StatusMsg:  "database error",
		}
		return &resp, nil
	}
	resp = mpb.MessageActionResp{
		StatusCode: 0,
		StatusMsg:  "success",
	}
	return &resp, nil

}

/*
	2.查看用户聊天记录的接口

*/

/*
	TODO 发送消息的第一次记录没有传递时间,如果时间传来==0那就说明是第一次请求
	如果不等于0，因为前端可能是push数组的，所以需要追加而不是全部查询出来，追加的数组就是新发送的消息，所以需要查询出来比PreMsgTime还要大的时间

	前端的请求就是一直查一直查，所以后端要处理请求的限制
*/

func (s *MessageServer) MessageChat(ctx context.Context, req *mpb.MessageChatReq) (*mpb.MessageChatResp, error) {
	uid := global.RS.Get(req.Token).Val()
	var messages []Message
	var messageList []*mpb.Message
	Intuid, _ := strconv.Atoi(uid)
	resp := mpb.MessageChatResp{}
	if req.PreMsgTime == 0 {
		//如果在redis中查询不出来就直接返回
		if Intuid == 0 {
			resp = mpb.MessageChatResp{
				StatusCode: 403,
				StatusMsg:  "请先登录",
			}
			return &resp, nil
		}
		//查询出来所有的聊天记录，然后转换消息时间的格式
		global.DB.Table("message").Where("(user_id = ? AND to_user_id = ?) OR (user_id = ? AND to_user_id = ?)",
			req.ToUserId, Intuid, Intuid, req.ToUserId).Find(&messages)

		for _, m := range messages {
			CreateTime := strconv.Itoa(int(m.CreateTime))
			message := &mpb.Message{
				Id:         m.MsgID,
				ToUserId:   m.ToUserID,
				FromUserId: m.UserId,
				Content:    m.Content,
				CreateTime: CreateTime,
			}
			messageList = append(messageList, message)
		}
		resp = mpb.MessageChatResp{
			StatusCode:  0,
			StatusMsg:   "success",
			MessageList: messageList,
		}
		return &resp, nil
	} else  {
		//此处需要查询出比LastCreateTime大于等于的记录
		var messages []Message
		_ = global.DB.Table("message").
			Where("user_id = ? AND to_user_id = ? AND create_time > ?", req.ToUserId, Intuid, req.PreMsgTime).
			Order("create_time ASC").
			Find(&messages)
		for _, m := range messages {
			CreateTime := strconv.Itoa(int(m.CreateTime))
			message := &mpb.Message{
				Id:         m.MsgID,
				ToUserId:   m.ToUserID,
				FromUserId: m.UserId,
				Content:    m.Content,
				CreateTime: CreateTime,
			}
			messageList = append(messageList, message)
		}
		resp = mpb.MessageChatResp{
			StatusCode:  0,
			StatusMsg:   "success",
			MessageList: messageList,
		}
		return &resp, nil
	}
}
