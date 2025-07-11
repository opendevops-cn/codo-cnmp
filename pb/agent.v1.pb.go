// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.31.0
// 	protoc        v5.27.2
// source: pb/agent.v1.proto

package pb

import (
	_ "github.com/Ccheers/protoc-gen-go-kratos-http/khttp"
	_ "google.golang.org/genproto/googleapis/api/annotations"
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

type ListAgentRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// 页码
	Page uint32 `protobuf:"varint,1,opt,name=page,proto3" json:"page,optional"`
	// 每页数量
	PageSize uint32 `protobuf:"varint,2,opt,name=page_size,json=pageSize,proto3" json:"page_size,optional"`
	// 模糊查询
	Keyword string `protobuf:"bytes,3,opt,name=keyword,proto3" json:"keyword,optional"`
	// 查询全部
	ListAll uint32 `protobuf:"varint,4,opt,name=list_all,json=listAll,proto3" json:"list_all,optional"`
}

func (x *ListAgentRequest) Reset() {
	*x = ListAgentRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pb_agent_v1_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListAgentRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListAgentRequest) ProtoMessage() {}

func (x *ListAgentRequest) ProtoReflect() protoreflect.Message {
	mi := &file_pb_agent_v1_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListAgentRequest.ProtoReflect.Descriptor instead.
func (*ListAgentRequest) Descriptor() ([]byte, []int) {
	return file_pb_agent_v1_proto_rawDescGZIP(), []int{0}
}

func (x *ListAgentRequest) GetPage() uint32 {
	if x != nil {
		return x.Page
	}
	return 0
}

func (x *ListAgentRequest) GetPageSize() uint32 {
	if x != nil {
		return x.PageSize
	}
	return 0
}

func (x *ListAgentRequest) GetKeyword() string {
	if x != nil {
		return x.Keyword
	}
	return ""
}

func (x *ListAgentRequest) GetListAll() uint32 {
	if x != nil {
		return x.ListAll
	}
	return 0
}

type CreateAgentRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Agent名称
	Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,optional"`
	// agent ID
	AgentId string `protobuf:"bytes,3,opt,name=agent_id,json=agentId,proto3" json:"agent_id,optional"`
}

func (x *CreateAgentRequest) Reset() {
	*x = CreateAgentRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pb_agent_v1_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreateAgentRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateAgentRequest) ProtoMessage() {}

func (x *CreateAgentRequest) ProtoReflect() protoreflect.Message {
	mi := &file_pb_agent_v1_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateAgentRequest.ProtoReflect.Descriptor instead.
func (*CreateAgentRequest) Descriptor() ([]byte, []int) {
	return file_pb_agent_v1_proto_rawDescGZIP(), []int{1}
}

func (x *CreateAgentRequest) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *CreateAgentRequest) GetAgentId() string {
	if x != nil {
		return x.AgentId
	}
	return ""
}

type CreateAgentResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id uint32 `protobuf:"varint,1,opt,name=id,proto3" json:"id,optional"`
}

func (x *CreateAgentResponse) Reset() {
	*x = CreateAgentResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pb_agent_v1_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreateAgentResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateAgentResponse) ProtoMessage() {}

func (x *CreateAgentResponse) ProtoReflect() protoreflect.Message {
	mi := &file_pb_agent_v1_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateAgentResponse.ProtoReflect.Descriptor instead.
func (*CreateAgentResponse) Descriptor() ([]byte, []int) {
	return file_pb_agent_v1_proto_rawDescGZIP(), []int{2}
}

func (x *CreateAgentResponse) GetId() uint32 {
	if x != nil {
		return x.Id
	}
	return 0
}

type DeleteAgentRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id uint32 `protobuf:"varint,1,opt,name=id,proto3" json:"id,optional"`
}

func (x *DeleteAgentRequest) Reset() {
	*x = DeleteAgentRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pb_agent_v1_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DeleteAgentRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteAgentRequest) ProtoMessage() {}

func (x *DeleteAgentRequest) ProtoReflect() protoreflect.Message {
	mi := &file_pb_agent_v1_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteAgentRequest.ProtoReflect.Descriptor instead.
func (*DeleteAgentRequest) Descriptor() ([]byte, []int) {
	return file_pb_agent_v1_proto_rawDescGZIP(), []int{3}
}

func (x *DeleteAgentRequest) GetId() uint32 {
	if x != nil {
		return x.Id
	}
	return 0
}

type DeleteAgentResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *DeleteAgentResponse) Reset() {
	*x = DeleteAgentResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pb_agent_v1_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DeleteAgentResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteAgentResponse) ProtoMessage() {}

