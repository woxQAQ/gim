// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.31.0
// 	protoc        v4.25.1
// source: kafka_messages.proto

package proto_pb

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

// SingleMessage 用于单聊消息
type SingleMessage struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Content    string `protobuf:"bytes,1,opt,name=content,proto3" json:"content,omitempty"`
	SenderId   uint64 `protobuf:"varint,2,opt,name=sender_id,json=senderId,proto3" json:"sender_id,omitempty"`
	ReceiverId uint64 `protobuf:"varint,3,opt,name=receiver_id,json=receiverId,proto3" json:"receiver_id,omitempty"`
	Timestamp  int64  `protobuf:"varint,4,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
}

func (x *SingleMessage) Reset() {
	*x = SingleMessage{}
	if protoimpl.UnsafeEnabled {
		mi := &file_kafka_messages_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SingleMessage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SingleMessage) ProtoMessage() {}

func (x *SingleMessage) ProtoReflect() protoreflect.Message {
	mi := &file_kafka_messages_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SingleMessage.ProtoReflect.Descriptor instead.
func (*SingleMessage) Descriptor() ([]byte, []int) {
	return file_kafka_messages_proto_rawDescGZIP(), []int{0}
}

func (x *SingleMessage) GetContent() string {
	if x != nil {
		return x.Content
	}
	return ""
}

func (x *SingleMessage) GetSenderId() uint64 {
	if x != nil {
		return x.SenderId
	}
	return 0
}

func (x *SingleMessage) GetReceiverId() uint64 {
	if x != nil {
		return x.ReceiverId
	}
	return 0
}

func (x *SingleMessage) GetTimestamp() int64 {
	if x != nil {
		return x.Timestamp
	}
	return 0
}

// GroupMessage 用于群聊消息
type GroupMessage struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	MessageId uint64 `protobuf:"varint,1,opt,name=message_id,json=messageId,proto3" json:"message_id,omitempty"`
	Content   string `protobuf:"bytes,2,opt,name=content,proto3" json:"content,omitempty"`
	SenderId  uint64 `protobuf:"varint,3,opt,name=sender_id,json=senderId,proto3" json:"sender_id,omitempty"`
	GroupId   uint64 `protobuf:"varint,4,opt,name=group_id,json=groupId,proto3" json:"group_id,omitempty"`
	Timestamp int64  `protobuf:"varint,5,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
}

func (x *GroupMessage) Reset() {
	*x = GroupMessage{}
	if protoimpl.UnsafeEnabled {
		mi := &file_kafka_messages_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GroupMessage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GroupMessage) ProtoMessage() {}

func (x *GroupMessage) ProtoReflect() protoreflect.Message {
	mi := &file_kafka_messages_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GroupMessage.ProtoReflect.Descriptor instead.
func (*GroupMessage) Descriptor() ([]byte, []int) {
	return file_kafka_messages_proto_rawDescGZIP(), []int{1}
}

func (x *GroupMessage) GetMessageId() uint64 {
	if x != nil {
		return x.MessageId
	}
	return 0
}

func (x *GroupMessage) GetContent() string {
	if x != nil {
		return x.Content
	}
	return ""
}

func (x *GroupMessage) GetSenderId() uint64 {
	if x != nil {
		return x.SenderId
	}
	return 0
}

func (x *GroupMessage) GetGroupId() uint64 {
	if x != nil {
		return x.GroupId
	}
	return 0
}

func (x *GroupMessage) GetTimestamp() int64 {
	if x != nil {
		return x.Timestamp
	}
	return 0
}

var File_kafka_messages_proto protoreflect.FileDescriptor

var file_kafka_messages_proto_rawDesc = []byte{
	0x0a, 0x14, 0x6b, 0x61, 0x66, 0x6b, 0x61, 0x5f, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x73,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x22,
	0x85, 0x01, 0x0a, 0x0d, 0x53, 0x69, 0x6e, 0x67, 0x6c, 0x65, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67,
	0x65, 0x12, 0x18, 0x0a, 0x07, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x07, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x12, 0x1b, 0x0a, 0x09, 0x73,
	0x65, 0x6e, 0x64, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x04, 0x52, 0x08,
	0x73, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x49, 0x64, 0x12, 0x1f, 0x0a, 0x0b, 0x72, 0x65, 0x63, 0x65,
	0x69, 0x76, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x04, 0x52, 0x0a, 0x72,
	0x65, 0x63, 0x65, 0x69, 0x76, 0x65, 0x72, 0x49, 0x64, 0x12, 0x1c, 0x0a, 0x09, 0x74, 0x69, 0x6d,
	0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x18, 0x04, 0x20, 0x01, 0x28, 0x03, 0x52, 0x09, 0x74, 0x69,
	0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x22, 0x9d, 0x01, 0x0a, 0x0c, 0x47, 0x72, 0x6f, 0x75,
	0x70, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x1d, 0x0a, 0x0a, 0x6d, 0x65, 0x73, 0x73,
	0x61, 0x67, 0x65, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x09, 0x6d, 0x65,
	0x73, 0x73, 0x61, 0x67, 0x65, 0x49, 0x64, 0x12, 0x18, 0x0a, 0x07, 0x63, 0x6f, 0x6e, 0x74, 0x65,
	0x6e, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e,
	0x74, 0x12, 0x1b, 0x0a, 0x09, 0x73, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x04, 0x52, 0x08, 0x73, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x49, 0x64, 0x12, 0x19,
	0x0a, 0x08, 0x67, 0x72, 0x6f, 0x75, 0x70, 0x5f, 0x69, 0x64, 0x18, 0x04, 0x20, 0x01, 0x28, 0x04,
	0x52, 0x07, 0x67, 0x72, 0x6f, 0x75, 0x70, 0x49, 0x64, 0x12, 0x1c, 0x0a, 0x09, 0x74, 0x69, 0x6d,
	0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x18, 0x05, 0x20, 0x01, 0x28, 0x03, 0x52, 0x09, 0x74, 0x69,
	0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x42, 0x0c, 0x5a, 0x0a, 0x2e, 0x2f, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x5f, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_kafka_messages_proto_rawDescOnce sync.Once
	file_kafka_messages_proto_rawDescData = file_kafka_messages_proto_rawDesc
)

func file_kafka_messages_proto_rawDescGZIP() []byte {
	file_kafka_messages_proto_rawDescOnce.Do(func() {
		file_kafka_messages_proto_rawDescData = protoimpl.X.CompressGZIP(file_kafka_messages_proto_rawDescData)
	})
	return file_kafka_messages_proto_rawDescData
}

var file_kafka_messages_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_kafka_messages_proto_goTypes = []interface{}{
	(*SingleMessage)(nil), // 0: message.SingleMessage
	(*GroupMessage)(nil),  // 1: message.GroupMessage
}
var file_kafka_messages_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_kafka_messages_proto_init() }
func file_kafka_messages_proto_init() {
	if File_kafka_messages_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_kafka_messages_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SingleMessage); i {
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
		file_kafka_messages_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GroupMessage); i {
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
			RawDescriptor: file_kafka_messages_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_kafka_messages_proto_goTypes,
		DependencyIndexes: file_kafka_messages_proto_depIdxs,
		MessageInfos:      file_kafka_messages_proto_msgTypes,
	}.Build()
	File_kafka_messages_proto = out.File
	file_kafka_messages_proto_rawDesc = nil
	file_kafka_messages_proto_goTypes = nil
	file_kafka_messages_proto_depIdxs = nil
}
