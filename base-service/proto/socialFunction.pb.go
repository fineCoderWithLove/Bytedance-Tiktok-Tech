package __

import (
	proto "github.com/golang/protobuf/proto"
)

type SocialFollowFunctionRequest struct {
	Token      string `protobuf:"bytes,1,opt,name=token,proto3" json:"token,omitempty"`
	ToUserId   int64  `protobuf:"varint,2,opt,name=to_user_id,proto3" json:"to_user_id,omitempty"`
	ActionType int32  `protobuf:"varint,3,opt,name=action_type,json=actionType,proto3" json:"action_type,omitempty"`
}

func (m *SocialFollowFunctionRequest) Reset()         { *m = SocialFollowFunctionRequest{} }
func (m *SocialFollowFunctionRequest) String() string { return proto.CompactTextString(m) }
func (*SocialFollowFunctionRequest) ProtoMessage()    {}

func (m *SocialFollowFunctionRequest) GetToken() string {
	if m != nil {
		return m.Token
	}
	return ""
}

func (m *SocialFollowFunctionRequest) GetToUserId() int64 {
	if m != nil {
		return m.ToUserId
	}
	return 0
}

func (m *SocialFollowFunctionRequest) GetActionType() int32 {
	if m != nil {
		return m.ActionType
	}
	return 0
}

type SocialFollowFunctionResponse struct {
	StatusCode int32  `protobuf:"varint,1,opt,name=status_code,json=statusCode,proto3" json:"status_code,omitempty"`
	StatusMsg  string `protobuf:"bytes,2,opt,name=status_msg,json=statusMsg,proto3" json:"status_msg,omitempty"`
}

func (m *SocialFollowFunctionResponse) Reset()         { *m = SocialFollowFunctionResponse{} }
func (m *SocialFollowFunctionResponse) String() string { return proto.CompactTextString(m) }
func (*SocialFollowFunctionResponse) ProtoMessage()    {}

func (m *SocialFollowFunctionResponse) GetStatusCode() int32 {
	if m != nil {
		return m.StatusCode
	}
	return 0
}

func (m *SocialFollowFunctionResponse) GetStatusMsg() string {
	if m != nil {
		return m.StatusMsg
	}
	return ""
}

type SocialFollowListRequest struct {
	UserId int64  `protobuf:"varint,1,opt,name=user_id,proto3" json:"user_id,omitempty"`
	Token  string `protobuf:"bytes,2,opt,name=token,proto3" json:"token,omitempty"`
}

func (m *SocialFollowListRequest) Reset()         { *m = SocialFollowListRequest{} }
func (m *SocialFollowListRequest) String() string { return proto.CompactTextString(m) }
func (*SocialFollowListRequest) ProtoMessage()    {}

func (m *SocialFollowListRequest) GetUserId() int64 {
	if m != nil {
		return m.UserId
	}
	return 0
}

func (m *SocialFollowListRequest) GetToken() string {
	if m != nil {
		return m.Token
	}
	return ""
}

type SocialFollowListResponse struct {
	StatusCode int32   `protobuf:"varint,1,opt,name=status_code,json=statusCode,proto3" json:"status_code,omitempty"`
	StatusMsg  string  `protobuf:"bytes,2,opt,name=status_msg,json=statusMsg,proto3" json:"status_msg,omitempty"`
	UserList   []*User `protobuf:"bytes,3,rep,name=user_list,json=userList,proto3" json:"user_list,omitempty"`
}

func (m *SocialFollowListResponse) Reset()         { *m = SocialFollowListResponse{} }
func (m *SocialFollowListResponse) String() string { return proto.CompactTextString(m) }
func (*SocialFollowListResponse) ProtoMessage()    {}

func (m *SocialFollowListResponse) GetStatusCode() int32 {
	if m != nil {
		return m.StatusCode
	}
	return 0
}

func (m *SocialFollowListResponse) GetStatusMsg() string {
	if m != nil {
		return m.StatusMsg
	}
	return ""
}

func (m *SocialFollowListResponse) GetUserList() []*User {
	if m != nil {
		return m.UserList
	}
	return nil
}

type SocialFollowerListRequest struct {
	UserId int64  `protobuf:"varint,1,opt,name=user_id,proto3" json:"user_id,omitempty"`
	Token  string `protobuf:"bytes,2,opt,name=token,proto3" json:"token,omitempty"`
}

func (m *SocialFollowerListRequest) Reset()         { *m = SocialFollowerListRequest{} }
func (m *SocialFollowerListRequest) String() string { return proto.CompactTextString(m) }
func (*SocialFollowerListRequest) ProtoMessage()    {}

func (m *SocialFollowerListRequest) GetUserId() int64 {
	if m != nil {
		return m.UserId
	}
	return 0
}

func (m *SocialFollowerListRequest) GetToken() string {
	if m != nil {
		return m.Token
	}
	return ""
}

