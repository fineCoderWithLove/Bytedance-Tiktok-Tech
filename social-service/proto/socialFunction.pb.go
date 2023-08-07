package proto

import (
	"errors"
	proto "github.com/golang/protobuf/proto"
)

type SocialRelationActionRequest struct {
	Token      string `protobuf:"bytes,1,opt,name=token,proto3" json:"token,omitempty"`
	ToUserId   int64  `protobuf:"varint,2,opt,name=to_user_id,proto3" json:"to_user_id,omitempty"`
	ActionType int32  `protobuf:"varint,3,opt,name=action_type,json=actionType,proto3" json:"action_type,omitempty"`
}

func (m *SocialRelationActionRequest) Reset()         { *m = SocialRelationActionRequest{} }
func (m *SocialRelationActionRequest) String() string { return proto.CompactTextString(m) }
func (*SocialRelationActionRequest) ProtoMessage()    {}

func (m *SocialRelationActionRequest) GetToken() string {
	if m != nil {
		return m.Token
	}
	return ""
}

func (m *SocialRelationActionRequest) GetToUserId() int64 {
	if m != nil {
		return m.ToUserId
	}
	return 0
}

func (m *SocialRelationActionRequest) GetActionType() int32 {
	if m != nil {
		return m.ActionType
	}
	return 0
}

type SocialRelationActionResponse struct {
	StatusCode int32  `protobuf:"varint,1,opt,name=status_code,json=statusCode,proto3" json:"status_code,omitempty"`
	StatusMsg  string `protobuf:"bytes,2,opt,name=status_msg,json=statusMsg,proto3" json:"status_msg,omitempty"`
}

func (m *SocialRelationActionResponse) Reset()         { *m = SocialRelationActionResponse{} }
func (m *SocialRelationActionResponse) String() string { return proto.CompactTextString(m) }
func (*SocialRelationActionResponse) ProtoMessage()    {}

func (m *SocialRelationActionResponse) GetStatusCode() int32 {
	if m != nil {
		return m.StatusCode
	}
	return 0
}

func (m *SocialRelationActionResponse) GetStatusMsg() string {
	if m != nil {
		return m.StatusMsg
	}
	return ""
}

type SocialRelationFollowListRequest struct {
	UserId int64  `protobuf:"varint,1,opt,name=user_id,proto3" json:"user_id,omitempty"`
	Token  string `protobuf:"bytes,2,opt,name=token,proto3" json:"token,omitempty"`
}

func (m *SocialRelationFollowListRequest) Reset()         { *m = SocialRelationFollowListRequest{} }
func (m *SocialRelationFollowListRequest) String() string { return proto.CompactTextString(m) }
func (*SocialRelationFollowListRequest) ProtoMessage()    {}

func (m *SocialRelationFollowListRequest) GetUserId() int64 {
	if m != nil {
		return m.UserId
	}
	return 0
}

func (m *SocialRelationFollowListRequest) GetToken() string {
	if m != nil {
		return m.Token
	}
	return ""
}

type SocialRelationFollowListResponse struct {
	StatusCode int32   `protobuf:"varint,1,opt,name=status_code,json=statusCode,proto3" json:"status_code,omitempty"`
	StatusMsg  string  `protobuf:"bytes,2,opt,name=status_msg,json=statusMsg,proto3" json:"status_msg,omitempty"`
	UserList   []*User `protobuf:"bytes,3,rep,name=user_list,json=userList,proto3" json:"user_list,omitempty"`
}

func (m *SocialRelationFollowListResponse) Reset()         { *m = SocialRelationFollowListResponse{} }
func (m *SocialRelationFollowListResponse) String() string { return proto.CompactTextString(m) }
func (*SocialRelationFollowListResponse) ProtoMessage()    {}

func (m *SocialRelationFollowListResponse) GetStatusCode() int32 {
	if m != nil {
		return m.StatusCode
	}
	return 0
}

func (m *SocialRelationFollowListResponse) GetStatusMsg() string {
	if m != nil {
		return m.StatusMsg
	}
	return ""
}

func (m *SocialRelationFollowListResponse) GetUserList() []*User {
	if m != nil {
		return m.UserList
	}
	return nil
}

type SocialRelationFollowerListRequest struct {
	UserId int64  `protobuf:"varint,1,opt,name=user_id,proto3" json:"user_id,omitempty"`
	Token  string `protobuf:"bytes,2,opt,name=token,proto3" json:"token,omitempty"`
}

func (m *SocialRelationFollowerListRequest) Reset()         { *m = SocialRelationFollowerListRequest{} }
func (m *SocialRelationFollowerListRequest) String() string { return proto.CompactTextString(m) }
func (*SocialRelationFollowerListRequest) ProtoMessage()    {}

func (m *SocialRelationFollowerListRequest) GetUserId() int64 {
	if m != nil {
		return m.UserId
	}
	return 0
}

