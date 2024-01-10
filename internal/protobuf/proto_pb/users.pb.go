// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.31.0
// 	protoc        v4.25.1
// source: users.proto

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

// Permission 标识本用户个人信息对于其他人的权限
type Permission int32

const (
	// 所有人可见
	Permission_PERMISSION_PUBLIC Permission = 0
	// 仅个人可见
	Permission_PERMISSION_PRIVATE Permission = 1
	// 仅指定可见
	Permission_PERMISSION_PROTECTED Permission = 2
)

// Enum value maps for Permission.
var (
	Permission_name = map[int32]string{
		0: "PERMISSION_PUBLIC",
		1: "PERMISSION_PRIVATE",
		2: "PERMISSION_PROTECTED",
	}
	Permission_value = map[string]int32{
		"PERMISSION_PUBLIC":    0,
		"PERMISSION_PRIVATE":   1,
		"PERMISSION_PROTECTED": 2,
	}
)

func (x Permission) Enum() *Permission {
	p := new(Permission)
	*p = x
	return p
}

func (x Permission) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Permission) Descriptor() protoreflect.EnumDescriptor {
	return file_users_proto_enumTypes[0].Descriptor()
}

func (Permission) Type() protoreflect.EnumType {
	return &file_users_proto_enumTypes[0]
}

func (x Permission) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Permission.Descriptor instead.
func (Permission) EnumDescriptor() ([]byte, []int) {
	return file_users_proto_rawDescGZIP(), []int{0}
}

type Protection int32

const (
	Protection_PROTECTION_FOR_FRIENDS_ALL Protection = 0
	Protection_PROTECTION_FOR_SOMEBODY    Protection = 1
)

// Enum value maps for Protection.
var (
	Protection_name = map[int32]string{
		0: "PROTECTION_FOR_FRIENDS_ALL",
		1: "PROTECTION_FOR_SOMEBODY",
	}
	Protection_value = map[string]int32{
		"PROTECTION_FOR_FRIENDS_ALL": 0,
		"PROTECTION_FOR_SOMEBODY":    1,
	}
)

func (x Protection) Enum() *Protection {
	p := new(Protection)
	*p = x
	return p
}

func (x Protection) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Protection) Descriptor() protoreflect.EnumDescriptor {
	return file_users_proto_enumTypes[1].Descriptor()
}

func (Protection) Type() protoreflect.EnumType {
	return &file_users_proto_enumTypes[1]
}

func (x Protection) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Protection.Descriptor instead.
func (Protection) EnumDescriptor() ([]byte, []int) {
	return file_users_proto_rawDescGZIP(), []int{1}
}