func (x *DeleteAgentResponse) ProtoReflect() protoreflect.Message {
	mi := &file_pb_agent_v1_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteAgentResponse.ProtoReflect.Descriptor instead.
func (*DeleteAgentResponse) Descriptor() ([]byte, []int) {
	return file_pb_agent_v1_proto_rawDescGZIP(), []int{4}
}

type UpdateAgentRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// ID
	Id uint32 `protobuf:"varint,1,opt,name=id,proto3" json:"id,optional"`
	// Agent名称
	Name string `protobuf:"bytes,2,opt,name=name,proto3" json:"name,optional"`
	// agent ID
	AgentId string `protobuf:"bytes,4,opt,name=agent_id,json=agentId,proto3" json:"agent_id,optional"`
}

func (x *UpdateAgentRequest) Reset() {
	*x = UpdateAgentRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pb_agent_v1_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UpdateAgentRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateAgentRequest) ProtoMessage() {}

func (x *UpdateAgentRequest) ProtoReflect() protoreflect.Message {
	mi := &file_pb_agent_v1_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateAgentRequest.ProtoReflect.Descriptor instead.
func (*UpdateAgentRequest) Descriptor() ([]byte, []int) {
	return file_pb_agent_v1_proto_rawDescGZIP(), []int{5}
}

func (x *UpdateAgentRequest) GetId() uint32 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *UpdateAgentRequest) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *UpdateAgentRequest) GetAgentId() string {
	if x != nil {
		return x.AgentId
	}
	return ""
}

type UpdateAgentResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *UpdateAgentResponse) Reset() {
	*x = UpdateAgentResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pb_agent_v1_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UpdateAgentResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateAgentResponse) ProtoMessage() {}

func (x *UpdateAgentResponse) ProtoReflect() protoreflect.Message {
	mi := &file_pb_agent_v1_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateAgentResponse.ProtoReflect.Descriptor instead.
func (*UpdateAgentResponse) Descriptor() ([]byte, []int) {
	return file_pb_agent_v1_proto_rawDescGZIP(), []int{6}
}

type AgentItem struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// ID
	Id uint32 `protobuf:"varint,1,opt,name=id,proto3" json:"id,optional"`
	// Agent名称
	Name string `protobuf:"bytes,2,opt,name=name,proto3" json:"name,optional"`
	// agent ID
	AgentId string `protobuf:"bytes,4,opt,name=agent_id,json=agentId,proto3" json:"agent_id,optional"`
}

func (x *AgentItem) Reset() {
	*x = AgentItem{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pb_agent_v1_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AgentItem) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AgentItem) ProtoMessage() {}

func (x *AgentItem) ProtoReflect() protoreflect.Message {
	mi := &file_pb_agent_v1_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AgentItem.ProtoReflect.Descriptor instead.
func (*AgentItem) Descriptor() ([]byte, []int) {
	return file_pb_agent_v1_proto_rawDescGZIP(), []int{7}
}

func (x *AgentItem) GetId() uint32 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *AgentItem) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *AgentItem) GetAgentId() string {
	if x != nil {
		return x.AgentId
	}
	return ""
}

type ListAgentResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// 组网列表
	List  []*AgentItem `protobuf:"bytes,1,rep,name=list,proto3" json:"list,optional"`
	Total uint32       `protobuf:"varint,2,opt,name=total,proto3" json:"total,optional"`
}

func (x *ListAgentResponse) Reset() {
	*x = ListAgentResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pb_agent_v1_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListAgentResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListAgentResponse) ProtoMessage() {}

