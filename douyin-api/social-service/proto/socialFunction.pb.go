package proto

import (
	proto "github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes/timestamp"
)

// 社交接口---关注操作请求
type SocialRelationActionRequest struct {
	Token      string `protobuf:"bytes,1,opt,name=token,proto3" json:"token,omitempty"`
	ToUserId   int64  `protobuf:"varint,2,opt,name=to_user_id,proto3" json:"to_user_id,omitempty"`
	ActionType int32  `protobuf:"varint,3,opt,name=action_type,json=actionType,proto3" json:"action_type,omitempty"`
}

func (m *SocialRelationActionRequest) Reset()         { *m = SocialRelationActionRequest{} }
func (m *SocialRelationActionRequest) String() string { return proto.CompactTextString(m) }
func (*SocialRelationActionRequest) ProtoMessage()    {}

// 社交接口---关注操作响应
type SocialRelationActionResponse struct {
	//StatusCode int32  `protobuf:"varint,1,opt,name=status_code,json=statusCode,proto3" json:"status_code,omitempty"`
	//StatusMsg  string `protobuf:"bytes,2,opt,name=status_msg,json=statusMsg,proto3" json:"status_msg,omitempty"`
	Response
}

func (m *SocialRelationActionResponse) Reset()         { *m = SocialRelationActionResponse{} }
func (m *SocialRelationActionResponse) String() string { return proto.CompactTextString(m) }
func (*SocialRelationActionResponse) ProtoMessage()    {}

// 社交接口---关注列表请求
type SocialRelationFollowListRequest struct {
	UserId int64  `protobuf:"varint,1,opt,name=user_id,proto3" json:"user_id,omitempty"`
	Token  string `protobuf:"bytes,2,opt,name=token,proto3" json:"token,omitempty"`
}

func (m *SocialRelationFollowListRequest) Reset()         { *m = SocialRelationFollowListRequest{} }
func (m *SocialRelationFollowListRequest) String() string { return proto.CompactTextString(m) }
func (*SocialRelationFollowListRequest) ProtoMessage()    {}

// 社交接口---关注列表响应
type SocialRelationFollowListResponse struct {
	//StatusCode int32   `protobuf:"varint,1,opt,name=status_code,json=statusCode,proto3" json:"status_code,omitempty"`
	//StatusMsg  string  `protobuf:"bytes,2,opt,name=status_msg,json=statusMsg,proto3" json:"status_msg,omitempty"`
	UserList []*User `protobuf:"bytes,3,rep,name=user_list,json=userList,proto3" json:"user_list,omitempty"`
	Response
}

func (m *SocialRelationFollowListResponse) Reset()         { *m = SocialRelationFollowListResponse{} }
func (m *SocialRelationFollowListResponse) String() string { return proto.CompactTextString(m) }
func (*SocialRelationFollowListResponse) ProtoMessage()    {}

// 社交接口---粉丝列表请求
type SocialRelationFollowerListRequest struct {
	UserId int64  `protobuf:"varint,1,opt,name=user_id,proto3" json:"user_id,omitempty"`
	Token  string `protobuf:"bytes,2,opt,name=token,proto3" json:"token,omitempty"`
}

func (m *SocialRelationFollowerListRequest) Reset()         { *m = SocialRelationFollowerListRequest{} }
func (m *SocialRelationFollowerListRequest) String() string { return proto.CompactTextString(m) }
func (*SocialRelationFollowerListRequest) ProtoMessage()    {}

// 社交接口---粉丝列表响应
type SocialRelationFollowerListResponse struct {
	//StatusCode int32   `protobuf:"varint,1,opt,name=status_code,json=statusCode,proto3" json:"status_code,omitempty"`
	//StatusMsg  string  `protobuf:"bytes,2,opt,name=status_msg,json=statusMsg,proto3" json:"status_msg,omitempty"`
	UserList []*User `protobuf:"bytes,3,rep,name=user_list,json=userList,proto3" json:"user_list,omitempty"`
	Response
}

func (m *SocialRelationFollowerListResponse) Reset()         { *m = SocialRelationFollowerListResponse{} }
func (m *SocialRelationFollowerListResponse) String() string { return proto.CompactTextString(m) }
func (*SocialRelationFollowerListResponse) ProtoMessage()    {}