func (m *SocialRelationFollowerListRequest) GetToken() string {
	if m != nil {
		return m.Token
	}
	return ""
}

type SocialRelationFollowerListResponse struct {
	StatusCode int32   `protobuf:"varint,1,opt,name=status_code,json=statusCode,proto3" json:"status_code,omitempty"`
	StatusMsg  string  `protobuf:"bytes,2,opt,name=status_msg,json=statusMsg,proto3" json:"status_msg,omitempty"`
	UserList   []*User `protobuf:"bytes,3,rep,name=user_list,json=userList,proto3" json:"user_list,omitempty"`
}

func (m *SocialRelationFollowerListResponse) Reset()         { *m = SocialRelationFollowerListResponse{} }
func (m *SocialRelationFollowerListResponse) String() string { return proto.CompactTextString(m) }
func (*SocialRelationFollowerListResponse) ProtoMessage()    {}

func (m *SocialRelationFollowerListResponse) GetStatusCode() int32 {
	if m != nil {
		return m.StatusCode
	}
	return 0
}

func (m *SocialRelationFollowerListResponse) GetStatusMsg() string {
	if m != nil {
		return m.StatusMsg
	}
	return ""
}

func (m *SocialRelationFollowerListResponse) GetUserList() []*User {
	if m != nil {
		return m.UserList
	}
	return nil
}

type SocialRelationFriendListRequest struct {
	UserId int64  `protobuf:"varint,1,opt,name=user_id,proto3" json:"user_id,omitempty"`
	Token  string `protobuf:"bytes,2,opt,name=token,proto3" json:"token,omitempty"`
}

func (m *SocialRelationFriendListRequest) Reset()         { *m = SocialRelationFriendListRequest{} }
func (m *SocialRelationFriendListRequest) String() string { return proto.CompactTextString(m) }
func (*SocialRelationFriendListRequest) ProtoMessage()    {}

func (m *SocialRelationFriendListRequest) GetUserId() int64 {
	if m != nil {
		return m.UserId
	}
	return 0
}

func (m *SocialRelationFriendListRequest) GetToken() string {
	if m != nil {
		return m.Token
	}
	return ""
}

type SocialRelationFriendListResponse struct {
	StatusCode int32   `protobuf:"varint,1,opt,name=status_code,json=statusCode,proto3" json:"status_code,omitempty"`
	StatusMsg  string  `protobuf:"bytes,2,opt,name=status_msg,json=statusMsg,proto3" json:"status_msg,omitempty"`
	UserList   []*User `protobuf:"bytes,3,rep,name=user_list,json=userList,proto3" json:"user_list,omitempty"`
}

func (m *SocialRelationFriendListResponse) Reset()         { *m = SocialRelationFriendListResponse{} }
func (m *SocialRelationFriendListResponse) String() string { return proto.CompactTextString(m) }
func (*SocialRelationFriendListResponse) ProtoMessage()    {}

func (m *SocialRelationFriendListResponse) GetStatusCode() int32 {
	if m != nil {
		return m.StatusCode
	}
	return 0
}

func (m *SocialRelationFriendListResponse) GetStatusMsg() string {
	if m != nil {
		return m.StatusMsg
	}
	return ""
}

func (m *SocialRelationFriendListResponse) GetUserList() []*User {
	if m != nil {
		return m.UserList
	}
	return nil
}

type User struct {
	Id   int64  `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Name string `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
}

func (m *User) Reset()         { *m = User{} }
func (m *User) String() string { return proto.CompactTextString(m) }
func (*User) ProtoMessage()    {}

func (m *User) GetId() int64 {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *User) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

type SocialMessageChatRequest struct {
	Token      string `protobuf:"bytes,1,opt,name=token,proto3" json:"token,omitempty"`
	ToUserId   int64  `protobuf:"varint,2,opt,name=to_user_id,proto3" json:"to_user_id,omitempty"`
	ActionType int32  `protobuf:"varint,3,opt,name=action_type,json=actionType,proto3" json:"action_type,omitempty"`
	Content    string `protobuf:"bytes,4,opt,name=content,proto3" json:"content,omitempty"`
}

func (m *SocialMessageChatRequest) Reset()         { *m = SocialMessageChatRequest{} }
func (m *SocialMessageChatRequest) String() string { return proto.CompactTextString(m) }
func (*SocialMessageChatRequest) ProtoMessage()    {}

func (m *SocialMessageChatRequest) GetToken() string {
	if m != nil {
		return m.Token
	}
	return ""
}

func (m *SocialMessageChatRequest) GetToUserId() int64 {
	if m != nil {
		return m.ToUserId
	}
	return 0
}

func (m *SocialMessageChatRequest) GetActionType() int32 {
	if m != nil {
		return m.ActionType
	}
	return 0
}

