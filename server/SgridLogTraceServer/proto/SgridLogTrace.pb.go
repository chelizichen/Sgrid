// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v4.25.1
// source: SgridLogTrace.proto

package protocol

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type LogTraceReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	LogServerName string `protobuf:"bytes,1,opt,name=logServerName,proto3" json:"logServerName,omitempty"` // 日志服务名
	LogHost       string `protobuf:"bytes,2,opt,name=logHost,proto3" json:"logHost,omitempty"`             // 日志服务所在主机
	LogGridId     int64  `protobuf:"varint,3,opt,name=logGridId,proto3" json:"logGridId,omitempty"`        // 日志服务所在网格id
	LogType       string `protobuf:"bytes,4,opt,name=logType,proto3" json:"logType,omitempty"`             // 日志类型
	LogContent    string `protobuf:"bytes,5,opt,name=logContent,proto3" json:"logContent,omitempty"`       // 日志内容
	CreateTime    string `protobuf:"bytes,6,opt,name=createTime,proto3" json:"createTime,omitempty"`       // 日志创建时间
	LogBytesLen   int64  `protobuf:"varint,7,opt,name=logBytesLen,proto3" json:"logBytesLen,omitempty"`    // 日志字节长度
}

func (x *LogTraceReq) Reset() {
	*x = LogTraceReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_SgridLogTrace_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LogTraceReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LogTraceReq) ProtoMessage() {}

func (x *LogTraceReq) ProtoReflect() protoreflect.Message {
	mi := &file_SgridLogTrace_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LogTraceReq.ProtoReflect.Descriptor instead.
func (*LogTraceReq) Descriptor() ([]byte, []int) {
	return file_SgridLogTrace_proto_rawDescGZIP(), []int{0}
}

func (x *LogTraceReq) GetLogServerName() string {
	if x != nil {
		return x.LogServerName
	}
	return ""
}

func (x *LogTraceReq) GetLogHost() string {
	if x != nil {
		return x.LogHost
	}
	return ""
}

func (x *LogTraceReq) GetLogGridId() int64 {
	if x != nil {
		return x.LogGridId
	}
	return 0
}

func (x *LogTraceReq) GetLogType() string {
	if x != nil {
		return x.LogType
	}
	return ""
}

func (x *LogTraceReq) GetLogContent() string {
	if x != nil {
		return x.LogContent
	}
	return ""
}

func (x *LogTraceReq) GetCreateTime() string {
	if x != nil {
		return x.CreateTime
	}
	return ""
}

func (x *LogTraceReq) GetLogBytesLen() int64 {
	if x != nil {
		return x.LogBytesLen
	}
	return 0
}

type BasicResp struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Code    int32  `protobuf:"varint,1,opt,name=code,proto3" json:"code,omitempty"`
	Message string `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"`
}

func (x *BasicResp) Reset() {
	*x = BasicResp{}
	if protoimpl.UnsafeEnabled {
		mi := &file_SgridLogTrace_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BasicResp) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BasicResp) ProtoMessage() {}