// 社交接口---好友列表请求
type SocialRelationFriendListRequest struct {
	UserId int64  `protobuf:"varint,1,opt,name=user_id,proto3" json:"user_id,omitempty"`
	Token  string `protobuf:"bytes,2,opt,name=token,proto3" json:"token,omitempty"`
}

func (m *SocialRelationFriendListRequest) Reset()         { *m = SocialRelationFriendListRequest{} }
func (m *SocialRelationFriendListRequest) String() string { return proto.CompactTextString(m) }
func (*SocialRelationFriendListRequest) ProtoMessage()    {}

// 社交接口---好友列表响应
type SocialRelationFriendListResponse struct {
	//StatusCode int32   `protobuf:"varint,1,opt,name=status_code,json=statusCode,proto3" json:"status_code,omitempty"`
	//StatusMsg  string  `protobuf:"bytes,2,opt,name=status_msg,json=statusMsg,proto3" json:"status_msg,omitempty"`
	UserList []*User `protobuf:"bytes,3,rep,name=user_list,json=userList,proto3" json:"user_list,omitempty"`
	Response
}

func (m *SocialRelationFriendListResponse) Reset()         { *m = SocialRelationFriendListResponse{} }
func (m *SocialRelationFriendListResponse) String() string { return proto.CompactTextString(m) }
func (*SocialRelationFriendListResponse) ProtoMessage()    {}

type User struct {
	UserId          int64  `protobuf:"varint,1,opt,name=user_id,json=userId" json:"user_id,omitempty"`
	UserName        string `protobuf:"bytes,2,opt,name=username" json:"username,omitempty"`
	Password        string `protobuf:"bytes,3,opt,name=password" json:"password,omitempty"`
	FollowCount     int64  `protobuf:"varint,4,opt,name=follow_count,json=followCount" json:"follow_count,omitempty"`
	FollowerCount   int64  `protobuf:"varint,5,opt,name=follower_count,json=followerCount" json:"follower_count,omitempty"`
	Avatar          string `protobuf:"bytes,6,opt,name=avatar" json:"avatar,omitempty"`
	BackgroundImage string `protobuf:"bytes,7,opt,name=background_image,json=backgroundImage" json:"background_image,omitempty"`
	Signature       string `protobuf:"bytes,8,opt,name=signature" json:"signature,omitempty"`
	TotalFavorited  int64  `protobuf:"varint,9,opt,name=total_favorited,json=totalFavorited" json:"total_favorited,omitempty"`
	WorkCount       int64  `protobuf:"varint,10,opt,name=work_count,json=workCount" json:"work_count,omitempty"`
	FavoriteCount   int64  `protobuf:"varint,11,opt,name=favorite_count,json=favoriteCount" json:"favorite_count,omitempty"`
	IsFollow        bool   `protobuf:"varint,12,opt,name=is_follow,json=isFollow" json:"is_follow"`
}

// 社交接口---发送信息请求
type SocialMessageChatRequest struct {
	Token      string `protobuf:"bytes,1,opt,name=token,proto3" json:"token,omitempty"`
	ToUserId   int64  `protobuf:"varint,2,opt,name=to_user_id,proto3" json:"to_user_id,omitempty"`
	ActionType int32  `protobuf:"varint,3,opt,name=action_type,json=actionType,proto3" json:"action_type,omitempty"`
	Content    string `protobuf:"bytes,4,opt,name=content,proto3" json:"content,omitempty"`
}

func (m *SocialMessageChatRequest) Reset()         { *m = SocialMessageChatRequest{} }
func (m *SocialMessageChatRequest) String() string { return proto.CompactTextString(m) }
func (*SocialMessageChatRequest) ProtoMessage()    {}

// 社交接口---发送信息响应
type SocialMessageChatResponse struct {
	//StatusCode int32  `protobuf:"varint,1,opt,name=status_code,json=statusCode,proto3" json:"status_code,omitempty"`
	//StatusMsg  string `protobuf:"bytes,2,opt,name=status_msg,json=statusMsg,proto3" json:"status_msg,omitempty"`
	Response
}

func (m *SocialMessageChatResponse) Reset()         { *m = SocialMessageChatResponse{} }
func (m *SocialMessageChatResponse) String() string { return proto.CompactTextString(m) }
func (*SocialMessageChatResponse) ProtoMessage()    {}

