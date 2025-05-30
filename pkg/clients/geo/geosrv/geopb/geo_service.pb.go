// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.2
// 	protoc        v3.6.1
// source: api/proto/geo_service.proto

package geopb

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

// Request
type GetGeolocationRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Street string `protobuf:"bytes,1,opt,name=Street,proto3" json:"Street,omitempty"`
}

func (x *GetGeolocationRequest) Reset() {
	*x = GetGeolocationRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_proto_geo_service_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetGeolocationRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetGeolocationRequest) ProtoMessage() {}

func (x *GetGeolocationRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_proto_geo_service_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetGeolocationRequest.ProtoReflect.Descriptor instead.
func (*GetGeolocationRequest) Descriptor() ([]byte, []int) {
	return file_api_proto_geo_service_proto_rawDescGZIP(), []int{0}
}

func (x *GetGeolocationRequest) GetStreet() string {
	if x != nil {
		return x.Street
	}
	return ""
}

// Response
type GetGeolocationReply struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Location *Location `protobuf:"bytes,1,opt,name=Location,proto3" json:"Location,omitempty"`
}

func (x *GetGeolocationReply) Reset() {
	*x = GetGeolocationReply{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_proto_geo_service_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetGeolocationReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetGeolocationReply) ProtoMessage() {}

func (x *GetGeolocationReply) ProtoReflect() protoreflect.Message {
	mi := &file_api_proto_geo_service_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetGeolocationReply.ProtoReflect.Descriptor instead.
func (*GetGeolocationReply) Descriptor() ([]byte, []int) {
	return file_api_proto_geo_service_proto_rawDescGZIP(), []int{1}
}

func (x *GetGeolocationReply) GetLocation() *Location {
	if x != nil {
		return x.Location
	}
	return nil
}

// Geolocation
type Location struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	X int32 `protobuf:"varint,1,opt,name=x,proto3" json:"x,omitempty"`
	Y int32 `protobuf:"varint,2,opt,name=y,proto3" json:"y,omitempty"`
}

func (x *Location) Reset() {
	*x = Location{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_proto_geo_service_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Location) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Location) ProtoMessage() {}

func (x *Location) ProtoReflect() protoreflect.Message {
	mi := &file_api_proto_geo_service_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Location.ProtoReflect.Descriptor instead.
func (*Location) Descriptor() ([]byte, []int) {
	return file_api_proto_geo_service_proto_rawDescGZIP(), []int{2}
}

func (x *Location) GetX() int32 {
	if x != nil {
		return x.X
	}
	return 0
}

func (x *Location) GetY() int32 {
	if x != nil {
		return x.Y
	}
	return 0
}

type ErrorResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Text string `protobuf:"bytes,1,opt,name=text,proto3" json:"text,omitempty"`
}

func (x *ErrorResponse) Reset() {
	*x = ErrorResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_proto_geo_service_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ErrorResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ErrorResponse) ProtoMessage() {}

func (x *ErrorResponse) ProtoReflect() protoreflect.Message {
	mi := &file_api_proto_geo_service_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ErrorResponse.ProtoReflect.Descriptor instead.
func (*ErrorResponse) Descriptor() ([]byte, []int) {
	return file_api_proto_geo_service_proto_rawDescGZIP(), []int{3}
}

func (x *ErrorResponse) GetText() string {
	if x != nil {
		return x.Text
	}
	return ""
}

var File_api_proto_geo_service_proto protoreflect.FileDescriptor

var file_api_proto_geo_service_proto_rawDesc = []byte{
	0x0a, 0x1b, 0x61, 0x70, 0x69, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x67, 0x65, 0x6f, 0x5f,
	0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x03, 0x67,
	0x65, 0x6f, 0x22, 0x2f, 0x0a, 0x15, 0x47, 0x65, 0x74, 0x47, 0x65, 0x6f, 0x6c, 0x6f, 0x63, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x16, 0x0a, 0x06, 0x53,
	0x74, 0x72, 0x65, 0x65, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x53, 0x74, 0x72,
	0x65, 0x65, 0x74, 0x22, 0x40, 0x0a, 0x13, 0x47, 0x65, 0x74, 0x47, 0x65, 0x6f, 0x6c, 0x6f, 0x63,
	0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x12, 0x29, 0x0a, 0x08, 0x4c, 0x6f,
	0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0d, 0x2e, 0x67,
	0x65, 0x6f, 0x2e, 0x4c, 0x6f, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x08, 0x4c, 0x6f, 0x63,
	0x61, 0x74, 0x69, 0x6f, 0x6e, 0x22, 0x26, 0x0a, 0x08, 0x4c, 0x6f, 0x63, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x12, 0x0c, 0x0a, 0x01, 0x78, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x01, 0x78, 0x12,
	0x0c, 0x0a, 0x01, 0x79, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x01, 0x79, 0x22, 0x23, 0x0a,
	0x0d, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x12,
	0x0a, 0x04, 0x74, 0x65, 0x78, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x74, 0x65,
	0x78, 0x74, 0x32, 0x4d, 0x0a, 0x03, 0x47, 0x65, 0x6f, 0x12, 0x46, 0x0a, 0x0e, 0x47, 0x65, 0x74,
	0x47, 0x65, 0x6f, 0x6c, 0x6f, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x1a, 0x2e, 0x67, 0x65,
	0x6f, 0x2e, 0x47, 0x65, 0x74, 0x47, 0x65, 0x6f, 0x6c, 0x6f, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x18, 0x2e, 0x67, 0x65, 0x6f, 0x2e, 0x47, 0x65,
	0x74, 0x47, 0x65, 0x6f, 0x6c, 0x6f, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x70, 0x6c,
	0x79, 0x42, 0x1b, 0x5a, 0x0c, 0x67, 0x65, 0x6f, 0x73, 0x72, 0x76, 0x2f, 0x67, 0x65, 0x6f, 0x70,
	0x62, 0xaa, 0x02, 0x0a, 0x47, 0x65, 0x6f, 0x41, 0x70, 0x70, 0x2e, 0x41, 0x70, 0x69, 0x62, 0x06,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_api_proto_geo_service_proto_rawDescOnce sync.Once
	file_api_proto_geo_service_proto_rawDescData = file_api_proto_geo_service_proto_rawDesc
)

func file_api_proto_geo_service_proto_rawDescGZIP() []byte {
	file_api_proto_geo_service_proto_rawDescOnce.Do(func() {
		file_api_proto_geo_service_proto_rawDescData = protoimpl.X.CompressGZIP(file_api_proto_geo_service_proto_rawDescData)
	})
	return file_api_proto_geo_service_proto_rawDescData
}

var file_api_proto_geo_service_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_api_proto_geo_service_proto_goTypes = []any{
	(*GetGeolocationRequest)(nil), // 0: geo.GetGeolocationRequest
	(*GetGeolocationReply)(nil),   // 1: geo.GetGeolocationReply
	(*Location)(nil),              // 2: geo.Location
	(*ErrorResponse)(nil),         // 3: geo.ErrorResponse
}
var file_api_proto_geo_service_proto_depIdxs = []int32{
	2, // 0: geo.GetGeolocationReply.Location:type_name -> geo.Location
	0, // 1: geo.Geo.GetGeolocation:input_type -> geo.GetGeolocationRequest
	1, // 2: geo.Geo.GetGeolocation:output_type -> geo.GetGeolocationReply
	2, // [2:3] is the sub-list for method output_type
	1, // [1:2] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_api_proto_geo_service_proto_init() }
func file_api_proto_geo_service_proto_init() {
	if File_api_proto_geo_service_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_api_proto_geo_service_proto_msgTypes[0].Exporter = func(v any, i int) any {
			switch v := v.(*GetGeolocationRequest); i {
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
		file_api_proto_geo_service_proto_msgTypes[1].Exporter = func(v any, i int) any {
			switch v := v.(*GetGeolocationReply); i {
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
		file_api_proto_geo_service_proto_msgTypes[2].Exporter = func(v any, i int) any {
			switch v := v.(*Location); i {
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
		file_api_proto_geo_service_proto_msgTypes[3].Exporter = func(v any, i int) any {
			switch v := v.(*ErrorResponse); i {
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
			RawDescriptor: file_api_proto_geo_service_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_api_proto_geo_service_proto_goTypes,
		DependencyIndexes: file_api_proto_geo_service_proto_depIdxs,
		MessageInfos:      file_api_proto_geo_service_proto_msgTypes,
	}.Build()
	File_api_proto_geo_service_proto = out.File
	file_api_proto_geo_service_proto_rawDesc = nil
	file_api_proto_geo_service_proto_goTypes = nil
	file_api_proto_geo_service_proto_depIdxs = nil
}
