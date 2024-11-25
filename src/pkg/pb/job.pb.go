// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.2
// 	protoc        v5.28.2
// source: pkg/pb/job.proto

package pb

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	durationpb "google.golang.org/protobuf/types/known/durationpb"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type Job struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	AgentId string                 `protobuf:"bytes,1,opt,name=agent_id,json=agentId,proto3" json:"agent_id,omitempty"`
	Id      string                 `protobuf:"bytes,2,opt,name=id,proto3" json:"id,omitempty"`
	Time    *timestamppb.Timestamp `protobuf:"bytes,3,opt,name=time,proto3" json:"time,omitempty"`
	Drift   *durationpb.Duration   `protobuf:"bytes,4,opt,name=drift,proto3" json:"drift,omitempty"`
	Took    *durationpb.Duration   `protobuf:"bytes,5,opt,name=took,proto3" json:"took,omitempty"`
	Name    string                 `protobuf:"bytes,6,opt,name=name,proto3" json:"name,omitempty"`
	Errors  uint32                 `protobuf:"varint,7,opt,name=errors,proto3" json:"errors,omitempty"`
}

func (x *Job) Reset() {
	*x = Job{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_pb_job_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Job) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Job) ProtoMessage() {}

func (x *Job) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_pb_job_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Job.ProtoReflect.Descriptor instead.
func (*Job) Descriptor() ([]byte, []int) {
	return file_pkg_pb_job_proto_rawDescGZIP(), []int{0}
}

func (x *Job) GetAgentId() string {
	if x != nil {
		return x.AgentId
	}
	return ""
}

func (x *Job) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Job) GetTime() *timestamppb.Timestamp {
	if x != nil {
		return x.Time
	}
	return nil
}

func (x *Job) GetDrift() *durationpb.Duration {
	if x != nil {
		return x.Drift
	}
	return nil
}

func (x *Job) GetTook() *durationpb.Duration {
	if x != nil {
		return x.Took
	}
	return nil
}

func (x *Job) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Job) GetErrors() uint32 {
	if x != nil {
		return x.Errors
	}
	return 0
}

var File_pkg_pb_job_proto protoreflect.FileDescriptor

var file_pkg_pb_job_proto_rawDesc = []byte{
	0x0a, 0x10, 0x70, 0x6b, 0x67, 0x2f, 0x70, 0x62, 0x2f, 0x6a, 0x6f, 0x62, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x1a, 0x1e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x62, 0x75, 0x66, 0x2f, 0x64, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x22, 0xec, 0x01, 0x0a, 0x03, 0x4a, 0x6f, 0x62, 0x12, 0x19, 0x0a, 0x08, 0x61,
	0x67, 0x65, 0x6e, 0x74, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x61,
	0x67, 0x65, 0x6e, 0x74, 0x49, 0x64, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x2e, 0x0a, 0x04, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70,
	0x52, 0x04, 0x74, 0x69, 0x6d, 0x65, 0x12, 0x2f, 0x0a, 0x05, 0x64, 0x72, 0x69, 0x66, 0x74, 0x18,
	0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x19, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x44, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x52, 0x05, 0x64, 0x72, 0x69, 0x66, 0x74, 0x12, 0x2d, 0x0a, 0x04, 0x74, 0x6f, 0x6f, 0x6b, 0x18,
	0x05, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x19, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x44, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x52, 0x04, 0x74, 0x6f, 0x6f, 0x6b, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x06,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x65, 0x72,
	0x72, 0x6f, 0x72, 0x73, 0x18, 0x07, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x06, 0x65, 0x72, 0x72, 0x6f,
	0x72, 0x73, 0x42, 0x0f, 0x5a, 0x0d, 0x63, 0x73, 0x70, 0x61, 0x67, 0x65, 0x2f, 0x70, 0x6b, 0x67,
	0x2f, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_pkg_pb_job_proto_rawDescOnce sync.Once
	file_pkg_pb_job_proto_rawDescData = file_pkg_pb_job_proto_rawDesc
)

func file_pkg_pb_job_proto_rawDescGZIP() []byte {
	file_pkg_pb_job_proto_rawDescOnce.Do(func() {
		file_pkg_pb_job_proto_rawDescData = protoimpl.X.CompressGZIP(file_pkg_pb_job_proto_rawDescData)
	})
	return file_pkg_pb_job_proto_rawDescData
}

var file_pkg_pb_job_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_pkg_pb_job_proto_goTypes = []any{
	(*Job)(nil),                   // 0: Job
	(*timestamppb.Timestamp)(nil), // 1: google.protobuf.Timestamp
	(*durationpb.Duration)(nil),   // 2: google.protobuf.Duration
}
var file_pkg_pb_job_proto_depIdxs = []int32{
	1, // 0: Job.time:type_name -> google.protobuf.Timestamp
	2, // 1: Job.drift:type_name -> google.protobuf.Duration
	2, // 2: Job.took:type_name -> google.protobuf.Duration
	3, // [3:3] is the sub-list for method output_type
	3, // [3:3] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_pkg_pb_job_proto_init() }
func file_pkg_pb_job_proto_init() {
	if File_pkg_pb_job_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_pkg_pb_job_proto_msgTypes[0].Exporter = func(v any, i int) any {
			switch v := v.(*Job); i {
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
			RawDescriptor: file_pkg_pb_job_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_pkg_pb_job_proto_goTypes,
		DependencyIndexes: file_pkg_pb_job_proto_depIdxs,
		MessageInfos:      file_pkg_pb_job_proto_msgTypes,
	}.Build()
	File_pkg_pb_job_proto = out.File
	file_pkg_pb_job_proto_rawDesc = nil
	file_pkg_pb_job_proto_goTypes = nil
	file_pkg_pb_job_proto_depIdxs = nil
}