//	type Message struct {
//		Id         int64  `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
//		Content    string `protobuf:"bytes,2,opt,name=content,proto3" json:"content,omitempty"`
//		CreateTime int64  `protobuf:"varint,3,opt,name=create_time,proto3" json:"create_time,omitempty"`
//		FromUserId int64  `protobuf:"varint,4,opt,name=from_user_id,proto3" json:"from_user_id,omitempty"`
//		ToUserId   int64  `protobuf:"varint,5,opt,name=to_user_id,proto3" json:"to_user_id,omitempty"`
//	}
type Message struct {
	MessageId  int64                `protobuf:"varint,1,opt,name=message_id,json=messageId" json:"message_id,omitempty"`
	UserId     int64                `protobuf:"varint,2,opt,name=user_id,json=userId" json:"user_id,omitempty"`
	ToUserId   int64                `protobuf:"varint,3,opt,name=to_user_id,json=toUserId" json:"to_user_id,omitempty"`
	Content    string               `protobuf:"bytes,4,opt,name=content" json:"content,omitempty"`
	ActionType string               `protobuf:"bytes,5,opt,name=action_type,json=actionType" json:"action_type,omitempty"`
	CreateTime *timestamp.Timestamp `protobuf:"bytes,6,opt,name=create_time,json=createTime" json:"create_time,omitempty"`
	UpdateTime *timestamp.Timestamp `protobuf:"bytes,7,opt,name=update_time,json=updateTime" json:"update_time,omitempty"`
	FromUserId int64                `protobuf:"varint,8,opt,name=from_user_id,proto3" json:"from_user_id,omitempty"`
}

func (Message) TableName() string {
	return "tb_message"
}

func (m *Message) Reset()         { *m = Message{} }
func (m *Message) String() string { return proto.CompactTextString(m) }
func (*Message) ProtoMessage()    {}

// 社交接口---聊天记录请求
type SocialMessageHistoryRequest struct {
	Token      string `protobuf:"bytes,1,opt,name=token,proto3" json:"token,omitempty"`
	ToUserId   int64  `protobuf:"varint,2,opt,name=to_user_id,proto3" json:"to_user_id,omitempty"`
	PreMsgTime int64  `protobuf:"varint,3,opt,name=pre_msg_time,proto3" json:"pre_msg_time,omitempty"`
}

func (m *SocialMessageHistoryRequest) Reset()         { *m = SocialMessageHistoryRequest{} }
func (m *SocialMessageHistoryRequest) String() string { return proto.CompactTextString(m) }
func (*SocialMessageHistoryRequest) ProtoMessage()    {}

// 社交接口---聊天记录响应
type SocialMessageHistoryResponse struct {
	//StatusCode  int32      `protobuf:"varint,1,opt,name=status_code,json=statusCode,proto3" json:"status_code,omitempty"`
	//StatusMsg   string     `protobuf:"bytes,2,opt,name=status_msg,json=statusMsg,proto3" json:"status_msg,omitempty"`
	MessageList []Message `protobuf:"bytes,3,rep,name=message_list,json=messageList,proto3" json:"message_list,omitempty"`
	Response
}

func (m *SocialMessageHistoryResponse) Reset()         { *m = SocialMessageHistoryResponse{} }
func (m *SocialMessageHistoryResponse) String() string { return proto.CompactTextString(m) }
func (*SocialMessageHistoryResponse) ProtoMessage()    {}

type Response struct {
	StatusMsg  string `protobuf:"bytes,1,opt,name=status_msg,json=statusMsg" json:"status_msg,omitempty"`
	StatusCode int32  `protobuf:"varint,2,opt,name=status_code,json=statusCode" json:"status_code,omitempty"`
}

type ErrorMessage struct {
	Response Response `protobuf:"bytes,1,opt,name=response" json:"response,omitempty"`
}

type Relation struct {
	RelationId int64 `protobuf:"varint,1,opt,name=relation_id,json=relationId" json:"relation_id,omitempty"`
	UserId     int64 `protobuf:"varint,2,opt,name=user_id,json=userId" json:"user_id,omitempty"`
	ToUserId   int64 `protobuf:"varint,3,opt,name=to_user_id,json=toUserId" json:"to_user_id,omitempty"`
}

func (m *Relation) Reset()         { *m = Relation{} }
func (m *Relation) String() string { return proto.CompactTextString(m) }
func (*Relation) ProtoMessage()    {}