type SocialFollowerListResponse struct {
	StatusCode int32   `protobuf:"varint,1,opt,name=status_code,json=statusCode,proto3" json:"status_code,omitempty"`
	StatusMsg  string  `protobuf:"bytes,2,opt,name=status_msg,json=statusMsg,proto3" json:"status_msg,omitempty"`
	UserList   []*User `protobuf:"bytes,3,rep,name=user_list,json=userList,proto3" json:"user_list,omitempty"`
}

func (m *SocialFollowerListResponse) Reset()         { *m = SocialFollowerListResponse{} }
func (m *SocialFollowerListResponse) String() string { return proto.CompactTextString(m) }
func (*SocialFollowerListResponse) ProtoMessage()    {}

func (m *SocialFollowerListResponse) GetStatusCode() int32 {
	if m != nil {
		return m.StatusCode
	}
	return 0
}

func (m *SocialFollowerListResponse) GetStatusMsg() string {
	if m != nil {
		return m.StatusMsg
	}
	return ""
}

func (m *SocialFollowerListResponse) GetUserList() []*User {
	if m != nil {
		return m.UserList
	}
	return nil
}

type SocialFriendListRequest struct {
	UserId int64  `protobuf:"varint,1,opt,name=user_id,proto3" json:"user_id,omitempty"`
	Token  string `protobuf:"bytes,2,opt,name=token,proto3" json:"token,omitempty"`
}

func (m *SocialFriendListRequest) Reset()         { *m = SocialFriendListRequest{} }
func (m *SocialFriendListRequest) String() string { return proto.CompactTextString(m) }
func (*SocialFriendListRequest) ProtoMessage()    {}

func (m *SocialFriendListRequest) GetUserId() int64 {
	if m != nil {
		return m.UserId
	}
	return 0
}

func (m *SocialFriendListRequest) GetToken() string {
	if m != nil {
		return m.Token
	}
	return ""
}

type SocialFriendListResponse struct {
	StatusCode int32   `protobuf:"varint,1,opt,name=status_code,json=statusCode,proto3" json:"status_code,omitempty"`
	StatusMsg  string  `protobuf:"bytes,2,opt,name=status_msg,json=statusMsg,proto3" json:"status_msg,omitempty"`
	UserList   []*User `protobuf:"bytes,3,rep,name=user_list,json=userList,proto3" json:"user_list,omitempty"`
}

func (m *SocialFriendListResponse) Reset()         { *m = SocialFriendListResponse{} }
func (m *SocialFriendListResponse) String() string { return proto.CompactTextString(m) }
func (*SocialFriendListResponse) ProtoMessage()    {}

func (m *SocialFriendListResponse) GetStatusCode() int32 {
	if m != nil {
		return m.StatusCode
	}
	return 0
}

func (m *SocialFriendListResponse) GetStatusMsg() string {
	if m != nil {
		return m.StatusMsg
	}
	return ""
}

func (m *SocialFriendListResponse) GetUserList() []*User {
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

type SocialMessagePostRequest struct {
	Token      string `protobuf:"bytes,1,opt,name=token,proto3" json:"token,omitempty"`
	ToUserId   int64  `protobuf:"varint,2,opt,name=to_user_id,proto3" json:"to_user_id,omitempty"`
	ActionType int32  `protobuf:"varint,3,opt,name=action_type,json=actionType,proto3" json:"action_type,omitempty"`
	Content    string `protobuf:"bytes,4,opt,name=content,proto3" json:"content,omitempty"`
}

func (m *SocialMessagePostRequest) Reset()         { *m = SocialMessagePostRequest{} }
func (m *SocialMessagePostRequest) String() string { return proto.CompactTextString(m) }
func (*SocialMessagePostRequest) ProtoMessage()    {}

func (m *SocialMessagePostRequest) GetToken() string {
	if m != nil {
		return m.Token
	}
	return ""
}

func (m *SocialMessagePostRequest) GetToUserId() int64 {
	if m != nil {
		return m.ToUserId
	}
	return 0
}

func (m *SocialMessagePostRequest) GetActionType() int32 {
	if m != nil {
		return m.ActionType
	}
	return 0
}

func (m *SocialMessagePostRequest) GetContent() string {
	if m != nil {
		return m.Content
	}
	return ""
}

type SocialMessagePostResponse struct {
	StatusCode int32  `protobuf:"varint,1,opt,name=status_code,json=statusCode,proto3" json:"status_code,omitempty"`
	StatusMsg  string `protobuf:"bytes,2,opt,name=status_msg,json=statusMsg,proto3" json:"status_msg,omitempty"`
}

func (m *SocialMessagePostResponse) Reset()         { *m = SocialMessagePostResponse{} }
func (m *SocialMessagePostResponse) String() string { return proto.CompactTextString(m) }
func (*SocialMessagePostResponse) ProtoMessage()    {}

func (m *SocialMessagePostResponse) GetStatusCode() int32 {
	if m != nil {
		return m.StatusCode
	}
	return 0
}

func (m *SocialMessagePostResponse) GetStatusMsg() string {
	if m != nil {
		return m.StatusMsg
	}
	return ""
}