func (x *BasicResp) ProtoReflect() protoreflect.Message {
	mi := &file_SgridLogTrace_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BasicResp.ProtoReflect.Descriptor instead.
func (*BasicResp) Descriptor() ([]byte, []int) {
	return file_SgridLogTrace_proto_rawDescGZIP(), []int{1}
}

func (x *BasicResp) GetCode() int32 {
	if x != nil {
		return x.Code
	}
	return 0
}

func (x *BasicResp) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

var File_SgridLogTrace_proto protoreflect.FileDescriptor

var file_SgridLogTrace_proto_rawDesc = []byte{
	0x0a, 0x13, 0x53, 0x67, 0x72, 0x69, 0x64, 0x4c, 0x6f, 0x67, 0x54, 0x72, 0x61, 0x63, 0x65, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0d, 0x53, 0x67, 0x72, 0x69, 0x64, 0x4c, 0x6f, 0x67, 0x54,
	0x72, 0x61, 0x63, 0x65, 0x1a, 0x1b, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x65, 0x6d, 0x70, 0x74, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x22, 0xe7, 0x01, 0x0a, 0x0b, 0x4c, 0x6f, 0x67, 0x54, 0x72, 0x61, 0x63, 0x65, 0x52, 0x65,
	0x71, 0x12, 0x24, 0x0a, 0x0d, 0x6c, 0x6f, 0x67, 0x53, 0x65, 0x72, 0x76, 0x65, 0x72, 0x4e, 0x61,
	0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0d, 0x6c, 0x6f, 0x67, 0x53, 0x65, 0x72,
	0x76, 0x65, 0x72, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x6c, 0x6f, 0x67, 0x48, 0x6f,
	0x73, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x6c, 0x6f, 0x67, 0x48, 0x6f, 0x73,
	0x74, 0x12, 0x1c, 0x0a, 0x09, 0x6c, 0x6f, 0x67, 0x47, 0x72, 0x69, 0x64, 0x49, 0x64, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x03, 0x52, 0x09, 0x6c, 0x6f, 0x67, 0x47, 0x72, 0x69, 0x64, 0x49, 0x64, 0x12,
	0x18, 0x0a, 0x07, 0x6c, 0x6f, 0x67, 0x54, 0x79, 0x70, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x07, 0x6c, 0x6f, 0x67, 0x54, 0x79, 0x70, 0x65, 0x12, 0x1e, 0x0a, 0x0a, 0x6c, 0x6f, 0x67,
	0x43, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x6c,
	0x6f, 0x67, 0x43, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x12, 0x1e, 0x0a, 0x0a, 0x63, 0x72, 0x65,
	0x61, 0x74, 0x65, 0x54, 0x69, 0x6d, 0x65, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x63,
	0x72, 0x65, 0x61, 0x74, 0x65, 0x54, 0x69, 0x6d, 0x65, 0x12, 0x20, 0x0a, 0x0b, 0x6c, 0x6f, 0x67,
	0x42, 0x79, 0x74, 0x65, 0x73, 0x4c, 0x65, 0x6e, 0x18, 0x07, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0b,
	0x6c, 0x6f, 0x67, 0x42, 0x79, 0x74, 0x65, 0x73, 0x4c, 0x65, 0x6e, 0x22, 0x39, 0x0a, 0x09, 0x62,
	0x61, 0x73, 0x69, 0x63, 0x52, 0x65, 0x73, 0x70, 0x12, 0x12, 0x0a, 0x04, 0x63, 0x6f, 0x64, 0x65,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x12, 0x18, 0x0a, 0x07,
	0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x6d,
	0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x32, 0x56, 0x0a, 0x14, 0x53, 0x67, 0x72, 0x69, 0x64, 0x4c,
	0x6f, 0x67, 0x54, 0x72, 0x61, 0x63, 0x65, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x3e,
	0x0a, 0x08, 0x4c, 0x6f, 0x67, 0x54, 0x72, 0x61, 0x63, 0x65, 0x12, 0x1a, 0x2e, 0x53, 0x67, 0x72,
	0x69, 0x64, 0x4c, 0x6f, 0x67, 0x54, 0x72, 0x61, 0x63, 0x65, 0x2e, 0x4c, 0x6f, 0x67, 0x54, 0x72,
	0x61, 0x63, 0x65, 0x52, 0x65, 0x71, 0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x42, 0x0c,
	0x5a, 0x0a, 0x2e, 0x3b, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6c, 0x62, 0x06, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_SgridLogTrace_proto_rawDescOnce sync.Once
	file_SgridLogTrace_proto_rawDescData = file_SgridLogTrace_proto_rawDesc
)

func file_SgridLogTrace_proto_rawDescGZIP() []byte {
	file_SgridLogTrace_proto_rawDescOnce.Do(func() {
		file_SgridLogTrace_proto_rawDescData = protoimpl.X.CompressGZIP(file_SgridLogTrace_proto_rawDescData)
	})
	return file_SgridLogTrace_proto_rawDescData
}

var file_SgridLogTrace_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_SgridLogTrace_proto_goTypes = []interface{}{
	(*LogTraceReq)(nil),   // 0: SgridLogTrace.LogTraceReq
	(*BasicResp)(nil),     // 1: SgridLogTrace.basicResp
	(*emptypb.Empty)(nil), // 2: google.protobuf.Empty
}
var file_SgridLogTrace_proto_depIdxs = []int32{
	0, // 0: SgridLogTrace.SgridLogTraceService.LogTrace:input_type -> SgridLogTrace.LogTraceReq
	2, // 1: SgridLogTrace.SgridLogTraceService.LogTrace:output_type -> google.protobuf.Empty
	1, // [1:2] is the sub-list for method output_type
	0, // [0:1] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_SgridLogTrace_proto_init() }
func file_SgridLogTrace_proto_init() {
	if File_SgridLogTrace_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_SgridLogTrace_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*LogTraceReq); i {
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
		file_SgridLogTrace_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BasicResp); i {
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
			RawDescriptor: file_SgridLogTrace_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_SgridLogTrace_proto_goTypes,
		DependencyIndexes: file_SgridLogTrace_proto_depIdxs,
		MessageInfos:      file_SgridLogTrace_proto_msgTypes,
	}.Build()
	File_SgridLogTrace_proto = out.File
	file_SgridLogTrace_proto_rawDesc = nil
	file_SgridLogTrace_proto_goTypes = nil
	file_SgridLogTrace_proto_depIdxs = nil
}
