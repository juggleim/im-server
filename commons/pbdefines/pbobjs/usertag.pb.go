// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.33.0
// 	protoc        v3.5.0
// source: usertag.proto

package pbobjs

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type UserIds struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UserIds []string `protobuf:"bytes,1,rep,name=user_ids,json=userIds,proto3" json:"user_ids,omitempty"`
}

func (x *UserIds) Reset() {
	*x = UserIds{}
	if protoimpl.UnsafeEnabled {
		mi := &file_usertag_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UserIds) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UserIds) ProtoMessage() {}

func (x *UserIds) ProtoReflect() protoreflect.Message {
	mi := &file_usertag_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UserIds.ProtoReflect.Descriptor instead.
func (*UserIds) Descriptor() ([]byte, []int) {
	return file_usertag_proto_rawDescGZIP(), []int{0}
}

func (x *UserIds) GetUserIds() []string {
	if x != nil {
		return x.UserIds
	}
	return nil
}

type UserTag struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UserId string   `protobuf:"bytes,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	Tags   []string `protobuf:"bytes,2,rep,name=tags,proto3" json:"tags,omitempty"`
}

func (x *UserTag) Reset() {
	*x = UserTag{}
	if protoimpl.UnsafeEnabled {
		mi := &file_usertag_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UserTag) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UserTag) ProtoMessage() {}

func (x *UserTag) ProtoReflect() protoreflect.Message {
	mi := &file_usertag_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UserTag.ProtoReflect.Descriptor instead.
func (*UserTag) Descriptor() ([]byte, []int) {
	return file_usertag_proto_rawDescGZIP(), []int{1}
}

func (x *UserTag) GetUserId() string {
	if x != nil {
		return x.UserId
	}
	return ""
}

func (x *UserTag) GetTags() []string {
	if x != nil {
		return x.Tags
	}
	return nil
}

type UserTagList struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UserTags []*UserTag `protobuf:"bytes,1,rep,name=user_tags,json=userTags,proto3" json:"user_tags,omitempty"`
}

func (x *UserTagList) Reset() {
	*x = UserTagList{}
	if protoimpl.UnsafeEnabled {
		mi := &file_usertag_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UserTagList) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UserTagList) ProtoMessage() {}

func (x *UserTagList) ProtoReflect() protoreflect.Message {
	mi := &file_usertag_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UserTagList.ProtoReflect.Descriptor instead.
func (*UserTagList) Descriptor() ([]byte, []int) {
	return file_usertag_proto_rawDescGZIP(), []int{2}
}

func (x *UserTagList) GetUserTags() []*UserTag {
	if x != nil {
		return x.UserTags
	}
	return nil
}

type PushNotificationWithTags struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	FromUserId   string                                 `protobuf:"bytes,1,opt,name=from_user_id,json=fromUserId,proto3" json:"from_user_id,omitempty"`
	Condition    *PushNotificationWithTags_Condition    `protobuf:"bytes,2,opt,name=condition,proto3" json:"condition,omitempty"`
	MsgBody      *PushNotificationWithTags_MsgBody      `protobuf:"bytes,3,opt,name=msg_body,json=msgBody,proto3" json:"msg_body,omitempty"`
	Notification *PushNotificationWithTags_Notification `protobuf:"bytes,4,opt,name=notification,proto3" json:"notification,omitempty"`
}

func (x *PushNotificationWithTags) Reset() {
	*x = PushNotificationWithTags{}
	if protoimpl.UnsafeEnabled {
		mi := &file_usertag_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PushNotificationWithTags) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PushNotificationWithTags) ProtoMessage() {}

func (x *PushNotificationWithTags) ProtoReflect() protoreflect.Message {
	mi := &file_usertag_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PushNotificationWithTags.ProtoReflect.Descriptor instead.
func (*PushNotificationWithTags) Descriptor() ([]byte, []int) {
	return file_usertag_proto_rawDescGZIP(), []int{3}
}

func (x *PushNotificationWithTags) GetFromUserId() string {
	if x != nil {
		return x.FromUserId
	}
	return ""
}

func (x *PushNotificationWithTags) GetCondition() *PushNotificationWithTags_Condition {
	if x != nil {
		return x.Condition
	}
	return nil
}

func (x *PushNotificationWithTags) GetMsgBody() *PushNotificationWithTags_MsgBody {
	if x != nil {
		return x.MsgBody
	}
	return nil
}

func (x *PushNotificationWithTags) GetNotification() *PushNotificationWithTags_Notification {
	if x != nil {
		return x.Notification
	}
	return nil
}

type PushNotificationWithTags_MsgBody struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	MsgType    string `protobuf:"bytes,1,opt,name=msg_type,json=msgType,proto3" json:"msg_type,omitempty"`
	MsgContent string `protobuf:"bytes,2,opt,name=msg_content,json=msgContent,proto3" json:"msg_content,omitempty"`
}

func (x *PushNotificationWithTags_MsgBody) Reset() {
	*x = PushNotificationWithTags_MsgBody{}
	if protoimpl.UnsafeEnabled {
		mi := &file_usertag_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PushNotificationWithTags_MsgBody) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PushNotificationWithTags_MsgBody) ProtoMessage() {}

func (x *PushNotificationWithTags_MsgBody) ProtoReflect() protoreflect.Message {
	mi := &file_usertag_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PushNotificationWithTags_MsgBody.ProtoReflect.Descriptor instead.
func (*PushNotificationWithTags_MsgBody) Descriptor() ([]byte, []int) {
	return file_usertag_proto_rawDescGZIP(), []int{3, 0}
}

func (x *PushNotificationWithTags_MsgBody) GetMsgType() string {
	if x != nil {
		return x.MsgType
	}
	return ""
}

func (x *PushNotificationWithTags_MsgBody) GetMsgContent() string {
	if x != nil {
		return x.MsgContent
	}
	return ""
}

type PushNotificationWithTags_Notification struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Title    string `protobuf:"bytes,1,opt,name=title,proto3" json:"title,omitempty"`
	PushText string `protobuf:"bytes,2,opt,name=push_text,json=pushText,proto3" json:"push_text,omitempty"`
}

func (x *PushNotificationWithTags_Notification) Reset() {
	*x = PushNotificationWithTags_Notification{}
	if protoimpl.UnsafeEnabled {
		mi := &file_usertag_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PushNotificationWithTags_Notification) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PushNotificationWithTags_Notification) ProtoMessage() {}

func (x *PushNotificationWithTags_Notification) ProtoReflect() protoreflect.Message {
	mi := &file_usertag_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PushNotificationWithTags_Notification.ProtoReflect.Descriptor instead.
func (*PushNotificationWithTags_Notification) Descriptor() ([]byte, []int) {
	return file_usertag_proto_rawDescGZIP(), []int{3, 1}
}

func (x *PushNotificationWithTags_Notification) GetTitle() string {
	if x != nil {
		return x.Title
	}
	return ""
}

func (x *PushNotificationWithTags_Notification) GetPushText() string {
	if x != nil {
		return x.PushText
	}
	return ""
}

type PushNotificationWithTags_Condition struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	TagsAnd []string `protobuf:"bytes,1,rep,name=tags_and,json=tagsAnd,proto3" json:"tags_and,omitempty"`
	TagsOr  []string `protobuf:"bytes,2,rep,name=tags_or,json=tagsOr,proto3" json:"tags_or,omitempty"`
}

func (x *PushNotificationWithTags_Condition) Reset() {
	*x = PushNotificationWithTags_Condition{}
	if protoimpl.UnsafeEnabled {
		mi := &file_usertag_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PushNotificationWithTags_Condition) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PushNotificationWithTags_Condition) ProtoMessage() {}

func (x *PushNotificationWithTags_Condition) ProtoReflect() protoreflect.Message {
	mi := &file_usertag_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PushNotificationWithTags_Condition.ProtoReflect.Descriptor instead.
func (*PushNotificationWithTags_Condition) Descriptor() ([]byte, []int) {
	return file_usertag_proto_rawDescGZIP(), []int{3, 2}
}

func (x *PushNotificationWithTags_Condition) GetTagsAnd() []string {
	if x != nil {
		return x.TagsAnd
	}
	return nil
}

func (x *PushNotificationWithTags_Condition) GetTagsOr() []string {
	if x != nil {
		return x.TagsOr
	}
	return nil
}

var File_usertag_proto protoreflect.FileDescriptor

var file_usertag_proto_rawDesc = []byte{
	0x0a, 0x0d, 0x75, 0x73, 0x65, 0x72, 0x74, 0x61, 0x67, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22,
	0x24, 0x0a, 0x07, 0x55, 0x73, 0x65, 0x72, 0x49, 0x64, 0x73, 0x12, 0x19, 0x0a, 0x08, 0x75, 0x73,
	0x65, 0x72, 0x5f, 0x69, 0x64, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x09, 0x52, 0x07, 0x75, 0x73,
	0x65, 0x72, 0x49, 0x64, 0x73, 0x22, 0x36, 0x0a, 0x07, 0x55, 0x73, 0x65, 0x72, 0x54, 0x61, 0x67,
	0x12, 0x17, 0x0a, 0x07, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x74, 0x61, 0x67,
	0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x09, 0x52, 0x04, 0x74, 0x61, 0x67, 0x73, 0x22, 0x34, 0x0a,
	0x0b, 0x55, 0x73, 0x65, 0x72, 0x54, 0x61, 0x67, 0x4c, 0x69, 0x73, 0x74, 0x12, 0x25, 0x0a, 0x09,
	0x75, 0x73, 0x65, 0x72, 0x5f, 0x74, 0x61, 0x67, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32,
	0x08, 0x2e, 0x55, 0x73, 0x65, 0x72, 0x54, 0x61, 0x67, 0x52, 0x08, 0x75, 0x73, 0x65, 0x72, 0x54,
	0x61, 0x67, 0x73, 0x22, 0xd4, 0x03, 0x0a, 0x18, 0x50, 0x75, 0x73, 0x68, 0x4e, 0x6f, 0x74, 0x69,
	0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x57, 0x69, 0x74, 0x68, 0x54, 0x61, 0x67, 0x73,
	0x12, 0x20, 0x0a, 0x0c, 0x66, 0x72, 0x6f, 0x6d, 0x5f, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x69, 0x64,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x66, 0x72, 0x6f, 0x6d, 0x55, 0x73, 0x65, 0x72,
	0x49, 0x64, 0x12, 0x41, 0x0a, 0x09, 0x63, 0x6f, 0x6e, 0x64, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x23, 0x2e, 0x50, 0x75, 0x73, 0x68, 0x4e, 0x6f, 0x74, 0x69,
	0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x57, 0x69, 0x74, 0x68, 0x54, 0x61, 0x67, 0x73,
	0x2e, 0x43, 0x6f, 0x6e, 0x64, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x09, 0x63, 0x6f, 0x6e, 0x64,
	0x69, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x3c, 0x0a, 0x08, 0x6d, 0x73, 0x67, 0x5f, 0x62, 0x6f, 0x64,
	0x79, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x21, 0x2e, 0x50, 0x75, 0x73, 0x68, 0x4e, 0x6f,
	0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x57, 0x69, 0x74, 0x68, 0x54, 0x61,
	0x67, 0x73, 0x2e, 0x4d, 0x73, 0x67, 0x42, 0x6f, 0x64, 0x79, 0x52, 0x07, 0x6d, 0x73, 0x67, 0x42,
	0x6f, 0x64, 0x79, 0x12, 0x4a, 0x0a, 0x0c, 0x6e, 0x6f, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x26, 0x2e, 0x50, 0x75, 0x73, 0x68,
	0x4e, 0x6f, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x57, 0x69, 0x74, 0x68,
	0x54, 0x61, 0x67, 0x73, 0x2e, 0x4e, 0x6f, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x52, 0x0c, 0x6e, 0x6f, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x1a,
	0x45, 0x0a, 0x07, 0x4d, 0x73, 0x67, 0x42, 0x6f, 0x64, 0x79, 0x12, 0x19, 0x0a, 0x08, 0x6d, 0x73,
	0x67, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x6d, 0x73,
	0x67, 0x54, 0x79, 0x70, 0x65, 0x12, 0x1f, 0x0a, 0x0b, 0x6d, 0x73, 0x67, 0x5f, 0x63, 0x6f, 0x6e,
	0x74, 0x65, 0x6e, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x6d, 0x73, 0x67, 0x43,
	0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x1a, 0x41, 0x0a, 0x0c, 0x4e, 0x6f, 0x74, 0x69, 0x66, 0x69,
	0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x14, 0x0a, 0x05, 0x74, 0x69, 0x74, 0x6c, 0x65, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x74, 0x69, 0x74, 0x6c, 0x65, 0x12, 0x1b, 0x0a, 0x09,
	0x70, 0x75, 0x73, 0x68, 0x5f, 0x74, 0x65, 0x78, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x08, 0x70, 0x75, 0x73, 0x68, 0x54, 0x65, 0x78, 0x74, 0x1a, 0x3f, 0x0a, 0x09, 0x43, 0x6f, 0x6e,
	0x64, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x19, 0x0a, 0x08, 0x74, 0x61, 0x67, 0x73, 0x5f, 0x61,
	0x6e, 0x64, 0x18, 0x01, 0x20, 0x03, 0x28, 0x09, 0x52, 0x07, 0x74, 0x61, 0x67, 0x73, 0x41, 0x6e,
	0x64, 0x12, 0x17, 0x0a, 0x07, 0x74, 0x61, 0x67, 0x73, 0x5f, 0x6f, 0x72, 0x18, 0x02, 0x20, 0x03,
	0x28, 0x09, 0x52, 0x06, 0x74, 0x61, 0x67, 0x73, 0x4f, 0x72, 0x42, 0x0a, 0x5a, 0x08, 0x2e, 0x2f,
	0x70, 0x62, 0x6f, 0x62, 0x6a, 0x73, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_usertag_proto_rawDescOnce sync.Once
	file_usertag_proto_rawDescData = file_usertag_proto_rawDesc
)

func file_usertag_proto_rawDescGZIP() []byte {
	file_usertag_proto_rawDescOnce.Do(func() {
		file_usertag_proto_rawDescData = protoimpl.X.CompressGZIP(file_usertag_proto_rawDescData)
	})
	return file_usertag_proto_rawDescData
}

var file_usertag_proto_msgTypes = make([]protoimpl.MessageInfo, 7)
var file_usertag_proto_goTypes = []interface{}{
	(*UserIds)(nil),                               // 0: UserIds
	(*UserTag)(nil),                               // 1: UserTag
	(*UserTagList)(nil),                           // 2: UserTagList
	(*PushNotificationWithTags)(nil),              // 3: PushNotificationWithTags
	(*PushNotificationWithTags_MsgBody)(nil),      // 4: PushNotificationWithTags.MsgBody
	(*PushNotificationWithTags_Notification)(nil), // 5: PushNotificationWithTags.Notification
	(*PushNotificationWithTags_Condition)(nil),    // 6: PushNotificationWithTags.Condition
}
var file_usertag_proto_depIdxs = []int32{
	1, // 0: UserTagList.user_tags:type_name -> UserTag
	6, // 1: PushNotificationWithTags.condition:type_name -> PushNotificationWithTags.Condition
	4, // 2: PushNotificationWithTags.msg_body:type_name -> PushNotificationWithTags.MsgBody
	5, // 3: PushNotificationWithTags.notification:type_name -> PushNotificationWithTags.Notification
	4, // [4:4] is the sub-list for method output_type
	4, // [4:4] is the sub-list for method input_type
	4, // [4:4] is the sub-list for extension type_name
	4, // [4:4] is the sub-list for extension extendee
	0, // [0:4] is the sub-list for field type_name
}

func init() { file_usertag_proto_init() }
func file_usertag_proto_init() {
	if File_usertag_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_usertag_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UserIds); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_usertag_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UserTag); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_usertag_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UserTagList); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_usertag_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PushNotificationWithTags); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_usertag_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PushNotificationWithTags_MsgBody); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_usertag_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PushNotificationWithTags_Notification); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_usertag_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PushNotificationWithTags_Condition); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_usertag_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   7,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_usertag_proto_goTypes,
		DependencyIndexes: file_usertag_proto_depIdxs,
		MessageInfos:      file_usertag_proto_msgTypes,
	}.Build()
	File_usertag_proto = out.File
	file_usertag_proto_rawDesc = nil
	file_usertag_proto_goTypes = nil
	file_usertag_proto_depIdxs = nil
}
