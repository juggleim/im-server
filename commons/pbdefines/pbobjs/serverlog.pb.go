// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.33.0
// 	protoc        v3.5.0
// source: commons/pbdefines/serverlog.proto

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

type QryServerLogsReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	LogType string `protobuf:"bytes,1,opt,name=logType,proto3" json:"logType,omitempty"`
	UserId  string `protobuf:"bytes,2,opt,name=userId,proto3" json:"userId,omitempty"`
	Session string `protobuf:"bytes,3,opt,name=session,proto3" json:"session,omitempty"`
	Index   int32  `protobuf:"varint,4,opt,name=index,proto3" json:"index,omitempty"`
	Start   int64  `protobuf:"varint,5,opt,name=start,proto3" json:"start,omitempty"`
	Count   int64  `protobuf:"varint,6,opt,name=count,proto3" json:"count,omitempty"`
}

func (x *QryServerLogsReq) Reset() {
	*x = QryServerLogsReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_commons_pbdefines_serverlog_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *QryServerLogsReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*QryServerLogsReq) ProtoMessage() {}

func (x *QryServerLogsReq) ProtoReflect() protoreflect.Message {
	mi := &file_commons_pbdefines_serverlog_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use QryServerLogsReq.ProtoReflect.Descriptor instead.
func (*QryServerLogsReq) Descriptor() ([]byte, []int) {
	return file_commons_pbdefines_serverlog_proto_rawDescGZIP(), []int{0}
}

func (x *QryServerLogsReq) GetLogType() string {
	if x != nil {
		return x.LogType
	}
	return ""
}

func (x *QryServerLogsReq) GetUserId() string {
	if x != nil {
		return x.UserId
	}
	return ""
}

func (x *QryServerLogsReq) GetSession() string {
	if x != nil {
		return x.Session
	}
	return ""
}

func (x *QryServerLogsReq) GetIndex() int32 {
	if x != nil {
		return x.Index
	}
	return 0
}

func (x *QryServerLogsReq) GetStart() int64 {
	if x != nil {
		return x.Start
	}
	return 0
}

func (x *QryServerLogsReq) GetCount() int64 {
	if x != nil {
		return x.Count
	}
	return 0
}

type QryServerLogsResp struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Logs []string `protobuf:"bytes,1,rep,name=logs,proto3" json:"logs,omitempty"`
}

func (x *QryServerLogsResp) Reset() {
	*x = QryServerLogsResp{}
	if protoimpl.UnsafeEnabled {
		mi := &file_commons_pbdefines_serverlog_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *QryServerLogsResp) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*QryServerLogsResp) ProtoMessage() {}

func (x *QryServerLogsResp) ProtoReflect() protoreflect.Message {
	mi := &file_commons_pbdefines_serverlog_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use QryServerLogsResp.ProtoReflect.Descriptor instead.
func (*QryServerLogsResp) Descriptor() ([]byte, []int) {
	return file_commons_pbdefines_serverlog_proto_rawDescGZIP(), []int{1}
}

func (x *QryServerLogsResp) GetLogs() []string {
	if x != nil {
		return x.Logs
	}
	return nil
}

type LogEntities struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Entities []*LogEntity `protobuf:"bytes,1,rep,name=entities,proto3" json:"entities,omitempty"`
}

func (x *LogEntities) Reset() {
	*x = LogEntities{}
	if protoimpl.UnsafeEnabled {
		mi := &file_commons_pbdefines_serverlog_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LogEntities) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LogEntities) ProtoMessage() {}

func (x *LogEntities) ProtoReflect() protoreflect.Message {
	mi := &file_commons_pbdefines_serverlog_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LogEntities.ProtoReflect.Descriptor instead.
func (*LogEntities) Descriptor() ([]byte, []int) {
	return file_commons_pbdefines_serverlog_proto_rawDescGZIP(), []int{2}
}

func (x *LogEntities) GetEntities() []*LogEntity {
	if x != nil {
		return x.Entities
	}
	return nil
}

type LogEntity struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Types that are assignable to LogOf:
	//
	//	*LogEntity_UserConnectLog
	//	*LogEntity_ConnectionLog
	//	*LogEntity_BusinessLog
	LogOf isLogEntity_LogOf `protobuf_oneof:"logOf"`
}