type UsersProfile struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id           uint64   `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	UserName     string   `protobuf:"bytes,2,opt,name=user_name,json=userName,proto3" json:"user_name,omitempty"`
	Email        string   `protobuf:"bytes,3,opt,name=email,proto3" json:"email,omitempty"`
	Profile      string   `protobuf:"bytes,4,opt,name=profile,proto3" json:"profile,omitempty"`
	Location     string   `protobuf:"bytes,5,opt,name=location,proto3" json:"location,omitempty"`
	Organization string   `protobuf:"bytes,6,opt,name=organization,proto3" json:"organization,omitempty"`
	Avatar       []byte   `protobuf:"bytes,7,opt,name=avatar,proto3" json:"avatar,omitempty"`
	GroupIds     []uint64 `protobuf:"varint,9,rep,packed,name=group_ids,json=groupIds,proto3" json:"group_ids,omitempty"`
	// 权限控制
	Permission Permission `protobuf:"varint,10,opt,name=permission,proto3,enum=message.Permission" json:"permission,omitempty"`
	// 仅当 permission = PERMISSION_PROTECTED 时生效
	Protection Protection `protobuf:"varint,11,opt,name=protection,proto3,enum=message.Protection" json:"protection,omitempty"` // 仅当 protection = PROTECTION_FOR_SOMEBODY 时生效
}

func (x *UsersProfile) Reset() {
	*x = UsersProfile{}
	if protoimpl.UnsafeEnabled {
		mi := &file_users_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UsersProfile) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UsersProfile) ProtoMessage() {}

func (x *UsersProfile) ProtoReflect() protoreflect.Message {
	mi := &file_users_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UsersProfile.ProtoReflect.Descriptor instead.
func (*UsersProfile) Descriptor() ([]byte, []int) {
	return file_users_proto_rawDescGZIP(), []int{0}
}

func (x *UsersProfile) GetId() uint64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *UsersProfile) GetUserName() string {
	if x != nil {
		return x.UserName
	}
	return ""
}

func (x *UsersProfile) GetEmail() string {
	if x != nil {
		return x.Email
	}
	return ""
}

func (x *UsersProfile) GetProfile() string {
	if x != nil {
		return x.Profile
	}
	return ""
}

func (x *UsersProfile) GetLocation() string {
	if x != nil {
		return x.Location
	}
	return ""
}

func (x *UsersProfile) GetOrganization() string {
	if x != nil {
		return x.Organization
	}
	return ""
}

func (x *UsersProfile) GetAvatar() []byte {
	if x != nil {
		return x.Avatar
	}
	return nil
}

func (x *UsersProfile) GetGroupIds() []uint64 {
	if x != nil {
		return x.GroupIds
	}
	return nil
}

func (x *UsersProfile) GetPermission() Permission {
	if x != nil {
		return x.Permission
	}
	return Permission_PERMISSION_PUBLIC
}

func (x *UsersProfile) GetProtection() Protection {
	if x != nil {
		return x.Protection
	}
	return Protection_PROTECTION_FOR_FRIENDS_ALL
}

type Friends struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id       uint64 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	UserName string `protobuf:"bytes,2,opt,name=user_name,json=userName,proto3" json:"user_name,omitempty"`
	Alias    string `protobuf:"bytes,3,opt,name=alias,proto3" json:"alias,omitempty"`
	Summary  string `protobuf:"bytes,4,opt,name=summary,proto3" json:"summary,omitempty"`
	Avatar   string `protobuf:"bytes,5,opt,name=avatar,proto3" json:"avatar,omitempty"`
	IsOnline bool   `protobuf:"varint,6,opt,name=is_online,json=isOnline,proto3" json:"is_online,omitempty"`
}

func (x *Friends) Reset() {
	*x = Friends{}
	if protoimpl.UnsafeEnabled {
		mi := &file_users_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Friends) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Friends) ProtoMessage() {}

func (x *Friends) ProtoReflect() protoreflect.Message {
	mi := &file_users_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Friends.ProtoReflect.Descriptor instead.
func (*Friends) Descriptor() ([]byte, []int) {
	return file_users_proto_rawDescGZIP(), []int{1}
}

func (x *Friends) GetId() uint64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *Friends) GetUserName() string {
	if x != nil {
		return x.UserName
	}
	return ""
}

func (x *Friends) GetAlias() string {
	if x != nil {
		return x.Alias
	}
	return ""
}

func (x *Friends) GetSummary() string {
	if x != nil {
		return x.Summary
	}
	return ""
}

func (x *Friends) GetAvatar() string {
	if x != nil {
		return x.Avatar
	}
	return ""
}

func (x *Friends) GetIsOnline() bool {
	if x != nil {
		return x.IsOnline
	}
	return false
}

type FriendList struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Friends []*Friends `protobuf:"bytes,1,rep,name=friends,proto3" json:"friends,omitempty"`
	// 仅当 permission = PERMISSION_PROTECTED 时生效
	ProtectedIds []uint64 `protobuf:"varint,2,rep,packed,name=protected_ids,json=protectedIds,proto3" json:"protected_ids,omitempty"`
}

func (x *FriendList) Reset() {
	*x = FriendList{}
	if protoimpl.UnsafeEnabled {
		mi := &file_users_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FriendList) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FriendList) ProtoMessage() {}

func (x *FriendList) ProtoReflect() protoreflect.Message {
	mi := &file_users_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FriendList.ProtoReflect.Descriptor instead.
func (*FriendList) Descriptor() ([]byte, []int) {
	return file_users_proto_rawDescGZIP(), []int{2}
}

func (x *FriendList) GetFriends() []*Friends {
	if x != nil {
		return x.Friends
	}
	return nil
}

func (x *FriendList) GetProtectedIds() []uint64 {
	if x != nil {
		return x.ProtectedIds
	}
	return nil
}

var File_users_proto protoreflect.FileDescriptor

var file_users_proto_rawDesc = []byte{
	0x0a, 0x0b, 0x75, 0x73, 0x65, 0x72, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x07, 0x6d,
	0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x22, 0xca, 0x02, 0x0a, 0x0c, 0x55, 0x73, 0x65, 0x72, 0x73,
	0x50, 0x72, 0x6f, 0x66, 0x69, 0x6c, 0x65, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x04, 0x52, 0x02, 0x69, 0x64, 0x12, 0x1b, 0x0a, 0x09, 0x75, 0x73, 0x65, 0x72, 0x5f,
	0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x75, 0x73, 0x65, 0x72,
	0x4e, 0x61, 0x6d, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x65, 0x6d, 0x61, 0x69, 0x6c, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x05, 0x65, 0x6d, 0x61, 0x69, 0x6c, 0x12, 0x18, 0x0a, 0x07, 0x70, 0x72,
	0x6f, 0x66, 0x69, 0x6c, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x70, 0x72, 0x6f,
	0x66, 0x69, 0x6c, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x6c, 0x6f, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x6c, 0x6f, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x12, 0x22, 0x0a, 0x0c, 0x6f, 0x72, 0x67, 0x61, 0x6e, 0x69, 0x7a, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x6f, 0x72, 0x67, 0x61, 0x6e, 0x69, 0x7a, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x12, 0x16, 0x0a, 0x06, 0x61, 0x76, 0x61, 0x74, 0x61, 0x72, 0x18, 0x07,
	0x20, 0x01, 0x28, 0x0c, 0x52, 0x06, 0x61, 0x76, 0x61, 0x74, 0x61, 0x72, 0x12, 0x1b, 0x0a, 0x09,
	0x67, 0x72, 0x6f, 0x75, 0x70, 0x5f, 0x69, 0x64, 0x73, 0x18, 0x09, 0x20, 0x03, 0x28, 0x04, 0x52,
	0x08, 0x67, 0x72, 0x6f, 0x75, 0x70, 0x49, 0x64, 0x73, 0x12, 0x33, 0x0a, 0x0a, 0x70, 0x65, 0x72,
	0x6d, 0x69, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x13, 0x2e,
	0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x2e, 0x50, 0x65, 0x72, 0x6d, 0x69, 0x73, 0x73, 0x69,
	0x6f, 0x6e, 0x52, 0x0a, 0x70, 0x65, 0x72, 0x6d, 0x69, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x33,
	0x0a, 0x0a, 0x70, 0x72, 0x6f, 0x74, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x0b, 0x20, 0x01,
	0x28, 0x0e, 0x32, 0x13, 0x2e, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x2e, 0x50, 0x72, 0x6f,
	0x74, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x0a, 0x70, 0x72, 0x6f, 0x74, 0x65, 0x63, 0x74,
	0x69, 0x6f, 0x6e, 0x22, 0x9b, 0x01, 0x0a, 0x07, 0x46, 0x72, 0x69, 0x65, 0x6e, 0x64, 0x73, 0x12,
	0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x02, 0x69, 0x64, 0x12,
	0x1b, 0x0a, 0x09, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x08, 0x75, 0x73, 0x65, 0x72, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x14, 0x0a, 0x05,
	0x61, 0x6c, 0x69, 0x61, 0x73, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x61, 0x6c, 0x69,
	0x61, 0x73, 0x12, 0x18, 0x0a, 0x07, 0x73, 0x75, 0x6d, 0x6d, 0x61, 0x72, 0x79, 0x18, 0x04, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x07, 0x73, 0x75, 0x6d, 0x6d, 0x61, 0x72, 0x79, 0x12, 0x16, 0x0a, 0x06,
	0x61, 0x76, 0x61, 0x74, 0x61, 0x72, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x61, 0x76,
	0x61, 0x74, 0x61, 0x72, 0x12, 0x1b, 0x0a, 0x09, 0x69, 0x73, 0x5f, 0x6f, 0x6e, 0x6c, 0x69, 0x6e,
	0x65, 0x18, 0x06, 0x20, 0x01, 0x28, 0x08, 0x52, 0x08, 0x69, 0x73, 0x4f, 0x6e, 0x6c, 0x69, 0x6e,
	0x65, 0x22, 0x5d, 0x0a, 0x0a, 0x46, 0x72, 0x69, 0x65, 0x6e, 0x64, 0x4c, 0x69, 0x73, 0x74, 0x12,
	0x2a, 0x0a, 0x07, 0x66, 0x72, 0x69, 0x65, 0x6e, 0x64, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b,
	0x32, 0x10, 0x2e, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x2e, 0x46, 0x72, 0x69, 0x65, 0x6e,
	0x64, 0x73, 0x52, 0x07, 0x66, 0x72, 0x69, 0x65, 0x6e, 0x64, 0x73, 0x12, 0x23, 0x0a, 0x0d, 0x70,
	0x72, 0x6f, 0x74, 0x65, 0x63, 0x74, 0x65, 0x64, 0x5f, 0x69, 0x64, 0x73, 0x18, 0x02, 0x20, 0x03,
	0x28, 0x04, 0x52, 0x0c, 0x70, 0x72, 0x6f, 0x74, 0x65, 0x63, 0x74, 0x65, 0x64, 0x49, 0x64, 0x73,
	0x2a, 0x55, 0x0a, 0x0a, 0x50, 0x65, 0x72, 0x6d, 0x69, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x15,
	0x0a, 0x11, 0x50, 0x45, 0x52, 0x4d, 0x49, 0x53, 0x53, 0x49, 0x4f, 0x4e, 0x5f, 0x50, 0x55, 0x42,
	0x4c, 0x49, 0x43, 0x10, 0x00, 0x12, 0x16, 0x0a, 0x12, 0x50, 0x45, 0x52, 0x4d, 0x49, 0x53, 0x53,
	0x49, 0x4f, 0x4e, 0x5f, 0x50, 0x52, 0x49, 0x56, 0x41, 0x54, 0x45, 0x10, 0x01, 0x12, 0x18, 0x0a,
	0x14, 0x50, 0x45, 0x52, 0x4d, 0x49, 0x53, 0x53, 0x49, 0x4f, 0x4e, 0x5f, 0x50, 0x52, 0x4f, 0x54,
	0x45, 0x43, 0x54, 0x45, 0x44, 0x10, 0x02, 0x2a, 0x49, 0x0a, 0x0a, 0x50, 0x72, 0x6f, 0x74, 0x65,
	0x63, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x1e, 0x0a, 0x1a, 0x50, 0x52, 0x4f, 0x54, 0x45, 0x43, 0x54,
	0x49, 0x4f, 0x4e, 0x5f, 0x46, 0x4f, 0x52, 0x5f, 0x46, 0x52, 0x49, 0x45, 0x4e, 0x44, 0x53, 0x5f,
	0x41, 0x4c, 0x4c, 0x10, 0x00, 0x12, 0x1b, 0x0a, 0x17, 0x50, 0x52, 0x4f, 0x54, 0x45, 0x43, 0x54,
	0x49, 0x4f, 0x4e, 0x5f, 0x46, 0x4f, 0x52, 0x5f, 0x53, 0x4f, 0x4d, 0x45, 0x42, 0x4f, 0x44, 0x59,
	0x10, 0x01, 0x42, 0x0c, 0x5a, 0x0a, 0x2e, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x5f, 0x70, 0x62,
	0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_users_proto_rawDescOnce sync.Once
	file_users_proto_rawDescData = file_users_proto_rawDesc
)

func file_users_proto_rawDescGZIP() []byte {
	file_users_proto_rawDescOnce.Do(func() {
		file_users_proto_rawDescData = protoimpl.X.CompressGZIP(file_users_proto_rawDescData)
	})
	return file_users_proto_rawDescData
}

var file_users_proto_enumTypes = make([]protoimpl.EnumInfo, 2)
var file_users_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_users_proto_goTypes = []interface{}{
	(Permission)(0),      // 0: message.Permission
	(Protection)(0),      // 1: message.Protection
	(*UsersProfile)(nil), // 2: message.UsersProfile
	(*Friends)(nil),      // 3: message.Friends
	(*FriendList)(nil),   // 4: message.FriendList
}
var file_users_proto_depIdxs = []int32{
	0, // 0: message.UsersProfile.permission:type_name -> message.Permission
	1, // 1: message.UsersProfile.protection:type_name -> message.Protection
	3, // 2: message.FriendList.friends:type_name -> message.Friends
	3, // [3:3] is the sub-list for method output_type
	3, // [3:3] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_users_proto_init() }
func file_users_proto_init() {
	if File_users_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_users_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UsersProfile); i {
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
		file_users_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Friends); i {
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
		file_users_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*FriendList); i {
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
			RawDescriptor: file_users_proto_rawDesc,
			NumEnums:      2,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_users_proto_goTypes,
		DependencyIndexes: file_users_proto_depIdxs,
		EnumInfos:         file_users_proto_enumTypes,
		MessageInfos:      file_users_proto_msgTypes,
	}.Build()
	File_users_proto = out.File
	file_users_proto_rawDesc = nil
	file_users_proto_goTypes = nil
	file_users_proto_depIdxs = nil
}