func (m *SocialMessageChatRequest) GetContent() string {
	if m != nil {
		return m.Content
	}
	return ""
}

type SocialMessageChatResponse struct {
	StatusCode int32  `protobuf:"varint,1,opt,name=status_code,json=statusCode,proto3" json:"status_code,omitempty"`
	StatusMsg  string `protobuf:"bytes,2,opt,name=status_msg,json=statusMsg,proto3" json:"status_msg,omitempty"`
}

func (m *SocialMessageChatResponse) Reset()         { *m = SocialMessageChatResponse{} }
func (m *SocialMessageChatResponse) String() string { return proto.CompactTextString(m) }
func (*SocialMessageChatResponse) ProtoMessage()    {}

func (m *SocialMessageChatResponse) GetStatusCode() int32 {
	if m != nil {
		return m.StatusCode
	}
	return 0
}

func (m *SocialMessageChatResponse) GetStatusMsg() string {
	if m != nil {
		return m.StatusMsg
	}
	return ""
}

type Message struct {
	Id         int64  `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Content    string `protobuf:"bytes,2,opt,name=content,proto3" json:"content,omitempty"`
	CreateTime int64  `protobuf:"varint,3,opt,name=create_time,proto3" json:"create_time,omitempty"`
	FromUserId int64  `protobuf:"varint,4,opt,name=from_user_id,proto3" json:"from_user_id,omitempty"`
	ToUserId   int64  `protobuf:"varint,5,opt,name=to_user_id,proto3" json:"to_user_id,omitempty"`
}

func (m *Message) Reset()         { *m = Message{} }
func (m *Message) String() string { return proto.CompactTextString(m) }
func (*Message) ProtoMessage()    {}

func (m *Message) GetId() (int64, error) {
	if m == nil {
		return 0, errors.New("Message is nil")
	}
	return m.Id, nil
}

func (m *Message) GetContent() (string, error) {
	if m == nil {
		return "", errors.New("Message is nil")
	}
	return m.Content, nil
}

func (m *Message) GetCreateTime() (int64, error) {
	if m == nil {
		return 0, errors.New("Message is nil")
	}
	return m.CreateTime, nil
}

func (m *Message) GetFromUserId() (int64, error) {
	if m == nil {
		return 0, errors.New("Message is nil")
	}
	return m.FromUserId, nil
}

func (m *Message) GetToUserId() (int64, error) {
	if m == nil {
		return 0, errors.New("Message is nil")
	}
	return m.ToUserId, nil
}

type SocialMessageHistoryRequest struct {
	Token      string `protobuf:"bytes,1,opt,name=token,proto3" json:"token,omitempty"`
	ToUserId   int64  `protobuf:"varint,2,opt,name=to_user_id,proto3" json:"to_user_id,omitempty"`
	PreMsgTime int64  `protobuf:"varint,3,opt,name=pre_msg_time,proto3" json:"pre_msg_time,omitempty"`
}

func (m *SocialMessageHistoryRequest) Reset()         { *m = SocialMessageHistoryRequest{} }
func (m *SocialMessageHistoryRequest) String() string { return proto.CompactTextString(m) }
func (*SocialMessageHistoryRequest) ProtoMessage()    {}

func (m *SocialMessageHistoryRequest) GetToken() string {
	if m == nil {
		return ""
	}
	return m.Token
}

func (m *SocialMessageHistoryRequest) GetToUserId() int64 {
	if m == nil {
		return 0
	}
	return m.ToUserId
}

func (m *SocialMessageHistoryRequest) GetPreMsgTime() int64 {
	if m == nil {
		return 0
	}
	return m.PreMsgTime
}

type SocialMessageHistoryResponse struct {
	StatusCode  int32      `protobuf:"varint,1,opt,name=status_code,json=statusCode,proto3" json:"status_code,omitempty"`
	StatusMsg   string     `protobuf:"bytes,2,opt,name=status_msg,json=statusMsg,proto3" json:"status_msg,omitempty"`
	MessageList []*Message `protobuf:"bytes,3,rep,name=message_list,json=messageList,proto3" json:"message_list,omitempty"`
}

func (m *SocialMessageHistoryResponse) Reset()         { *m = SocialMessageHistoryResponse{} }
func (m *SocialMessageHistoryResponse) String() string { return proto.CompactTextString(m) }
func (*SocialMessageHistoryResponse) ProtoMessage()    {}

func (m *SocialMessageHistoryResponse) GetStatusCode() int32 {
	if m == nil {
		return 0
	}
	return m.StatusCode
}

func (m *SocialMessageHistoryResponse) GetStatusMsg() string {
	if m == nil {
		return ""
	}
	return m.StatusMsg
}

func (m *SocialMessageHistoryResponse) GetMessageList() []*Message {
	if m == nil {
		return nil
	}
	return m.MessageList
}