func (x *LogEntity) Reset() {
	*x = LogEntity{}
	if protoimpl.UnsafeEnabled {
		mi := &file_commons_pbdefines_serverlog_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LogEntity) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LogEntity) ProtoMessage() {}

func (x *LogEntity) ProtoReflect() protoreflect.Message {
	mi := &file_commons_pbdefines_serverlog_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LogEntity.ProtoReflect.Descriptor instead.
func (*LogEntity) Descriptor() ([]byte, []int) {
	return file_commons_pbdefines_serverlog_proto_rawDescGZIP(), []int{3}
}

func (m *LogEntity) GetLogOf() isLogEntity_LogOf {
	if m != nil {
		return m.LogOf
	}
	return nil
}

func (x *LogEntity) GetUserConnectLog() *UserConnectLog {
	if x, ok := x.GetLogOf().(*LogEntity_UserConnectLog); ok {
		return x.UserConnectLog
	}
	return nil
}

func (x *LogEntity) GetConnectionLog() *ConnectionLog {
	if x, ok := x.GetLogOf().(*LogEntity_ConnectionLog); ok {
		return x.ConnectionLog
	}
	return nil
}

func (x *LogEntity) GetBusinessLog() *BusinessLog {
	if x, ok := x.GetLogOf().(*LogEntity_BusinessLog); ok {
		return x.BusinessLog
	}
	return nil
}

type isLogEntity_LogOf interface {
	isLogEntity_LogOf()
}

type LogEntity_UserConnectLog struct {
	UserConnectLog *UserConnectLog `protobuf:"bytes,11,opt,name=userConnectLog,proto3,oneof"`
}

type LogEntity_ConnectionLog struct {
	ConnectionLog *ConnectionLog `protobuf:"bytes,12,opt,name=connectionLog,proto3,oneof"`
}

type LogEntity_BusinessLog struct {
	BusinessLog *BusinessLog `protobuf:"bytes,13,opt,name=businessLog,proto3,oneof"`
}

func (*LogEntity_UserConnectLog) isLogEntity_LogOf() {}

func (*LogEntity_ConnectionLog) isLogEntity_LogOf() {}

func (*LogEntity_BusinessLog) isLogEntity_LogOf() {}

type UserConnectLog struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Timestamp int64  `protobuf:"varint,1,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
	RealTime  int64  `protobuf:"varint,2,opt,name=real_time,json=realTime,proto3" json:"real_time,omitempty"`
	AppKey    string `protobuf:"bytes,3,opt,name=app_key,json=appKey,proto3" json:"app_key,omitempty"`
	UserId    string `protobuf:"bytes,4,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	Session   string `protobuf:"bytes,5,opt,name=session,proto3" json:"session,omitempty"`
	Code      int32  `protobuf:"varint,6,opt,name=code,proto3" json:"code,omitempty"`
	Platform  string `protobuf:"bytes,7,opt,name=platform,proto3" json:"platform,omitempty"`
	Version   string `protobuf:"bytes,8,opt,name=version,proto3" json:"version,omitempty"`
	ClientIp  string `protobuf:"bytes,9,opt,name=client_ip,json=clientIp,proto3" json:"client_ip,omitempty"`
}

func (x *UserConnectLog) Reset() {
	*x = UserConnectLog{}
	if protoimpl.UnsafeEnabled {
		mi := &file_commons_pbdefines_serverlog_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UserConnectLog) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UserConnectLog) ProtoMessage() {}

