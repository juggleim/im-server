// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.33.0
// 	protoc        v3.5.0
// source: commons/pbdefines/rtcroom.proto

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

type RtcRoomType int32

const (
	RtcRoomType_OneOne  RtcRoomType = 0
	RtcRoomType_OneMore RtcRoomType = 1
)

// Enum value maps for RtcRoomType.
var (
	RtcRoomType_name = map[int32]string{
		0: "OneOne",
		1: "OneMore",
	}
	RtcRoomType_value = map[string]int32{
		"OneOne":  0,
		"OneMore": 1,
	}
)

func (x RtcRoomType) Enum() *RtcRoomType {
	p := new(RtcRoomType)
	*p = x
	return p
}

func (x RtcRoomType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (RtcRoomType) Descriptor() protoreflect.EnumDescriptor {
	return file_commons_pbdefines_rtcroom_proto_enumTypes[0].Descriptor()
}

func (RtcRoomType) Type() protoreflect.EnumType {
	return &file_commons_pbdefines_rtcroom_proto_enumTypes[0]
}

func (x RtcRoomType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use RtcRoomType.Descriptor instead.
func (RtcRoomType) EnumDescriptor() ([]byte, []int) {
	return file_commons_pbdefines_rtcroom_proto_rawDescGZIP(), []int{0}
}

type RtcState int32

const (
	RtcState_RtcStateDefault RtcState = 0
	RtcState_RtcIncoming     RtcState = 1
	RtcState_RtcOutgoing     RtcState = 2
	RtcState_RtcConnecting   RtcState = 3
	RtcState_RtcConnected    RtcState = 4
)

// Enum value maps for RtcState.
var (
	RtcState_name = map[int32]string{
		0: "RtcStateDefault",
		1: "RtcIncoming",
		2: "RtcOutgoing",
		3: "RtcConnecting",
		4: "RtcConnected",
	}
	RtcState_value = map[string]int32{
		"RtcStateDefault": 0,
		"RtcIncoming":     1,
		"RtcOutgoing":     2,
		"RtcConnecting":   3,
		"RtcConnected":    4,
	}
)

func (x RtcState) Enum() *RtcState {
	p := new(RtcState)
	*p = x
	return p
}

func (x RtcState) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (RtcState) Descriptor() protoreflect.EnumDescriptor {
	return file_commons_pbdefines_rtcroom_proto_enumTypes[1].Descriptor()
}

func (RtcState) Type() protoreflect.EnumType {
	return &file_commons_pbdefines_rtcroom_proto_enumTypes[1]
}

func (x RtcState) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use RtcState.Descriptor instead.
func (RtcState) EnumDescriptor() ([]byte, []int) {
	return file_commons_pbdefines_rtcroom_proto_rawDescGZIP(), []int{1}
}

type RtcRoomEventType int32

const (
	RtcRoomEventType_DefaultRtcRoomEvent RtcRoomEventType = 0
	RtcRoomEventType_RtcJoin             RtcRoomEventType = 1
	RtcRoomEventType_RtcQuit             RtcRoomEventType = 2
	RtcRoomEventType_RtcDestroy          RtcRoomEventType = 3
	RtcRoomEventType_RtcStateChg         RtcRoomEventType = 4
)

// Enum value maps for RtcRoomEventType.
var (
	RtcRoomEventType_name = map[int32]string{
		0: "DefaultRtcRoomEvent",
		1: "RtcJoin",
		2: "RtcQuit",
		3: "RtcDestroy",
		4: "RtcStateChg",
	}
	RtcRoomEventType_value = map[string]int32{
		"DefaultRtcRoomEvent": 0,
		"RtcJoin":             1,
		"RtcQuit":             2,
		"RtcDestroy":          3,
		"RtcStateChg":         4,
	}
)

func (x RtcRoomEventType) Enum() *RtcRoomEventType {
	p := new(RtcRoomEventType)
	*p = x
	return p
}

func (x RtcRoomEventType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (RtcRoomEventType) Descriptor() protoreflect.EnumDescriptor {
	return file_commons_pbdefines_rtcroom_proto_enumTypes[2].Descriptor()
}

func (RtcRoomEventType) Type() protoreflect.EnumType {
	return &file_commons_pbdefines_rtcroom_proto_enumTypes[2]
}

func (x RtcRoomEventType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use RtcRoomEventType.Descriptor instead.
func (RtcRoomEventType) EnumDescriptor() ([]byte, []int) {
	return file_commons_pbdefines_rtcroom_proto_rawDescGZIP(), []int{2}
}

type RtcRoomQuitReason int32

const (
	RtcRoomQuitReason_Active      RtcRoomQuitReason = 0
	RtcRoomQuitReason_CallTimeout RtcRoomQuitReason = 1
	RtcRoomQuitReason_PingTimeout RtcRoomQuitReason = 2
)

// Enum value maps for RtcRoomQuitReason.
var (
	RtcRoomQuitReason_name = map[int32]string{
		0: "Active",
		1: "CallTimeout",
		2: "PingTimeout",
	}
	RtcRoomQuitReason_value = map[string]int32{
		"Active":      0,
		"CallTimeout": 1,
		"PingTimeout": 2,
	}
)

func (x RtcRoomQuitReason) Enum() *RtcRoomQuitReason {
	p := new(RtcRoomQuitReason)
	*p = x
	return p
}

func (x RtcRoomQuitReason) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (RtcRoomQuitReason) Descriptor() protoreflect.EnumDescriptor {
	return file_commons_pbdefines_rtcroom_proto_enumTypes[3].Descriptor()
}

func (RtcRoomQuitReason) Type() protoreflect.EnumType {
	return &file_commons_pbdefines_rtcroom_proto_enumTypes[3]
}

func (x RtcRoomQuitReason) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use RtcRoomQuitReason.Descriptor instead.
func (RtcRoomQuitReason) EnumDescriptor() ([]byte, []int) {
	return file_commons_pbdefines_rtcroom_proto_rawDescGZIP(), []int{3}
}

type RtcRoomReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	RoomType   RtcRoomType `protobuf:"varint,1,opt,name=roomType,proto3,enum=RtcRoomType" json:"roomType,omitempty"`
	RoomId     string      `protobuf:"bytes,2,opt,name=roomId,proto3" json:"roomId,omitempty"`
	JoinMember *RtcMember  `protobuf:"bytes,3,opt,name=joinMember,proto3" json:"joinMember,omitempty"`
}

func (x *RtcRoomReq) Reset() {
	*x = RtcRoomReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_commons_pbdefines_rtcroom_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RtcRoomReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RtcRoomReq) ProtoMessage() {}

func (x *RtcRoomReq) ProtoReflect() protoreflect.Message {
	mi := &file_commons_pbdefines_rtcroom_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RtcRoomReq.ProtoReflect.Descriptor instead.
func (*RtcRoomReq) Descriptor() ([]byte, []int) {
	return file_commons_pbdefines_rtcroom_proto_rawDescGZIP(), []int{0}
}

func (x *RtcRoomReq) GetRoomType() RtcRoomType {
	if x != nil {
		return x.RoomType
	}
	return RtcRoomType_OneOne
}

func (x *RtcRoomReq) GetRoomId() string {
	if x != nil {
		return x.RoomId
	}
	return ""
}

func (x *RtcRoomReq) GetJoinMember() *RtcMember {
	if x != nil {
		return x.JoinMember
	}
	return nil
}

type RtcRoom struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	RoomType RtcRoomType  `protobuf:"varint,1,opt,name=roomType,proto3,enum=RtcRoomType" json:"roomType,omitempty"`
	RoomId   string       `protobuf:"bytes,2,opt,name=roomId,proto3" json:"roomId,omitempty"`
	Owner    *UserInfo    `protobuf:"bytes,3,opt,name=owner,proto3" json:"owner,omitempty"`
	Members  []*RtcMember `protobuf:"bytes,51,rep,name=members,proto3" json:"members,omitempty"`
}

func (x *RtcRoom) Reset() {
	*x = RtcRoom{}
	if protoimpl.UnsafeEnabled {
		mi := &file_commons_pbdefines_rtcroom_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RtcRoom) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RtcRoom) ProtoMessage() {}

func (x *RtcRoom) ProtoReflect() protoreflect.Message {
	mi := &file_commons_pbdefines_rtcroom_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RtcRoom.ProtoReflect.Descriptor instead.
func (*RtcRoom) Descriptor() ([]byte, []int) {
	return file_commons_pbdefines_rtcroom_proto_rawDescGZIP(), []int{1}
}

func (x *RtcRoom) GetRoomType() RtcRoomType {
	if x != nil {
		return x.RoomType
	}
	return RtcRoomType_OneOne
}

func (x *RtcRoom) GetRoomId() string {
	if x != nil {
		return x.RoomId
	}
	return ""
}

func (x *RtcRoom) GetOwner() *UserInfo {
	if x != nil {
		return x.Owner
	}
	return nil
}

func (x *RtcRoom) GetMembers() []*RtcMember {
	if x != nil {
		return x.Members
	}
	return nil
}

type RtcMember struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Member      *UserInfo `protobuf:"bytes,1,opt,name=member,proto3" json:"member,omitempty"`
	RtcState    RtcState  `protobuf:"varint,2,opt,name=rtcState,proto3,enum=RtcState" json:"rtcState,omitempty"`
	CallTime    int64     `protobuf:"varint,3,opt,name=callTime,proto3" json:"callTime,omitempty"`
	ConnectTime int64     `protobuf:"varint,4,opt,name=connectTime,proto3" json:"connectTime,omitempty"`
	HangupTime  int64     `protobuf:"varint,5,opt,name=hangupTime,proto3" json:"hangupTime,omitempty"`
	Inviter     *UserInfo `protobuf:"bytes,6,opt,name=inviter,proto3" json:"inviter,omitempty"`
}

func (x *RtcMember) Reset() {
	*x = RtcMember{}
	if protoimpl.UnsafeEnabled {
		mi := &file_commons_pbdefines_rtcroom_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RtcMember) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RtcMember) ProtoMessage() {}

func (x *RtcMember) ProtoReflect() protoreflect.Message {
	mi := &file_commons_pbdefines_rtcroom_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RtcMember.ProtoReflect.Descriptor instead.
func (*RtcMember) Descriptor() ([]byte, []int) {
	return file_commons_pbdefines_rtcroom_proto_rawDescGZIP(), []int{2}
}

func (x *RtcMember) GetMember() *UserInfo {
	if x != nil {
		return x.Member
	}
	return nil
}

func (x *RtcMember) GetRtcState() RtcState {
	if x != nil {
		return x.RtcState
	}
	return RtcState_RtcStateDefault
}

func (x *RtcMember) GetCallTime() int64 {
	if x != nil {
		return x.CallTime
	}
	return 0
}

func (x *RtcMember) GetConnectTime() int64 {
	if x != nil {
		return x.ConnectTime
	}
	return 0
}

func (x *RtcMember) GetHangupTime() int64 {
	if x != nil {
		return x.HangupTime
	}
	return 0
}

func (x *RtcMember) GetInviter() *UserInfo {
	if x != nil {
		return x.Inviter
	}
	return nil
}

type MemberState struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	RoomId   string      `protobuf:"bytes,1,opt,name=roomId,proto3" json:"roomId,omitempty"`
	RoomType RtcRoomType `protobuf:"varint,2,opt,name=roomType,proto3,enum=RtcRoomType" json:"roomType,omitempty"`
	MemberId string      `protobuf:"bytes,3,opt,name=memberId,proto3" json:"memberId,omitempty"`
	DeviceId string      `protobuf:"bytes,4,opt,name=deviceId,proto3" json:"deviceId,omitempty"`
	RtcState RtcState    `protobuf:"varint,5,opt,name=rtcState,proto3,enum=RtcState" json:"rtcState,omitempty"`
}

func (x *MemberState) Reset() {
	*x = MemberState{}
	if protoimpl.UnsafeEnabled {
		mi := &file_commons_pbdefines_rtcroom_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MemberState) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MemberState) ProtoMessage() {}

func (x *MemberState) ProtoReflect() protoreflect.Message {
	mi := &file_commons_pbdefines_rtcroom_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MemberState.ProtoReflect.Descriptor instead.
func (*MemberState) Descriptor() ([]byte, []int) {
	return file_commons_pbdefines_rtcroom_proto_rawDescGZIP(), []int{3}
}

func (x *MemberState) GetRoomId() string {
	if x != nil {
		return x.RoomId
	}
	return ""
}

func (x *MemberState) GetRoomType() RtcRoomType {
	if x != nil {
		return x.RoomType
	}
	return RtcRoomType_OneOne
}

func (x *MemberState) GetMemberId() string {
	if x != nil {
		return x.MemberId
	}
	return ""
}

func (x *MemberState) GetDeviceId() string {
	if x != nil {
		return x.DeviceId
	}
	return ""
}

func (x *MemberState) GetRtcState() RtcState {
	if x != nil {
		return x.RtcState
	}
	return RtcState_RtcStateDefault
}

type SyncMemberStateReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	IsDelete bool         `protobuf:"varint,1,opt,name=isDelete,proto3" json:"isDelete,omitempty"`
	Member   *MemberState `protobuf:"bytes,2,opt,name=member,proto3" json:"member,omitempty"`
}

func (x *SyncMemberStateReq) Reset() {
	*x = SyncMemberStateReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_commons_pbdefines_rtcroom_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SyncMemberStateReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SyncMemberStateReq) ProtoMessage() {}

func (x *SyncMemberStateReq) ProtoReflect() protoreflect.Message {
	mi := &file_commons_pbdefines_rtcroom_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SyncMemberStateReq.ProtoReflect.Descriptor instead.
func (*SyncMemberStateReq) Descriptor() ([]byte, []int) {
	return file_commons_pbdefines_rtcroom_proto_rawDescGZIP(), []int{4}
}

func (x *SyncMemberStateReq) GetIsDelete() bool {
	if x != nil {
		return x.IsDelete
	}
	return false
}

func (x *SyncMemberStateReq) GetMember() *MemberState {
	if x != nil {
		return x.Member
	}
	return nil
}

type RtcRoomEvent struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	RoomEventType RtcRoomEventType  `protobuf:"varint,1,opt,name=roomEventType,proto3,enum=RtcRoomEventType" json:"roomEventType,omitempty"`
	Member        *RtcMember        `protobuf:"bytes,2,opt,name=member,proto3" json:"member,omitempty"`
	Room          *RtcRoom          `protobuf:"bytes,3,opt,name=room,proto3" json:"room,omitempty"`
	Reason        RtcRoomQuitReason `protobuf:"varint,4,opt,name=reason,proto3,enum=RtcRoomQuitReason" json:"reason,omitempty"`
}

func (x *RtcRoomEvent) Reset() {
	*x = RtcRoomEvent{}
	if protoimpl.UnsafeEnabled {
		mi := &file_commons_pbdefines_rtcroom_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RtcRoomEvent) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RtcRoomEvent) ProtoMessage() {}

func (x *RtcRoomEvent) ProtoReflect() protoreflect.Message {
	mi := &file_commons_pbdefines_rtcroom_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RtcRoomEvent.ProtoReflect.Descriptor instead.
func (*RtcRoomEvent) Descriptor() ([]byte, []int) {
	return file_commons_pbdefines_rtcroom_proto_rawDescGZIP(), []int{5}
}

func (x *RtcRoomEvent) GetRoomEventType() RtcRoomEventType {
	if x != nil {
		return x.RoomEventType
	}
	return RtcRoomEventType_DefaultRtcRoomEvent
}

func (x *RtcRoomEvent) GetMember() *RtcMember {
	if x != nil {
		return x.Member
	}
	return nil
}

func (x *RtcRoomEvent) GetRoom() *RtcRoom {
	if x != nil {
		return x.Room
	}
	return nil
}

func (x *RtcRoomEvent) GetReason() RtcRoomQuitReason {
	if x != nil {
		return x.Reason
	}
	return RtcRoomQuitReason_Active
}

type RtcInviteReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	TargetIds []string    `protobuf:"bytes,1,rep,name=targetIds,proto3" json:"targetIds,omitempty"`
	RoomType  RtcRoomType `protobuf:"varint,2,opt,name=roomType,proto3,enum=RtcRoomType" json:"roomType,omitempty"`
	RoomId    string      `protobuf:"bytes,3,opt,name=roomId,proto3" json:"roomId,omitempty"`
}

func (x *RtcInviteReq) Reset() {
	*x = RtcInviteReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_commons_pbdefines_rtcroom_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RtcInviteReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RtcInviteReq) ProtoMessage() {}

func (x *RtcInviteReq) ProtoReflect() protoreflect.Message {
	mi := &file_commons_pbdefines_rtcroom_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RtcInviteReq.ProtoReflect.Descriptor instead.
func (*RtcInviteReq) Descriptor() ([]byte, []int) {
	return file_commons_pbdefines_rtcroom_proto_rawDescGZIP(), []int{6}
}

func (x *RtcInviteReq) GetTargetIds() []string {
	if x != nil {
		return x.TargetIds
	}
	return nil
}

func (x *RtcInviteReq) GetRoomType() RtcRoomType {
	if x != nil {
		return x.RoomType
	}
	return RtcRoomType_OneOne
}

func (x *RtcInviteReq) GetRoomId() string {
	if x != nil {
		return x.RoomId
	}
	return ""
}

type RtcMemberRooms struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Rooms []*RtcMemberRoom `protobuf:"bytes,1,rep,name=rooms,proto3" json:"rooms,omitempty"`
}

func (x *RtcMemberRooms) Reset() {
	*x = RtcMemberRooms{}
	if protoimpl.UnsafeEnabled {
		mi := &file_commons_pbdefines_rtcroom_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RtcMemberRooms) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RtcMemberRooms) ProtoMessage() {}

func (x *RtcMemberRooms) ProtoReflect() protoreflect.Message {
	mi := &file_commons_pbdefines_rtcroom_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RtcMemberRooms.ProtoReflect.Descriptor instead.
func (*RtcMemberRooms) Descriptor() ([]byte, []int) {
	return file_commons_pbdefines_rtcroom_proto_rawDescGZIP(), []int{7}
}

func (x *RtcMemberRooms) GetRooms() []*RtcMemberRoom {
	if x != nil {
		return x.Rooms
	}
	return nil
}

type RtcMemberRoom struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	RoomType RtcRoomType `protobuf:"varint,1,opt,name=roomType,proto3,enum=RtcRoomType" json:"roomType,omitempty"`
	RoomId   string      `protobuf:"bytes,2,opt,name=roomId,proto3" json:"roomId,omitempty"`
	Owner    *UserInfo   `protobuf:"bytes,3,opt,name=owner,proto3" json:"owner,omitempty"`
	RtcState RtcState    `protobuf:"varint,4,opt,name=rtcState,proto3,enum=RtcState" json:"rtcState,omitempty"`
}

func (x *RtcMemberRoom) Reset() {
	*x = RtcMemberRoom{}
	if protoimpl.UnsafeEnabled {
		mi := &file_commons_pbdefines_rtcroom_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RtcMemberRoom) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RtcMemberRoom) ProtoMessage() {}

func (x *RtcMemberRoom) ProtoReflect() protoreflect.Message {
	mi := &file_commons_pbdefines_rtcroom_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RtcMemberRoom.ProtoReflect.Descriptor instead.
func (*RtcMemberRoom) Descriptor() ([]byte, []int) {
	return file_commons_pbdefines_rtcroom_proto_rawDescGZIP(), []int{8}
}

func (x *RtcMemberRoom) GetRoomType() RtcRoomType {
	if x != nil {
		return x.RoomType
	}
	return RtcRoomType_OneOne
}

func (x *RtcMemberRoom) GetRoomId() string {
	if x != nil {
		return x.RoomId
	}
	return ""
}

func (x *RtcMemberRoom) GetOwner() *UserInfo {
	if x != nil {
		return x.Owner
	}
	return nil
}

func (x *RtcMemberRoom) GetRtcState() RtcState {
	if x != nil {
		return x.RtcState
	}
	return RtcState_RtcStateDefault
}

var File_commons_pbdefines_rtcroom_proto protoreflect.FileDescriptor

var file_commons_pbdefines_rtcroom_proto_rawDesc = []byte{
	0x0a, 0x1f, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x73, 0x2f, 0x70, 0x62, 0x64, 0x65, 0x66, 0x69,
	0x6e, 0x65, 0x73, 0x2f, 0x72, 0x74, 0x63, 0x72, 0x6f, 0x6f, 0x6d, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x1a, 0x23, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x73, 0x2f, 0x70, 0x62, 0x64, 0x65, 0x66,
	0x69, 0x6e, 0x65, 0x73, 0x2f, 0x61, 0x70, 0x70, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x73,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x7a, 0x0a, 0x0a, 0x52, 0x74, 0x63, 0x52, 0x6f, 0x6f,
	0x6d, 0x52, 0x65, 0x71, 0x12, 0x28, 0x0a, 0x08, 0x72, 0x6f, 0x6f, 0x6d, 0x54, 0x79, 0x70, 0x65,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x0c, 0x2e, 0x52, 0x74, 0x63, 0x52, 0x6f, 0x6f, 0x6d,
	0x54, 0x79, 0x70, 0x65, 0x52, 0x08, 0x72, 0x6f, 0x6f, 0x6d, 0x54, 0x79, 0x70, 0x65, 0x12, 0x16,
	0x0a, 0x06, 0x72, 0x6f, 0x6f, 0x6d, 0x49, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06,
	0x72, 0x6f, 0x6f, 0x6d, 0x49, 0x64, 0x12, 0x2a, 0x0a, 0x0a, 0x6a, 0x6f, 0x69, 0x6e, 0x4d, 0x65,
	0x6d, 0x62, 0x65, 0x72, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0a, 0x2e, 0x52, 0x74, 0x63,
	0x4d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x52, 0x0a, 0x6a, 0x6f, 0x69, 0x6e, 0x4d, 0x65, 0x6d, 0x62,
	0x65, 0x72, 0x22, 0x92, 0x01, 0x0a, 0x07, 0x52, 0x74, 0x63, 0x52, 0x6f, 0x6f, 0x6d, 0x12, 0x28,
	0x0a, 0x08, 0x72, 0x6f, 0x6f, 0x6d, 0x54, 0x79, 0x70, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e,
	0x32, 0x0c, 0x2e, 0x52, 0x74, 0x63, 0x52, 0x6f, 0x6f, 0x6d, 0x54, 0x79, 0x70, 0x65, 0x52, 0x08,
	0x72, 0x6f, 0x6f, 0x6d, 0x54, 0x79, 0x70, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x72, 0x6f, 0x6f, 0x6d,
	0x49, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x72, 0x6f, 0x6f, 0x6d, 0x49, 0x64,
	0x12, 0x1f, 0x0a, 0x05, 0x6f, 0x77, 0x6e, 0x65, 0x72, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x09, 0x2e, 0x55, 0x73, 0x65, 0x72, 0x49, 0x6e, 0x66, 0x6f, 0x52, 0x05, 0x6f, 0x77, 0x6e, 0x65,
	0x72, 0x12, 0x24, 0x0a, 0x07, 0x6d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x73, 0x18, 0x33, 0x20, 0x03,
	0x28, 0x0b, 0x32, 0x0a, 0x2e, 0x52, 0x74, 0x63, 0x4d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x52, 0x07,
	0x6d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x73, 0x22, 0xd8, 0x01, 0x0a, 0x09, 0x52, 0x74, 0x63, 0x4d,
	0x65, 0x6d, 0x62, 0x65, 0x72, 0x12, 0x21, 0x0a, 0x06, 0x6d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x09, 0x2e, 0x55, 0x73, 0x65, 0x72, 0x49, 0x6e, 0x66, 0x6f,
	0x52, 0x06, 0x6d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x12, 0x25, 0x0a, 0x08, 0x72, 0x74, 0x63, 0x53,
	0x74, 0x61, 0x74, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x09, 0x2e, 0x52, 0x74, 0x63,
	0x53, 0x74, 0x61, 0x74, 0x65, 0x52, 0x08, 0x72, 0x74, 0x63, 0x53, 0x74, 0x61, 0x74, 0x65, 0x12,
	0x1a, 0x0a, 0x08, 0x63, 0x61, 0x6c, 0x6c, 0x54, 0x69, 0x6d, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28,
	0x03, 0x52, 0x08, 0x63, 0x61, 0x6c, 0x6c, 0x54, 0x69, 0x6d, 0x65, 0x12, 0x20, 0x0a, 0x0b, 0x63,
	0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x54, 0x69, 0x6d, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x03,
	0x52, 0x0b, 0x63, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x54, 0x69, 0x6d, 0x65, 0x12, 0x1e, 0x0a,
	0x0a, 0x68, 0x61, 0x6e, 0x67, 0x75, 0x70, 0x54, 0x69, 0x6d, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28,
	0x03, 0x52, 0x0a, 0x68, 0x61, 0x6e, 0x67, 0x75, 0x70, 0x54, 0x69, 0x6d, 0x65, 0x12, 0x23, 0x0a,
	0x07, 0x69, 0x6e, 0x76, 0x69, 0x74, 0x65, 0x72, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x09,
	0x2e, 0x55, 0x73, 0x65, 0x72, 0x49, 0x6e, 0x66, 0x6f, 0x52, 0x07, 0x69, 0x6e, 0x76, 0x69, 0x74,
	0x65, 0x72, 0x22, 0xae, 0x01, 0x0a, 0x0b, 0x4d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x53, 0x74, 0x61,
	0x74, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x72, 0x6f, 0x6f, 0x6d, 0x49, 0x64, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x06, 0x72, 0x6f, 0x6f, 0x6d, 0x49, 0x64, 0x12, 0x28, 0x0a, 0x08, 0x72, 0x6f,
	0x6f, 0x6d, 0x54, 0x79, 0x70, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x0c, 0x2e, 0x52,
	0x74, 0x63, 0x52, 0x6f, 0x6f, 0x6d, 0x54, 0x79, 0x70, 0x65, 0x52, 0x08, 0x72, 0x6f, 0x6f, 0x6d,
	0x54, 0x79, 0x70, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x6d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x49, 0x64,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x6d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x49, 0x64,
	0x12, 0x1a, 0x0a, 0x08, 0x64, 0x65, 0x76, 0x69, 0x63, 0x65, 0x49, 0x64, 0x18, 0x04, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x08, 0x64, 0x65, 0x76, 0x69, 0x63, 0x65, 0x49, 0x64, 0x12, 0x25, 0x0a, 0x08,
	0x72, 0x74, 0x63, 0x53, 0x74, 0x61, 0x74, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x09,
	0x2e, 0x52, 0x74, 0x63, 0x53, 0x74, 0x61, 0x74, 0x65, 0x52, 0x08, 0x72, 0x74, 0x63, 0x53, 0x74,
	0x61, 0x74, 0x65, 0x22, 0x56, 0x0a, 0x12, 0x53, 0x79, 0x6e, 0x63, 0x4d, 0x65, 0x6d, 0x62, 0x65,
	0x72, 0x53, 0x74, 0x61, 0x74, 0x65, 0x52, 0x65, 0x71, 0x12, 0x1a, 0x0a, 0x08, 0x69, 0x73, 0x44,
	0x65, 0x6c, 0x65, 0x74, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x08, 0x69, 0x73, 0x44,
	0x65, 0x6c, 0x65, 0x74, 0x65, 0x12, 0x24, 0x0a, 0x06, 0x6d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0c, 0x2e, 0x4d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x53, 0x74,
	0x61, 0x74, 0x65, 0x52, 0x06, 0x6d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x22, 0xb5, 0x01, 0x0a, 0x0c,
	0x52, 0x74, 0x63, 0x52, 0x6f, 0x6f, 0x6d, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x12, 0x37, 0x0a, 0x0d,
	0x72, 0x6f, 0x6f, 0x6d, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x54, 0x79, 0x70, 0x65, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x0e, 0x32, 0x11, 0x2e, 0x52, 0x74, 0x63, 0x52, 0x6f, 0x6f, 0x6d, 0x45, 0x76, 0x65,
	0x6e, 0x74, 0x54, 0x79, 0x70, 0x65, 0x52, 0x0d, 0x72, 0x6f, 0x6f, 0x6d, 0x45, 0x76, 0x65, 0x6e,
	0x74, 0x54, 0x79, 0x70, 0x65, 0x12, 0x22, 0x0a, 0x06, 0x6d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0a, 0x2e, 0x52, 0x74, 0x63, 0x4d, 0x65, 0x6d, 0x62, 0x65,
	0x72, 0x52, 0x06, 0x6d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x12, 0x1c, 0x0a, 0x04, 0x72, 0x6f, 0x6f,
	0x6d, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x08, 0x2e, 0x52, 0x74, 0x63, 0x52, 0x6f, 0x6f,
	0x6d, 0x52, 0x04, 0x72, 0x6f, 0x6f, 0x6d, 0x12, 0x2a, 0x0a, 0x06, 0x72, 0x65, 0x61, 0x73, 0x6f,
	0x6e, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x12, 0x2e, 0x52, 0x74, 0x63, 0x52, 0x6f, 0x6f,
	0x6d, 0x51, 0x75, 0x69, 0x74, 0x52, 0x65, 0x61, 0x73, 0x6f, 0x6e, 0x52, 0x06, 0x72, 0x65, 0x61,
	0x73, 0x6f, 0x6e, 0x22, 0x6e, 0x0a, 0x0c, 0x52, 0x74, 0x63, 0x49, 0x6e, 0x76, 0x69, 0x74, 0x65,
	0x52, 0x65, 0x71, 0x12, 0x1c, 0x0a, 0x09, 0x74, 0x61, 0x72, 0x67, 0x65, 0x74, 0x49, 0x64, 0x73,
	0x18, 0x01, 0x20, 0x03, 0x28, 0x09, 0x52, 0x09, 0x74, 0x61, 0x72, 0x67, 0x65, 0x74, 0x49, 0x64,
	0x73, 0x12, 0x28, 0x0a, 0x08, 0x72, 0x6f, 0x6f, 0x6d, 0x54, 0x79, 0x70, 0x65, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x0e, 0x32, 0x0c, 0x2e, 0x52, 0x74, 0x63, 0x52, 0x6f, 0x6f, 0x6d, 0x54, 0x79, 0x70,
	0x65, 0x52, 0x08, 0x72, 0x6f, 0x6f, 0x6d, 0x54, 0x79, 0x70, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x72,
	0x6f, 0x6f, 0x6d, 0x49, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x72, 0x6f, 0x6f,
	0x6d, 0x49, 0x64, 0x22, 0x36, 0x0a, 0x0e, 0x52, 0x74, 0x63, 0x4d, 0x65, 0x6d, 0x62, 0x65, 0x72,
	0x52, 0x6f, 0x6f, 0x6d, 0x73, 0x12, 0x24, 0x0a, 0x05, 0x72, 0x6f, 0x6f, 0x6d, 0x73, 0x18, 0x01,
	0x20, 0x03, 0x28, 0x0b, 0x32, 0x0e, 0x2e, 0x52, 0x74, 0x63, 0x4d, 0x65, 0x6d, 0x62, 0x65, 0x72,
	0x52, 0x6f, 0x6f, 0x6d, 0x52, 0x05, 0x72, 0x6f, 0x6f, 0x6d, 0x73, 0x22, 0x99, 0x01, 0x0a, 0x0d,
	0x52, 0x74, 0x63, 0x4d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x52, 0x6f, 0x6f, 0x6d, 0x12, 0x28, 0x0a,
	0x08, 0x72, 0x6f, 0x6f, 0x6d, 0x54, 0x79, 0x70, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32,
	0x0c, 0x2e, 0x52, 0x74, 0x63, 0x52, 0x6f, 0x6f, 0x6d, 0x54, 0x79, 0x70, 0x65, 0x52, 0x08, 0x72,
	0x6f, 0x6f, 0x6d, 0x54, 0x79, 0x70, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x72, 0x6f, 0x6f, 0x6d, 0x49,
	0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x72, 0x6f, 0x6f, 0x6d, 0x49, 0x64, 0x12,
	0x1f, 0x0a, 0x05, 0x6f, 0x77, 0x6e, 0x65, 0x72, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x09,
	0x2e, 0x55, 0x73, 0x65, 0x72, 0x49, 0x6e, 0x66, 0x6f, 0x52, 0x05, 0x6f, 0x77, 0x6e, 0x65, 0x72,
	0x12, 0x25, 0x0a, 0x08, 0x72, 0x74, 0x63, 0x53, 0x74, 0x61, 0x74, 0x65, 0x18, 0x04, 0x20, 0x01,
	0x28, 0x0e, 0x32, 0x09, 0x2e, 0x52, 0x74, 0x63, 0x53, 0x74, 0x61, 0x74, 0x65, 0x52, 0x08, 0x72,
	0x74, 0x63, 0x53, 0x74, 0x61, 0x74, 0x65, 0x2a, 0x26, 0x0a, 0x0b, 0x52, 0x74, 0x63, 0x52, 0x6f,
	0x6f, 0x6d, 0x54, 0x79, 0x70, 0x65, 0x12, 0x0a, 0x0a, 0x06, 0x4f, 0x6e, 0x65, 0x4f, 0x6e, 0x65,
	0x10, 0x00, 0x12, 0x0b, 0x0a, 0x07, 0x4f, 0x6e, 0x65, 0x4d, 0x6f, 0x72, 0x65, 0x10, 0x01, 0x2a,
	0x66, 0x0a, 0x08, 0x52, 0x74, 0x63, 0x53, 0x74, 0x61, 0x74, 0x65, 0x12, 0x13, 0x0a, 0x0f, 0x52,
	0x74, 0x63, 0x53, 0x74, 0x61, 0x74, 0x65, 0x44, 0x65, 0x66, 0x61, 0x75, 0x6c, 0x74, 0x10, 0x00,
	0x12, 0x0f, 0x0a, 0x0b, 0x52, 0x74, 0x63, 0x49, 0x6e, 0x63, 0x6f, 0x6d, 0x69, 0x6e, 0x67, 0x10,
	0x01, 0x12, 0x0f, 0x0a, 0x0b, 0x52, 0x74, 0x63, 0x4f, 0x75, 0x74, 0x67, 0x6f, 0x69, 0x6e, 0x67,
	0x10, 0x02, 0x12, 0x11, 0x0a, 0x0d, 0x52, 0x74, 0x63, 0x43, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74,
	0x69, 0x6e, 0x67, 0x10, 0x03, 0x12, 0x10, 0x0a, 0x0c, 0x52, 0x74, 0x63, 0x43, 0x6f, 0x6e, 0x6e,
	0x65, 0x63, 0x74, 0x65, 0x64, 0x10, 0x04, 0x2a, 0x66, 0x0a, 0x10, 0x52, 0x74, 0x63, 0x52, 0x6f,
	0x6f, 0x6d, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x54, 0x79, 0x70, 0x65, 0x12, 0x17, 0x0a, 0x13, 0x44,
	0x65, 0x66, 0x61, 0x75, 0x6c, 0x74, 0x52, 0x74, 0x63, 0x52, 0x6f, 0x6f, 0x6d, 0x45, 0x76, 0x65,
	0x6e, 0x74, 0x10, 0x00, 0x12, 0x0b, 0x0a, 0x07, 0x52, 0x74, 0x63, 0x4a, 0x6f, 0x69, 0x6e, 0x10,
	0x01, 0x12, 0x0b, 0x0a, 0x07, 0x52, 0x74, 0x63, 0x51, 0x75, 0x69, 0x74, 0x10, 0x02, 0x12, 0x0e,
	0x0a, 0x0a, 0x52, 0x74, 0x63, 0x44, 0x65, 0x73, 0x74, 0x72, 0x6f, 0x79, 0x10, 0x03, 0x12, 0x0f,
	0x0a, 0x0b, 0x52, 0x74, 0x63, 0x53, 0x74, 0x61, 0x74, 0x65, 0x43, 0x68, 0x67, 0x10, 0x04, 0x2a,
	0x41, 0x0a, 0x11, 0x52, 0x74, 0x63, 0x52, 0x6f, 0x6f, 0x6d, 0x51, 0x75, 0x69, 0x74, 0x52, 0x65,
	0x61, 0x73, 0x6f, 0x6e, 0x12, 0x0a, 0x0a, 0x06, 0x41, 0x63, 0x74, 0x69, 0x76, 0x65, 0x10, 0x00,
	0x12, 0x0f, 0x0a, 0x0b, 0x43, 0x61, 0x6c, 0x6c, 0x54, 0x69, 0x6d, 0x65, 0x6f, 0x75, 0x74, 0x10,
	0x01, 0x12, 0x0f, 0x0a, 0x0b, 0x50, 0x69, 0x6e, 0x67, 0x54, 0x69, 0x6d, 0x65, 0x6f, 0x75, 0x74,
	0x10, 0x02, 0x42, 0x1a, 0x5a, 0x18, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x73, 0x2f, 0x70, 0x62,
	0x64, 0x65, 0x66, 0x69, 0x6e, 0x65, 0x73, 0x2f, 0x70, 0x62, 0x6f, 0x62, 0x6a, 0x73, 0x62, 0x06,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_commons_pbdefines_rtcroom_proto_rawDescOnce sync.Once
	file_commons_pbdefines_rtcroom_proto_rawDescData = file_commons_pbdefines_rtcroom_proto_rawDesc
)

func file_commons_pbdefines_rtcroom_proto_rawDescGZIP() []byte {
	file_commons_pbdefines_rtcroom_proto_rawDescOnce.Do(func() {
		file_commons_pbdefines_rtcroom_proto_rawDescData = protoimpl.X.CompressGZIP(file_commons_pbdefines_rtcroom_proto_rawDescData)
	})
	return file_commons_pbdefines_rtcroom_proto_rawDescData
}

var file_commons_pbdefines_rtcroom_proto_enumTypes = make([]protoimpl.EnumInfo, 4)
var file_commons_pbdefines_rtcroom_proto_msgTypes = make([]protoimpl.MessageInfo, 9)
var file_commons_pbdefines_rtcroom_proto_goTypes = []interface{}{
	(RtcRoomType)(0),           // 0: RtcRoomType
	(RtcState)(0),              // 1: RtcState
	(RtcRoomEventType)(0),      // 2: RtcRoomEventType
	(RtcRoomQuitReason)(0),     // 3: RtcRoomQuitReason
	(*RtcRoomReq)(nil),         // 4: RtcRoomReq
	(*RtcRoom)(nil),            // 5: RtcRoom
	(*RtcMember)(nil),          // 6: RtcMember
	(*MemberState)(nil),        // 7: MemberState
	(*SyncMemberStateReq)(nil), // 8: SyncMemberStateReq
	(*RtcRoomEvent)(nil),       // 9: RtcRoomEvent
	(*RtcInviteReq)(nil),       // 10: RtcInviteReq
	(*RtcMemberRooms)(nil),     // 11: RtcMemberRooms
	(*RtcMemberRoom)(nil),      // 12: RtcMemberRoom
	(*UserInfo)(nil),           // 13: UserInfo
}
var file_commons_pbdefines_rtcroom_proto_depIdxs = []int32{
	0,  // 0: RtcRoomReq.roomType:type_name -> RtcRoomType
	6,  // 1: RtcRoomReq.joinMember:type_name -> RtcMember
	0,  // 2: RtcRoom.roomType:type_name -> RtcRoomType
	13, // 3: RtcRoom.owner:type_name -> UserInfo
	6,  // 4: RtcRoom.members:type_name -> RtcMember
	13, // 5: RtcMember.member:type_name -> UserInfo
	1,  // 6: RtcMember.rtcState:type_name -> RtcState
	13, // 7: RtcMember.inviter:type_name -> UserInfo
	0,  // 8: MemberState.roomType:type_name -> RtcRoomType
	1,  // 9: MemberState.rtcState:type_name -> RtcState
	7,  // 10: SyncMemberStateReq.member:type_name -> MemberState
	2,  // 11: RtcRoomEvent.roomEventType:type_name -> RtcRoomEventType
	6,  // 12: RtcRoomEvent.member:type_name -> RtcMember
	5,  // 13: RtcRoomEvent.room:type_name -> RtcRoom
	3,  // 14: RtcRoomEvent.reason:type_name -> RtcRoomQuitReason
	0,  // 15: RtcInviteReq.roomType:type_name -> RtcRoomType
	12, // 16: RtcMemberRooms.rooms:type_name -> RtcMemberRoom
	0,  // 17: RtcMemberRoom.roomType:type_name -> RtcRoomType
	13, // 18: RtcMemberRoom.owner:type_name -> UserInfo
	1,  // 19: RtcMemberRoom.rtcState:type_name -> RtcState
	20, // [20:20] is the sub-list for method output_type
	20, // [20:20] is the sub-list for method input_type
	20, // [20:20] is the sub-list for extension type_name
	20, // [20:20] is the sub-list for extension extendee
	0,  // [0:20] is the sub-list for field type_name
}

func init() { file_commons_pbdefines_rtcroom_proto_init() }
func file_commons_pbdefines_rtcroom_proto_init() {
	if File_commons_pbdefines_rtcroom_proto != nil {
		return
	}
	file_commons_pbdefines_appmessages_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_commons_pbdefines_rtcroom_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RtcRoomReq); i {
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
		file_commons_pbdefines_rtcroom_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RtcRoom); i {
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
		file_commons_pbdefines_rtcroom_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RtcMember); i {
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
		file_commons_pbdefines_rtcroom_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MemberState); i {
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
		file_commons_pbdefines_rtcroom_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SyncMemberStateReq); i {
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
		file_commons_pbdefines_rtcroom_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RtcRoomEvent); i {
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
		file_commons_pbdefines_rtcroom_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RtcInviteReq); i {
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
		file_commons_pbdefines_rtcroom_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RtcMemberRooms); i {
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
		file_commons_pbdefines_rtcroom_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RtcMemberRoom); i {
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
			RawDescriptor: file_commons_pbdefines_rtcroom_proto_rawDesc,
			NumEnums:      4,
			NumMessages:   9,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_commons_pbdefines_rtcroom_proto_goTypes,
		DependencyIndexes: file_commons_pbdefines_rtcroom_proto_depIdxs,
		EnumInfos:         file_commons_pbdefines_rtcroom_proto_enumTypes,
		MessageInfos:      file_commons_pbdefines_rtcroom_proto_msgTypes,
	}.Build()
	File_commons_pbdefines_rtcroom_proto = out.File
	file_commons_pbdefines_rtcroom_proto_rawDesc = nil
	file_commons_pbdefines_rtcroom_proto_goTypes = nil
	file_commons_pbdefines_rtcroom_proto_depIdxs = nil
}