func (x *ListAgentResponse) ProtoReflect() protoreflect.Message {
	mi := &file_pb_agent_v1_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListAgentResponse.ProtoReflect.Descriptor instead.
func (*ListAgentResponse) Descriptor() ([]byte, []int) {
	return file_pb_agent_v1_proto_rawDescGZIP(), []int{8}
}

func (x *ListAgentResponse) GetList() []*AgentItem {
	if x != nil {
		return x.List
	}
	return nil
}

func (x *ListAgentResponse) GetTotal() uint32 {
	if x != nil {
		return x.Total
	}
	return 0
}

var File_pb_agent_v1_proto protoreflect.FileDescriptor

var file_pb_agent_v1_proto_rawDesc = []byte{
	0x0a, 0x11, 0x70, 0x62, 0x2f, 0x61, 0x67, 0x65, 0x6e, 0x74, 0x2e, 0x76, 0x31, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x12, 0x05, 0x61, 0x67, 0x65, 0x6e, 0x74, 0x1a, 0x1c, 0x67, 0x6f, 0x6f, 0x67,
	0x6c, 0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x61, 0x6e, 0x6e, 0x6f, 0x74, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x16, 0x6b, 0x68, 0x74, 0x74, 0x70, 0x2f,
	0x6d, 0x69, 0x64, 0x64, 0x6c, 0x65, 0x77, 0x61, 0x72, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x66, 0x69, 0x65,
	0x6c, 0x64, 0x5f, 0x62, 0x65, 0x68, 0x61, 0x76, 0x69, 0x6f, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x22, 0x78, 0x0a, 0x10, 0x4c, 0x69, 0x73, 0x74, 0x41, 0x67, 0x65, 0x6e, 0x74, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x70, 0x61, 0x67, 0x65, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x0d, 0x52, 0x04, 0x70, 0x61, 0x67, 0x65, 0x12, 0x1b, 0x0a, 0x09, 0x70, 0x61, 0x67,
	0x65, 0x5f, 0x73, 0x69, 0x7a, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x08, 0x70, 0x61,
	0x67, 0x65, 0x53, 0x69, 0x7a, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x6b, 0x65, 0x79, 0x77, 0x6f, 0x72,
	0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x6b, 0x65, 0x79, 0x77, 0x6f, 0x72, 0x64,
	0x12, 0x19, 0x0a, 0x08, 0x6c, 0x69, 0x73, 0x74, 0x5f, 0x61, 0x6c, 0x6c, 0x18, 0x04, 0x20, 0x01,
	0x28, 0x0d, 0x52, 0x07, 0x6c, 0x69, 0x73, 0x74, 0x41, 0x6c, 0x6c, 0x22, 0x4d, 0x0a, 0x12, 0x43,
	0x72, 0x65, 0x61, 0x74, 0x65, 0x41, 0x67, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x12, 0x17, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x42,
	0x03, 0xe0, 0x41, 0x02, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x1e, 0x0a, 0x08, 0x61, 0x67,
	0x65, 0x6e, 0x74, 0x5f, 0x69, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x42, 0x03, 0xe0, 0x41,
	0x02, 0x52, 0x07, 0x61, 0x67, 0x65, 0x6e, 0x74, 0x49, 0x64, 0x22, 0x2a, 0x0a, 0x13, 0x43, 0x72,
	0x65, 0x61, 0x74, 0x65, 0x41, 0x67, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x12, 0x13, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0d, 0x42, 0x03, 0xe0,
	0x41, 0x02, 0x52, 0x02, 0x69, 0x64, 0x22, 0x29, 0x0a, 0x12, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65,
	0x41, 0x67, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x13, 0x0a, 0x02,
	0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0d, 0x42, 0x03, 0xe0, 0x41, 0x02, 0x52, 0x02, 0x69,
	0x64, 0x22, 0x15, 0x0a, 0x13, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x41, 0x67, 0x65, 0x6e, 0x74,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x62, 0x0a, 0x12, 0x55, 0x70, 0x64, 0x61,
	0x74, 0x65, 0x41, 0x67, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x13,
	0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0d, 0x42, 0x03, 0xe0, 0x41, 0x02, 0x52,
	0x02, 0x69, 0x64, 0x12, 0x17, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x42, 0x03, 0xe0, 0x41, 0x02, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x1e, 0x0a, 0x08,
	0x61, 0x67, 0x65, 0x6e, 0x74, 0x5f, 0x69, 0x64, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x42, 0x03,
	0xe0, 0x41, 0x02, 0x52, 0x07, 0x61, 0x67, 0x65, 0x6e, 0x74, 0x49, 0x64, 0x22, 0x15, 0x0a, 0x13,
	0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x41, 0x67, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x22, 0x59, 0x0a, 0x09, 0x41, 0x67, 0x65, 0x6e, 0x74, 0x49, 0x74, 0x65, 0x6d,
	0x12, 0x13, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0d, 0x42, 0x03, 0xe0, 0x41,
	0x02, 0x52, 0x02, 0x69, 0x64, 0x12, 0x17, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x09, 0x42, 0x03, 0xe0, 0x41, 0x02, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x1e,
	0x0a, 0x08, 0x61, 0x67, 0x65, 0x6e, 0x74, 0x5f, 0x69, 0x64, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09,
	0x42, 0x03, 0xe0, 0x41, 0x02, 0x52, 0x07, 0x61, 0x67, 0x65, 0x6e, 0x74, 0x49, 0x64, 0x22, 0x4f,
	0x0a, 0x11, 0x4c, 0x69, 0x73, 0x74, 0x41, 0x67, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x12, 0x24, 0x0a, 0x04, 0x6c, 0x69, 0x73, 0x74, 0x18, 0x01, 0x20, 0x03, 0x28,
	0x0b, 0x32, 0x10, 0x2e, 0x61, 0x67, 0x65, 0x6e, 0x74, 0x2e, 0x41, 0x67, 0x65, 0x6e, 0x74, 0x49,
	0x74, 0x65, 0x6d, 0x52, 0x04, 0x6c, 0x69, 0x73, 0x74, 0x12, 0x14, 0x0a, 0x05, 0x74, 0x6f, 0x74,
	0x61, 0x6c, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x05, 0x74, 0x6f, 0x74, 0x61, 0x6c, 0x32,
	0xac, 0x03, 0x0a, 0x05, 0x41, 0x67, 0x65, 0x6e, 0x74, 0x12, 0x5f, 0x0a, 0x09, 0x4c, 0x69, 0x73,
	0x74, 0x41, 0x67, 0x65, 0x6e, 0x74, 0x12, 0x17, 0x2e, 0x61, 0x67, 0x65, 0x6e, 0x74, 0x2e, 0x4c,
	0x69, 0x73, 0x74, 0x41, 0x67, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a,
	0x18, 0x2e, 0x61, 0x67, 0x65, 0x6e, 0x74, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x41, 0x67, 0x65, 0x6e,
	0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x1f, 0xc2, 0xdb, 0xaa, 0x03, 0x00,
	0x82, 0xd3, 0xe4, 0x93, 0x02, 0x14, 0x12, 0x12, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x76, 0x31, 0x2f,
	0x61, 0x67, 0x65, 0x6e, 0x74, 0x2f, 0x6c, 0x69, 0x73, 0x74, 0x12, 0x6a, 0x0a, 0x0b, 0x43, 0x72,
	0x65, 0x61, 0x74, 0x65, 0x41, 0x67, 0x65, 0x6e, 0x74, 0x12, 0x19, 0x2e, 0x61, 0x67, 0x65, 0x6e,
	0x74, 0x2e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x41, 0x67, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x1a, 0x1a, 0x2e, 0x61, 0x67, 0x65, 0x6e, 0x74, 0x2e, 0x43, 0x72, 0x65,
	0x61, 0x74, 0x65, 0x41, 0x67, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x22, 0x24, 0xc2, 0xdb, 0xaa, 0x03, 0x00, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x19, 0x3a, 0x01, 0x2a,
	0x22, 0x14, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x76, 0x31, 0x2f, 0x61, 0x67, 0x65, 0x6e, 0x74, 0x2f,
	0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x12, 0x6a, 0x0a, 0x0b, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65,
	0x41, 0x67, 0x65, 0x6e, 0x74, 0x12, 0x19, 0x2e, 0x61, 0x67, 0x65, 0x6e, 0x74, 0x2e, 0x44, 0x65,
	0x6c, 0x65, 0x74, 0x65, 0x41, 0x67, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x1a, 0x1a, 0x2e, 0x61, 0x67, 0x65, 0x6e, 0x74, 0x2e, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x41,
	0x67, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x24, 0xc2, 0xdb,
	0xaa, 0x03, 0x00, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x19, 0x3a, 0x01, 0x2a, 0x22, 0x14, 0x2f, 0x61,
	0x70, 0x69, 0x2f, 0x76, 0x31, 0x2f, 0x61, 0x67, 0x65, 0x6e, 0x74, 0x2f, 0x64, 0x65, 0x6c, 0x65,
	0x74, 0x65, 0x12, 0x6a, 0x0a, 0x0b, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x41, 0x67, 0x65, 0x6e,
	0x74, 0x12, 0x19, 0x2e, 0x61, 0x67, 0x65, 0x6e, 0x74, 0x2e, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65,
	0x41, 0x67, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1a, 0x2e, 0x61,
	0x67, 0x65, 0x6e, 0x74, 0x2e, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x41, 0x67, 0x65, 0x6e, 0x74,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x24, 0xc2, 0xdb, 0xaa, 0x03, 0x00, 0x82,
	0xd3, 0xe4, 0x93, 0x02, 0x19, 0x3a, 0x01, 0x2a, 0x22, 0x14, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x76,
	0x31, 0x2f, 0x61, 0x67, 0x65, 0x6e, 0x74, 0x2f, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x42, 0x11,
	0x5a, 0x0f, 0x63, 0x6f, 0x64, 0x6f, 0x2d, 0x63, 0x6e, 0x6d, 0x70, 0x2f, 0x70, 0x62, 0x3b, 0x70,
	0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_pb_agent_v1_proto_rawDescOnce sync.Once
	file_pb_agent_v1_proto_rawDescData = file_pb_agent_v1_proto_rawDesc
)

func file_pb_agent_v1_proto_rawDescGZIP() []byte {
	file_pb_agent_v1_proto_rawDescOnce.Do(func() {
		file_pb_agent_v1_proto_rawDescData = protoimpl.X.CompressGZIP(file_pb_agent_v1_proto_rawDescData)
	})
	return file_pb_agent_v1_proto_rawDescData
}

var file_pb_agent_v1_proto_msgTypes = make([]protoimpl.MessageInfo, 9)
var file_pb_agent_v1_proto_goTypes = []interface{}{
	(*ListAgentRequest)(nil),    // 0: agent.ListAgentRequest
	(*CreateAgentRequest)(nil),  // 1: agent.CreateAgentRequest
	(*CreateAgentResponse)(nil), // 2: agent.CreateAgentResponse
	(*DeleteAgentRequest)(nil),  // 3: agent.DeleteAgentRequest
	(*DeleteAgentResponse)(nil), // 4: agent.DeleteAgentResponse
	(*UpdateAgentRequest)(nil),  // 5: agent.UpdateAgentRequest
	(*UpdateAgentResponse)(nil), // 6: agent.UpdateAgentResponse
	(*AgentItem)(nil),           // 7: agent.AgentItem
	(*ListAgentResponse)(nil),   // 8: agent.ListAgentResponse
}
var file_pb_agent_v1_proto_depIdxs = []int32{
	7, // 0: agent.ListAgentResponse.list:type_name -> agent.AgentItem
	0, // 1: agent.Agent.ListAgent:input_type -> agent.ListAgentRequest
	1, // 2: agent.Agent.CreateAgent:input_type -> agent.CreateAgentRequest
	3, // 3: agent.Agent.DeleteAgent:input_type -> agent.DeleteAgentRequest
	5, // 4: agent.Agent.UpdateAgent:input_type -> agent.UpdateAgentRequest
	8, // 5: agent.Agent.ListAgent:output_type -> agent.ListAgentResponse
	2, // 6: agent.Agent.CreateAgent:output_type -> agent.CreateAgentResponse
	4, // 7: agent.Agent.DeleteAgent:output_type -> agent.DeleteAgentResponse
	6, // 8: agent.Agent.UpdateAgent:output_type -> agent.UpdateAgentResponse
	5, // [5:9] is the sub-list for method output_type
	1, // [1:5] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_pb_agent_v1_proto_init() }
func file_pb_agent_v1_proto_init() {
	if File_pb_agent_v1_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_pb_agent_v1_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ListAgentRequest); i {
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
		file_pb_agent_v1_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CreateAgentRequest); i {
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
		file_pb_agent_v1_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CreateAgentResponse); i {
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
		file_pb_agent_v1_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DeleteAgentRequest); i {
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
		file_pb_agent_v1_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DeleteAgentResponse); i {
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
		file_pb_agent_v1_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UpdateAgentRequest); i {
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
		file_pb_agent_v1_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UpdateAgentResponse); i {
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
		file_pb_agent_v1_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AgentItem); i {
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
		file_pb_agent_v1_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ListAgentResponse); i {
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
			RawDescriptor: file_pb_agent_v1_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   9,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_pb_agent_v1_proto_goTypes,
		DependencyIndexes: file_pb_agent_v1_proto_depIdxs,
		MessageInfos:      file_pb_agent_v1_proto_msgTypes,
	}.Build()
	File_pb_agent_v1_proto = out.File
	file_pb_agent_v1_proto_rawDesc = nil
	file_pb_agent_v1_proto_goTypes = nil
	file_pb_agent_v1_proto_depIdxs = nil
}