func (x *UserConnectLog) ProtoReflect() protoreflect.Message {
	mi := &file_commons_pbdefines_serverlog_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UserConnectLog.ProtoReflect.Descriptor instead.
func (*UserConnectLog) Descriptor() ([]byte, []int) {
	return file_commons_pbdefines_serverlog_proto_rawDescGZIP(), []int{4}
}

func (x *UserConnectLog) GetTimestamp() int64 {
	if x != nil {
		return x.Timestamp
	}
	return 0
}

func (x *UserConnectLog) GetRealTime() int64 {
	if x != nil {
		return x.RealTime
	}
	return 0
}

func (x *UserConnectLog) GetAppKey() string {
	if x != nil {
		return x.AppKey
	}
	return ""
}

func (x *UserConnectLog) GetUserId() string {
	if x != nil {
		return x.UserId
	}
	return ""
}

func (x *UserConnectLog) GetSession() string {
	if x != nil {
		return x.Session
	}
	return ""
}

func (x *UserConnectLog) GetCode() int32 {
	if x != nil {
		return x.Code
	}
	return 0
}

func (x *UserConnectLog) GetPlatform() string {
	if x != nil {
		return x.Platform
	}
	return ""
}

func (x *UserConnectLog) GetVersion() string {
	if x != nil {
		return x.Version
	}
	return ""
}

func (x *UserConnectLog) GetClientIp() string {
	if x != nil {
		return x.ClientIp
	}
	return ""
}

type ConnectionLog struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Timestamp int64  `protobuf:"varint,1,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
	RealTime  int64  `protobuf:"varint,2,opt,name=real_time,json=realTime,proto3" json:"real_time,omitempty"`
	AppKey    string `protobuf:"bytes,3,opt,name=app_key,json=appKey,proto3" json:"app_key,omitempty"`
	Session   string `protobuf:"bytes,4,opt,name=session,proto3" json:"session,omitempty"`
	Index     int32  `protobuf:"varint,5,opt,name=index,proto3" json:"index,omitempty"`
	Action    string `protobuf:"bytes,6,opt,name=action,proto3" json:"action,omitempty"`
	Method    string `protobuf:"bytes,7,opt,name=method,proto3" json:"method,omitempty"`
	TargetId  string `protobuf:"bytes,8,opt,name=target_id,json=targetId,proto3" json:"target_id,omitempty"`
	DataLen   int32  `protobuf:"varint,9,opt,name=data_len,json=dataLen,proto3" json:"data_len,omitempty"`
	Code      int32  `protobuf:"varint,10,opt,name=code,proto3" json:"code,omitempty"`
}

func (x *ConnectionLog) Reset() {
	*x = ConnectionLog{}
	if protoimpl.UnsafeEnabled {
		mi := &file_commons_pbdefines_serverlog_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ConnectionLog) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ConnectionLog) ProtoMessage() {}

func (x *ConnectionLog) ProtoReflect() protoreflect.Message {
	mi := &file_commons_pbdefines_serverlog_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ConnectionLog.ProtoReflect.Descriptor instead.
func (*ConnectionLog) Descriptor() ([]byte, []int) {
	return file_commons_pbdefines_serverlog_proto_rawDescGZIP(), []int{5}
}

func (x *ConnectionLog) GetTimestamp() int64 {
	if x != nil {
		return x.Timestamp
	}
	return 0
}

func (x *ConnectionLog) GetRealTime() int64 {
	if x != nil {
		return x.RealTime
	}
	return 0
}

func (x *ConnectionLog) GetAppKey() string {
	if x != nil {
		return x.AppKey
	}
	return ""
}

func (x *ConnectionLog) GetSession() string {
	if x != nil {
		return x.Session
	}
	return ""
}

func (x *ConnectionLog) GetIndex() int32 {
	if x != nil {
		return x.Index
	}
	return 0
}

func (x *ConnectionLog) GetAction() string {
	if x != nil {
		return x.Action
	}
	return ""
}

func (x *ConnectionLog) GetMethod() string {
	if x != nil {
		return x.Method
	}
	return ""
}

func (x *ConnectionLog) GetTargetId() string {
	if x != nil {
		return x.TargetId
	}
	return ""
}

func (x *ConnectionLog) GetDataLen() int32 {
	if x != nil {
		return x.DataLen
	}
	return 0
}

func (x *ConnectionLog) GetCode() int32 {
	if x != nil {
		return x.Code
	}
	return 0
}

type BusinessLog struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Timestamp   string `protobuf:"bytes,1,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
	RealTime    int64  `protobuf:"varint,2,opt,name=real_time,json=realTime,proto3" json:"real_time,omitempty"`
	Session     string `protobuf:"bytes,3,opt,name=session,proto3" json:"session,omitempty"`
	Index       uint32 `protobuf:"varint,4,opt,name=index,proto3" json:"index,omitempty"`
	ServiceName string `protobuf:"bytes,5,opt,name=service_name,json=serviceName,proto3" json:"service_name,omitempty"`
	Method      string `protobuf:"bytes,6,opt,name=method,proto3" json:"method,omitempty"`
	Message     string `protobuf:"bytes,7,opt,name=message,proto3" json:"message,omitempty"`
	AppKey      string `protobuf:"bytes,8,opt,name=appKey,proto3" json:"appKey,omitempty"`
}

func (x *BusinessLog) Reset() {
	*x = BusinessLog{}
	if protoimpl.UnsafeEnabled {
		mi := &file_commons_pbdefines_serverlog_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BusinessLog) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BusinessLog) ProtoMessage() {}

func (x *BusinessLog) ProtoReflect() protoreflect.Message {
	mi := &file_commons_pbdefines_serverlog_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BusinessLog.ProtoReflect.Descriptor instead.
func (*BusinessLog) Descriptor() ([]byte, []int) {
	return file_commons_pbdefines_serverlog_proto_rawDescGZIP(), []int{6}
}

func (x *BusinessLog) GetTimestamp() string {
	if x != nil {
		return x.Timestamp
	}
	return ""
}

func (x *BusinessLog) GetRealTime() int64 {
	if x != nil {
		return x.RealTime
	}
	return 0
}

func (x *BusinessLog) GetSession() string {
	if x != nil {
		return x.Session
	}
	return ""
}

func (x *BusinessLog) GetIndex() uint32 {
	if x != nil {
		return x.Index
	}
	return 0
}

func (x *BusinessLog) GetServiceName() string {
	if x != nil {
		return x.ServiceName
	}
	return ""
}

func (x *BusinessLog) GetMethod() string {
	if x != nil {
		return x.Method
	}
	return ""
}

func (x *BusinessLog) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

func (x *BusinessLog) GetAppKey() string {
	if x != nil {
		return x.AppKey
	}
	return ""
}

var File_commons_pbdefines_serverlog_proto protoreflect.FileDescriptor

var file_commons_pbdefines_serverlog_proto_rawDesc = []byte{
	0x0a, 0x21, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x73, 0x2f, 0x70, 0x62, 0x64, 0x65, 0x66, 0x69,
	0x6e, 0x65, 0x73, 0x2f, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x6c, 0x6f, 0x67, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x22, 0xa0, 0x01, 0x0a, 0x10, 0x51, 0x72, 0x79, 0x53, 0x65, 0x72, 0x76, 0x65,
	0x72, 0x4c, 0x6f, 0x67, 0x73, 0x52, 0x65, 0x71, 0x12, 0x18, 0x0a, 0x07, 0x6c, 0x6f, 0x67, 0x54,
	0x79, 0x70, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x6c, 0x6f, 0x67, 0x54, 0x79,
	0x70, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x12, 0x18, 0x0a, 0x07, 0x73, 0x65,
	0x73, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x73, 0x65, 0x73,
	0x73, 0x69, 0x6f, 0x6e, 0x12, 0x14, 0x0a, 0x05, 0x69, 0x6e, 0x64, 0x65, 0x78, 0x18, 0x04, 0x20,
	0x01, 0x28, 0x05, 0x52, 0x05, 0x69, 0x6e, 0x64, 0x65, 0x78, 0x12, 0x14, 0x0a, 0x05, 0x73, 0x74,
	0x61, 0x72, 0x74, 0x18, 0x05, 0x20, 0x01, 0x28, 0x03, 0x52, 0x05, 0x73, 0x74, 0x61, 0x72, 0x74,
	0x12, 0x14, 0x0a, 0x05, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x06, 0x20, 0x01, 0x28, 0x03, 0x52,
	0x05, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x22, 0x27, 0x0a, 0x11, 0x51, 0x72, 0x79, 0x53, 0x65, 0x72,
	0x76, 0x65, 0x72, 0x4c, 0x6f, 0x67, 0x73, 0x52, 0x65, 0x73, 0x70, 0x12, 0x12, 0x0a, 0x04, 0x6c,
	0x6f, 0x67, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x09, 0x52, 0x04, 0x6c, 0x6f, 0x67, 0x73, 0x22,
	0x35, 0x0a, 0x0b, 0x4c, 0x6f, 0x67, 0x45, 0x6e, 0x74, 0x69, 0x74, 0x69, 0x65, 0x73, 0x12, 0x26,
	0x0a, 0x08, 0x65, 0x6e, 0x74, 0x69, 0x74, 0x69, 0x65, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b,
	0x32, 0x0a, 0x2e, 0x4c, 0x6f, 0x67, 0x45, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x52, 0x08, 0x65, 0x6e,
	0x74, 0x69, 0x74, 0x69, 0x65, 0x73, 0x22, 0xb9, 0x01, 0x0a, 0x09, 0x4c, 0x6f, 0x67, 0x45, 0x6e,
	0x74, 0x69, 0x74, 0x79, 0x12, 0x39, 0x0a, 0x0e, 0x75, 0x73, 0x65, 0x72, 0x43, 0x6f, 0x6e, 0x6e,
	0x65, 0x63, 0x74, 0x4c, 0x6f, 0x67, 0x18, 0x0b, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0f, 0x2e, 0x55,
	0x73, 0x65, 0x72, 0x43, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x4c, 0x6f, 0x67, 0x48, 0x00, 0x52,
	0x0e, 0x75, 0x73, 0x65, 0x72, 0x43, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x4c, 0x6f, 0x67, 0x12,
	0x36, 0x0a, 0x0d, 0x63, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x4c, 0x6f, 0x67,
	0x18, 0x0c, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0e, 0x2e, 0x43, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74,
	0x69, 0x6f, 0x6e, 0x4c, 0x6f, 0x67, 0x48, 0x00, 0x52, 0x0d, 0x63, 0x6f, 0x6e, 0x6e, 0x65, 0x63,
	0x74, 0x69, 0x6f, 0x6e, 0x4c, 0x6f, 0x67, 0x12, 0x30, 0x0a, 0x0b, 0x62, 0x75, 0x73, 0x69, 0x6e,
	0x65, 0x73, 0x73, 0x4c, 0x6f, 0x67, 0x18, 0x0d, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0c, 0x2e, 0x42,
	0x75, 0x73, 0x69, 0x6e, 0x65, 0x73, 0x73, 0x4c, 0x6f, 0x67, 0x48, 0x00, 0x52, 0x0b, 0x62, 0x75,
	0x73, 0x69, 0x6e, 0x65, 0x73, 0x73, 0x4c, 0x6f, 0x67, 0x42, 0x07, 0x0a, 0x05, 0x6c, 0x6f, 0x67,
	0x4f, 0x66, 0x22, 0xfe, 0x01, 0x0a, 0x0e, 0x55, 0x73, 0x65, 0x72, 0x43, 0x6f, 0x6e, 0x6e, 0x65,
	0x63, 0x74, 0x4c, 0x6f, 0x67, 0x12, 0x1c, 0x0a, 0x09, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61,
	0x6d, 0x70, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x09, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74,
	0x61, 0x6d, 0x70, 0x12, 0x1b, 0x0a, 0x09, 0x72, 0x65, 0x61, 0x6c, 0x5f, 0x74, 0x69, 0x6d, 0x65,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x08, 0x72, 0x65, 0x61, 0x6c, 0x54, 0x69, 0x6d, 0x65,
	0x12, 0x17, 0x0a, 0x07, 0x61, 0x70, 0x70, 0x5f, 0x6b, 0x65, 0x79, 0x18, 0x03, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x06, 0x61, 0x70, 0x70, 0x4b, 0x65, 0x79, 0x12, 0x17, 0x0a, 0x07, 0x75, 0x73, 0x65,
	0x72, 0x5f, 0x69, 0x64, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x75, 0x73, 0x65, 0x72,
	0x49, 0x64, 0x12, 0x18, 0x0a, 0x07, 0x73, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x05, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x07, 0x73, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x12, 0x0a, 0x04,
	0x63, 0x6f, 0x64, 0x65, 0x18, 0x06, 0x20, 0x01, 0x28, 0x05, 0x52, 0x04, 0x63, 0x6f, 0x64, 0x65,
	0x12, 0x1a, 0x0a, 0x08, 0x70, 0x6c, 0x61, 0x74, 0x66, 0x6f, 0x72, 0x6d, 0x18, 0x07, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x08, 0x70, 0x6c, 0x61, 0x74, 0x66, 0x6f, 0x72, 0x6d, 0x12, 0x18, 0x0a, 0x07,
	0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x08, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x76,
	0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x1b, 0x0a, 0x09, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74,
	0x5f, 0x69, 0x70, 0x18, 0x09, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x63, 0x6c, 0x69, 0x65, 0x6e,
	0x74, 0x49, 0x70, 0x22, 0x8f, 0x02, 0x0a, 0x0d, 0x43, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x69,
	0x6f, 0x6e, 0x4c, 0x6f, 0x67, 0x12, 0x1c, 0x0a, 0x09, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61,
	0x6d, 0x70, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x09, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74,
	0x61, 0x6d, 0x70, 0x12, 0x1b, 0x0a, 0x09, 0x72, 0x65, 0x61, 0x6c, 0x5f, 0x74, 0x69, 0x6d, 0x65,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x08, 0x72, 0x65, 0x61, 0x6c, 0x54, 0x69, 0x6d, 0x65,
	0x12, 0x17, 0x0a, 0x07, 0x61, 0x70, 0x70, 0x5f, 0x6b, 0x65, 0x79, 0x18, 0x03, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x06, 0x61, 0x70, 0x70, 0x4b, 0x65, 0x79, 0x12, 0x18, 0x0a, 0x07, 0x73, 0x65, 0x73,
	0x73, 0x69, 0x6f, 0x6e, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x73, 0x65, 0x73, 0x73,
	0x69, 0x6f, 0x6e, 0x12, 0x14, 0x0a, 0x05, 0x69, 0x6e, 0x64, 0x65, 0x78, 0x18, 0x05, 0x20, 0x01,
	0x28, 0x05, 0x52, 0x05, 0x69, 0x6e, 0x64, 0x65, 0x78, 0x12, 0x16, 0x0a, 0x06, 0x61, 0x63, 0x74,
	0x69, 0x6f, 0x6e, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x61, 0x63, 0x74, 0x69, 0x6f,
	0x6e, 0x12, 0x16, 0x0a, 0x06, 0x6d, 0x65, 0x74, 0x68, 0x6f, 0x64, 0x18, 0x07, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x06, 0x6d, 0x65, 0x74, 0x68, 0x6f, 0x64, 0x12, 0x1b, 0x0a, 0x09, 0x74, 0x61, 0x72,
	0x67, 0x65, 0x74, 0x5f, 0x69, 0x64, 0x18, 0x08, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x74, 0x61,
	0x72, 0x67, 0x65, 0x74, 0x49, 0x64, 0x12, 0x19, 0x0a, 0x08, 0x64, 0x61, 0x74, 0x61, 0x5f, 0x6c,
	0x65, 0x6e, 0x18, 0x09, 0x20, 0x01, 0x28, 0x05, 0x52, 0x07, 0x64, 0x61, 0x74, 0x61, 0x4c, 0x65,
	0x6e, 0x12, 0x12, 0x0a, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x05, 0x52,
	0x04, 0x63, 0x6f, 0x64, 0x65, 0x22, 0xe5, 0x01, 0x0a, 0x0b, 0x42, 0x75, 0x73, 0x69, 0x6e, 0x65,
	0x73, 0x73, 0x4c, 0x6f, 0x67, 0x12, 0x1c, 0x0a, 0x09, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61,
	0x6d, 0x70, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74,
	0x61, 0x6d, 0x70, 0x12, 0x1b, 0x0a, 0x09, 0x72, 0x65, 0x61, 0x6c, 0x5f, 0x74, 0x69, 0x6d, 0x65,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x08, 0x72, 0x65, 0x61, 0x6c, 0x54, 0x69, 0x6d, 0x65,
	0x12, 0x18, 0x0a, 0x07, 0x73, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x03, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x07, 0x73, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x14, 0x0a, 0x05, 0x69, 0x6e,
	0x64, 0x65, 0x78, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x05, 0x69, 0x6e, 0x64, 0x65, 0x78,
	0x12, 0x21, 0x0a, 0x0c, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x5f, 0x6e, 0x61, 0x6d, 0x65,
	0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x4e,
	0x61, 0x6d, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x6d, 0x65, 0x74, 0x68, 0x6f, 0x64, 0x18, 0x06, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x06, 0x6d, 0x65, 0x74, 0x68, 0x6f, 0x64, 0x12, 0x18, 0x0a, 0x07, 0x6d,
	0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x07, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x6d, 0x65,
	0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x61, 0x70, 0x70, 0x4b, 0x65, 0x79, 0x18,
	0x08, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x61, 0x70, 0x70, 0x4b, 0x65, 0x79, 0x42, 0x1a, 0x5a,
	0x18, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x73, 0x2f, 0x70, 0x62, 0x64, 0x65, 0x66, 0x69, 0x6e,
	0x65, 0x73, 0x2f, 0x70, 0x62, 0x6f, 0x62, 0x6a, 0x73, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x33,
}

var (
	file_commons_pbdefines_serverlog_proto_rawDescOnce sync.Once
	file_commons_pbdefines_serverlog_proto_rawDescData = file_commons_pbdefines_serverlog_proto_rawDesc
)

func file_commons_pbdefines_serverlog_proto_rawDescGZIP() []byte {
	file_commons_pbdefines_serverlog_proto_rawDescOnce.Do(func() {
		file_commons_pbdefines_serverlog_proto_rawDescData = protoimpl.X.CompressGZIP(file_commons_pbdefines_serverlog_proto_rawDescData)
	})
	return file_commons_pbdefines_serverlog_proto_rawDescData
}

var file_commons_pbdefines_serverlog_proto_msgTypes = make([]protoimpl.MessageInfo, 7)
var file_commons_pbdefines_serverlog_proto_goTypes = []interface{}{
	(*QryServerLogsReq)(nil),  // 0: QryServerLogsReq
	(*QryServerLogsResp)(nil), // 1: QryServerLogsResp
	(*LogEntities)(nil),       // 2: LogEntities
	(*LogEntity)(nil),         // 3: LogEntity
	(*UserConnectLog)(nil),    // 4: UserConnectLog
	(*ConnectionLog)(nil),     // 5: ConnectionLog
	(*BusinessLog)(nil),       // 6: BusinessLog
}
var file_commons_pbdefines_serverlog_proto_depIdxs = []int32{
	3, // 0: LogEntities.entities:type_name -> LogEntity
	4, // 1: LogEntity.userConnectLog:type_name -> UserConnectLog
	5, // 2: LogEntity.connectionLog:type_name -> ConnectionLog
	6, // 3: LogEntity.businessLog:type_name -> BusinessLog
	4, // [4:4] is the sub-list for method output_type
	4, // [4:4] is the sub-list for method input_type
	4, // [4:4] is the sub-list for extension type_name
	4, // [4:4] is the sub-list for extension extendee
	0, // [0:4] is the sub-list for field type_name
}

func init() { file_commons_pbdefines_serverlog_proto_init() }
func file_commons_pbdefines_serverlog_proto_init() {
	if File_commons_pbdefines_serverlog_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_commons_pbdefines_serverlog_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*QryServerLogsReq); i {
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
		file_commons_pbdefines_serverlog_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*QryServerLogsResp); i {
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
		file_commons_pbdefines_serverlog_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*LogEntities); i {
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
		file_commons_pbdefines_serverlog_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*LogEntity); i {
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
		file_commons_pbdefines_serverlog_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UserConnectLog); i {
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
		file_commons_pbdefines_serverlog_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ConnectionLog); i {
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
		file_commons_pbdefines_serverlog_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BusinessLog); i {
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
	file_commons_pbdefines_serverlog_proto_msgTypes[3].OneofWrappers = []interface{}{
		(*LogEntity_UserConnectLog)(nil),
		(*LogEntity_ConnectionLog)(nil),
		(*LogEntity_BusinessLog)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_commons_pbdefines_serverlog_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   7,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_commons_pbdefines_serverlog_proto_goTypes,
		DependencyIndexes: file_commons_pbdefines_serverlog_proto_depIdxs,
		MessageInfos:      file_commons_pbdefines_serverlog_proto_msgTypes,
	}.Build()
	File_commons_pbdefines_serverlog_proto = out.File
	file_commons_pbdefines_serverlog_proto_rawDesc = nil
	file_commons_pbdefines_serverlog_proto_goTypes = nil
	file_commons_pbdefines_serverlog_proto_depIdxs = nil
}