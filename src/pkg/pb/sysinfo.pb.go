// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.2
// 	protoc        v5.28.2
// source: pkg/pb/sysinfo.proto

package pb

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

type OSStat struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Mem *OSStat_Memory `protobuf:"bytes,1,opt,name=mem,proto3" json:"mem,omitempty"`
	Cpu *OSStat_CPU    `protobuf:"bytes,2,opt,name=cpu,proto3" json:"cpu,omitempty"`
}

func (x *OSStat) Reset() {
	*x = OSStat{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_pb_sysinfo_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *OSStat) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*OSStat) ProtoMessage() {}

func (x *OSStat) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_pb_sysinfo_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use OSStat.ProtoReflect.Descriptor instead.
func (*OSStat) Descriptor() ([]byte, []int) {
	return file_pkg_pb_sysinfo_proto_rawDescGZIP(), []int{0}
}

func (x *OSStat) GetMem() *OSStat_Memory {
	if x != nil {
		return x.Mem
	}
	return nil
}

func (x *OSStat) GetCpu() *OSStat_CPU {
	if x != nil {
		return x.Cpu
	}
	return nil
}

type ProcStat struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Threads    int32            `protobuf:"varint,1,opt,name=threads,proto3" json:"threads,omitempty"`
	Fds        int32            `protobuf:"varint,2,opt,name=fds,proto3" json:"fds,omitempty"`
	CpuPercent float32          `protobuf:"fixed32,3,opt,name=cpu_percent,json=cpuPercent,proto3" json:"cpu_percent,omitempty"`
	Mem        *ProcStat_Memory `protobuf:"bytes,4,opt,name=mem,proto3" json:"mem,omitempty"`
	Io         *ProcStat_IO     `protobuf:"bytes,5,opt,name=io,proto3" json:"io,omitempty"`
}

func (x *ProcStat) Reset() {
	*x = ProcStat{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_pb_sysinfo_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ProcStat) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ProcStat) ProtoMessage() {}

func (x *ProcStat) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_pb_sysinfo_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ProcStat.ProtoReflect.Descriptor instead.
func (*ProcStat) Descriptor() ([]byte, []int) {
	return file_pkg_pb_sysinfo_proto_rawDescGZIP(), []int{1}
}

func (x *ProcStat) GetThreads() int32 {
	if x != nil {
		return x.Threads
	}
	return 0
}

func (x *ProcStat) GetFds() int32 {
	if x != nil {
		return x.Fds
	}
	return 0
}

func (x *ProcStat) GetCpuPercent() float32 {
	if x != nil {
		return x.CpuPercent
	}
	return 0
}

func (x *ProcStat) GetMem() *ProcStat_Memory {
	if x != nil {
		return x.Mem
	}
	return nil
}

func (x *ProcStat) GetIo() *ProcStat_IO {
	if x != nil {
		return x.Io
	}
	return nil
}

type SysStat struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Os   *OSStat   `protobuf:"bytes,1,opt,name=os,proto3" json:"os,omitempty"`
	Proc *ProcStat `protobuf:"bytes,2,opt,name=proc,proto3" json:"proc,omitempty"`
}

func (x *SysStat) Reset() {
	*x = SysStat{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_pb_sysinfo_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SysStat) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SysStat) ProtoMessage() {}

func (x *SysStat) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_pb_sysinfo_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SysStat.ProtoReflect.Descriptor instead.
func (*SysStat) Descriptor() ([]byte, []int) {
	return file_pkg_pb_sysinfo_proto_rawDescGZIP(), []int{2}
}

func (x *SysStat) GetOs() *OSStat {
	if x != nil {
		return x.Os
	}
	return nil
}

func (x *SysStat) GetProc() *ProcStat {
	if x != nil {
		return x.Proc
	}
	return nil
}

type OSStat_Memory struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Total        uint64 `protobuf:"varint,1,opt,name=total,proto3" json:"total,omitempty"`
	Available    uint64 `protobuf:"varint,2,opt,name=available,proto3" json:"available,omitempty"`
	Used         uint64 `protobuf:"varint,3,opt,name=used,proto3" json:"used,omitempty"`
	Free         uint64 `protobuf:"varint,4,opt,name=free,proto3" json:"free,omitempty"`
	Active       uint64 `protobuf:"varint,5,opt,name=active,proto3" json:"active,omitempty"`
	Inactive     uint64 `protobuf:"varint,6,opt,name=inactive,proto3" json:"inactive,omitempty"`
	Wired        uint64 `protobuf:"varint,7,opt,name=wired,proto3" json:"wired,omitempty"`
	Laundry      uint64 `protobuf:"varint,8,opt,name=laundry,proto3" json:"laundry,omitempty"`
	Buffers      uint64 `protobuf:"varint,9,opt,name=buffers,proto3" json:"buffers,omitempty"`
	Cached       uint64 `protobuf:"varint,10,opt,name=cached,proto3" json:"cached,omitempty"`
	WriteBack    uint64 `protobuf:"varint,11,opt,name=write_back,json=writeBack,proto3" json:"write_back,omitempty"`
	Dirty        uint64 `protobuf:"varint,12,opt,name=dirty,proto3" json:"dirty,omitempty"`
	WriteBackTmp uint64 `protobuf:"varint,13,opt,name=write_back_tmp,json=writeBackTmp,proto3" json:"write_back_tmp,omitempty"`
	Shared       uint64 `protobuf:"varint,14,opt,name=shared,proto3" json:"shared,omitempty"`
	Slab         uint64 `protobuf:"varint,15,opt,name=slab,proto3" json:"slab,omitempty"`
}

func (x *OSStat_Memory) Reset() {
	*x = OSStat_Memory{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_pb_sysinfo_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *OSStat_Memory) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*OSStat_Memory) ProtoMessage() {}

func (x *OSStat_Memory) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_pb_sysinfo_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use OSStat_Memory.ProtoReflect.Descriptor instead.
func (*OSStat_Memory) Descriptor() ([]byte, []int) {
	return file_pkg_pb_sysinfo_proto_rawDescGZIP(), []int{0, 0}
}

func (x *OSStat_Memory) GetTotal() uint64 {
	if x != nil {
		return x.Total
	}
	return 0
}

func (x *OSStat_Memory) GetAvailable() uint64 {
	if x != nil {
		return x.Available
	}
	return 0
}

func (x *OSStat_Memory) GetUsed() uint64 {
	if x != nil {
		return x.Used
	}
	return 0
}

func (x *OSStat_Memory) GetFree() uint64 {
	if x != nil {
		return x.Free
	}
	return 0
}

func (x *OSStat_Memory) GetActive() uint64 {
	if x != nil {
		return x.Active
	}
	return 0
}

func (x *OSStat_Memory) GetInactive() uint64 {
	if x != nil {
		return x.Inactive
	}
	return 0
}

func (x *OSStat_Memory) GetWired() uint64 {
	if x != nil {
		return x.Wired
	}
	return 0
}

func (x *OSStat_Memory) GetLaundry() uint64 {
	if x != nil {
		return x.Laundry
	}
	return 0
}

func (x *OSStat_Memory) GetBuffers() uint64 {
	if x != nil {
		return x.Buffers
	}
	return 0
}

func (x *OSStat_Memory) GetCached() uint64 {
	if x != nil {
		return x.Cached
	}
	return 0
}

func (x *OSStat_Memory) GetWriteBack() uint64 {
	if x != nil {
		return x.WriteBack
	}
	return 0
}

func (x *OSStat_Memory) GetDirty() uint64 {
	if x != nil {
		return x.Dirty
	}
	return 0
}

func (x *OSStat_Memory) GetWriteBackTmp() uint64 {
	if x != nil {
		return x.WriteBackTmp
	}
	return 0
}

func (x *OSStat_Memory) GetShared() uint64 {
	if x != nil {
		return x.Shared
	}
	return 0
}

func (x *OSStat_Memory) GetSlab() uint64 {
	if x != nil {
		return x.Slab
	}
	return 0
}

type OSStat_CPU struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	User    float32 `protobuf:"fixed32,1,opt,name=user,proto3" json:"user,omitempty"`
	System  float32 `protobuf:"fixed32,2,opt,name=system,proto3" json:"system,omitempty"`
	Idle    float32 `protobuf:"fixed32,3,opt,name=idle,proto3" json:"idle,omitempty"`
	Nice    float32 `protobuf:"fixed32,4,opt,name=nice,proto3" json:"nice,omitempty"`
	Iowait  float32 `protobuf:"fixed32,5,opt,name=iowait,proto3" json:"iowait,omitempty"`
	Irq     float32 `protobuf:"fixed32,6,opt,name=irq,proto3" json:"irq,omitempty"`
	Softirq float32 `protobuf:"fixed32,7,opt,name=softirq,proto3" json:"softirq,omitempty"`
	Steal   float32 `protobuf:"fixed32,8,opt,name=steal,proto3" json:"steal,omitempty"`
}

func (x *OSStat_CPU) Reset() {
	*x = OSStat_CPU{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_pb_sysinfo_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *OSStat_CPU) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*OSStat_CPU) ProtoMessage() {}

func (x *OSStat_CPU) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_pb_sysinfo_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use OSStat_CPU.ProtoReflect.Descriptor instead.
func (*OSStat_CPU) Descriptor() ([]byte, []int) {
	return file_pkg_pb_sysinfo_proto_rawDescGZIP(), []int{0, 1}
}

func (x *OSStat_CPU) GetUser() float32 {
	if x != nil {
		return x.User
	}
	return 0
}

func (x *OSStat_CPU) GetSystem() float32 {
	if x != nil {
		return x.System
	}
	return 0
}

func (x *OSStat_CPU) GetIdle() float32 {
	if x != nil {
		return x.Idle
	}
	return 0
}

func (x *OSStat_CPU) GetNice() float32 {
	if x != nil {
		return x.Nice
	}
	return 0
}

func (x *OSStat_CPU) GetIowait() float32 {
	if x != nil {
		return x.Iowait
	}
	return 0
}

func (x *OSStat_CPU) GetIrq() float32 {
	if x != nil {
		return x.Irq
	}
	return 0
}

func (x *OSStat_CPU) GetSoftirq() float32 {
	if x != nil {
		return x.Softirq
	}
	return 0
}

func (x *OSStat_CPU) GetSteal() float32 {
	if x != nil {
		return x.Steal
	}
	return 0
}

type ProcStat_Memory struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Rss    uint64 `protobuf:"varint,1,opt,name=rss,proto3" json:"rss,omitempty"`
	Vms    uint64 `protobuf:"varint,2,opt,name=vms,proto3" json:"vms,omitempty"`
	Hwm    uint64 `protobuf:"varint,3,opt,name=hwm,proto3" json:"hwm,omitempty"`
	Data   uint64 `protobuf:"varint,4,opt,name=data,proto3" json:"data,omitempty"`
	Stack  uint64 `protobuf:"varint,5,opt,name=stack,proto3" json:"stack,omitempty"`
	Locked uint64 `protobuf:"varint,6,opt,name=locked,proto3" json:"locked,omitempty"`
	Swap   uint64 `protobuf:"varint,7,opt,name=swap,proto3" json:"swap,omitempty"`
}

func (x *ProcStat_Memory) Reset() {
	*x = ProcStat_Memory{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_pb_sysinfo_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ProcStat_Memory) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ProcStat_Memory) ProtoMessage() {}

func (x *ProcStat_Memory) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_pb_sysinfo_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ProcStat_Memory.ProtoReflect.Descriptor instead.
func (*ProcStat_Memory) Descriptor() ([]byte, []int) {
	return file_pkg_pb_sysinfo_proto_rawDescGZIP(), []int{1, 0}
}

func (x *ProcStat_Memory) GetRss() uint64 {
	if x != nil {
		return x.Rss
	}
	return 0
}

func (x *ProcStat_Memory) GetVms() uint64 {
	if x != nil {
		return x.Vms
	}
	return 0
}

func (x *ProcStat_Memory) GetHwm() uint64 {
	if x != nil {
		return x.Hwm
	}
	return 0
}

func (x *ProcStat_Memory) GetData() uint64 {
	if x != nil {
		return x.Data
	}
	return 0
}

func (x *ProcStat_Memory) GetStack() uint64 {
	if x != nil {
		return x.Stack
	}
	return 0
}

func (x *ProcStat_Memory) GetLocked() uint64 {
	if x != nil {
		return x.Locked
	}
	return 0
}

func (x *ProcStat_Memory) GetSwap() uint64 {
	if x != nil {
		return x.Swap
	}
	return 0
}

type ProcStat_IO struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ReadCount  uint64 `protobuf:"varint,1,opt,name=read_count,json=readCount,proto3" json:"read_count,omitempty"`
	WriteCount uint64 `protobuf:"varint,2,opt,name=write_count,json=writeCount,proto3" json:"write_count,omitempty"`
	ReadBytes  uint64 `protobuf:"varint,3,opt,name=read_bytes,json=readBytes,proto3" json:"read_bytes,omitempty"`
	WriteBytes uint64 `protobuf:"varint,4,opt,name=write_bytes,json=writeBytes,proto3" json:"write_bytes,omitempty"`
}

func (x *ProcStat_IO) Reset() {
	*x = ProcStat_IO{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_pb_sysinfo_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ProcStat_IO) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ProcStat_IO) ProtoMessage() {}

func (x *ProcStat_IO) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_pb_sysinfo_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ProcStat_IO.ProtoReflect.Descriptor instead.
func (*ProcStat_IO) Descriptor() ([]byte, []int) {
	return file_pkg_pb_sysinfo_proto_rawDescGZIP(), []int{1, 1}
}

func (x *ProcStat_IO) GetReadCount() uint64 {
	if x != nil {
		return x.ReadCount
	}
	return 0
}

func (x *ProcStat_IO) GetWriteCount() uint64 {
	if x != nil {
		return x.WriteCount
	}
	return 0
}

func (x *ProcStat_IO) GetReadBytes() uint64 {
	if x != nil {
		return x.ReadBytes
	}
	return 0
}

func (x *ProcStat_IO) GetWriteBytes() uint64 {
	if x != nil {
		return x.WriteBytes
	}
	return 0
}

var File_pkg_pb_sysinfo_proto protoreflect.FileDescriptor

var file_pkg_pb_sysinfo_proto_rawDesc = []byte{
	0x0a, 0x14, 0x70, 0x6b, 0x67, 0x2f, 0x70, 0x62, 0x2f, 0x73, 0x79, 0x73, 0x69, 0x6e, 0x66, 0x6f,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x83, 0x05, 0x0a, 0x06, 0x4f, 0x53, 0x53, 0x74, 0x61,
	0x74, 0x12, 0x20, 0x0a, 0x03, 0x6d, 0x65, 0x6d, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0e,
	0x2e, 0x4f, 0x53, 0x53, 0x74, 0x61, 0x74, 0x2e, 0x4d, 0x65, 0x6d, 0x6f, 0x72, 0x79, 0x52, 0x03,
	0x6d, 0x65, 0x6d, 0x12, 0x1d, 0x0a, 0x03, 0x63, 0x70, 0x75, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x0b, 0x2e, 0x4f, 0x53, 0x53, 0x74, 0x61, 0x74, 0x2e, 0x43, 0x50, 0x55, 0x52, 0x03, 0x63,
	0x70, 0x75, 0x1a, 0x81, 0x03, 0x0a, 0x06, 0x4d, 0x65, 0x6d, 0x6f, 0x72, 0x79, 0x12, 0x14, 0x0a,
	0x05, 0x74, 0x6f, 0x74, 0x61, 0x6c, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x05, 0x74, 0x6f,
	0x74, 0x61, 0x6c, 0x12, 0x1c, 0x0a, 0x09, 0x61, 0x76, 0x61, 0x69, 0x6c, 0x61, 0x62, 0x6c, 0x65,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x04, 0x52, 0x09, 0x61, 0x76, 0x61, 0x69, 0x6c, 0x61, 0x62, 0x6c,
	0x65, 0x12, 0x12, 0x0a, 0x04, 0x75, 0x73, 0x65, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x04, 0x52,
	0x04, 0x75, 0x73, 0x65, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x66, 0x72, 0x65, 0x65, 0x18, 0x04, 0x20,
	0x01, 0x28, 0x04, 0x52, 0x04, 0x66, 0x72, 0x65, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x61, 0x63, 0x74,
	0x69, 0x76, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x04, 0x52, 0x06, 0x61, 0x63, 0x74, 0x69, 0x76,
	0x65, 0x12, 0x1a, 0x0a, 0x08, 0x69, 0x6e, 0x61, 0x63, 0x74, 0x69, 0x76, 0x65, 0x18, 0x06, 0x20,
	0x01, 0x28, 0x04, 0x52, 0x08, 0x69, 0x6e, 0x61, 0x63, 0x74, 0x69, 0x76, 0x65, 0x12, 0x14, 0x0a,
	0x05, 0x77, 0x69, 0x72, 0x65, 0x64, 0x18, 0x07, 0x20, 0x01, 0x28, 0x04, 0x52, 0x05, 0x77, 0x69,
	0x72, 0x65, 0x64, 0x12, 0x18, 0x0a, 0x07, 0x6c, 0x61, 0x75, 0x6e, 0x64, 0x72, 0x79, 0x18, 0x08,
	0x20, 0x01, 0x28, 0x04, 0x52, 0x07, 0x6c, 0x61, 0x75, 0x6e, 0x64, 0x72, 0x79, 0x12, 0x18, 0x0a,
	0x07, 0x62, 0x75, 0x66, 0x66, 0x65, 0x72, 0x73, 0x18, 0x09, 0x20, 0x01, 0x28, 0x04, 0x52, 0x07,
	0x62, 0x75, 0x66, 0x66, 0x65, 0x72, 0x73, 0x12, 0x16, 0x0a, 0x06, 0x63, 0x61, 0x63, 0x68, 0x65,
	0x64, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x04, 0x52, 0x06, 0x63, 0x61, 0x63, 0x68, 0x65, 0x64, 0x12,
	0x1d, 0x0a, 0x0a, 0x77, 0x72, 0x69, 0x74, 0x65, 0x5f, 0x62, 0x61, 0x63, 0x6b, 0x18, 0x0b, 0x20,
	0x01, 0x28, 0x04, 0x52, 0x09, 0x77, 0x72, 0x69, 0x74, 0x65, 0x42, 0x61, 0x63, 0x6b, 0x12, 0x14,
	0x0a, 0x05, 0x64, 0x69, 0x72, 0x74, 0x79, 0x18, 0x0c, 0x20, 0x01, 0x28, 0x04, 0x52, 0x05, 0x64,
	0x69, 0x72, 0x74, 0x79, 0x12, 0x24, 0x0a, 0x0e, 0x77, 0x72, 0x69, 0x74, 0x65, 0x5f, 0x62, 0x61,
	0x63, 0x6b, 0x5f, 0x74, 0x6d, 0x70, 0x18, 0x0d, 0x20, 0x01, 0x28, 0x04, 0x52, 0x0c, 0x77, 0x72,
	0x69, 0x74, 0x65, 0x42, 0x61, 0x63, 0x6b, 0x54, 0x6d, 0x70, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x68,
	0x61, 0x72, 0x65, 0x64, 0x18, 0x0e, 0x20, 0x01, 0x28, 0x04, 0x52, 0x06, 0x73, 0x68, 0x61, 0x72,
	0x65, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x73, 0x6c, 0x61, 0x62, 0x18, 0x0f, 0x20, 0x01, 0x28, 0x04,
	0x52, 0x04, 0x73, 0x6c, 0x61, 0x62, 0x1a, 0xb3, 0x01, 0x0a, 0x03, 0x43, 0x50, 0x55, 0x12, 0x12,
	0x0a, 0x04, 0x75, 0x73, 0x65, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x02, 0x52, 0x04, 0x75, 0x73,
	0x65, 0x72, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x79, 0x73, 0x74, 0x65, 0x6d, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x02, 0x52, 0x06, 0x73, 0x79, 0x73, 0x74, 0x65, 0x6d, 0x12, 0x12, 0x0a, 0x04, 0x69, 0x64,
	0x6c, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x02, 0x52, 0x04, 0x69, 0x64, 0x6c, 0x65, 0x12, 0x12,
	0x0a, 0x04, 0x6e, 0x69, 0x63, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x02, 0x52, 0x04, 0x6e, 0x69,
	0x63, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x69, 0x6f, 0x77, 0x61, 0x69, 0x74, 0x18, 0x05, 0x20, 0x01,
	0x28, 0x02, 0x52, 0x06, 0x69, 0x6f, 0x77, 0x61, 0x69, 0x74, 0x12, 0x10, 0x0a, 0x03, 0x69, 0x72,
	0x71, 0x18, 0x06, 0x20, 0x01, 0x28, 0x02, 0x52, 0x03, 0x69, 0x72, 0x71, 0x12, 0x18, 0x0a, 0x07,
	0x73, 0x6f, 0x66, 0x74, 0x69, 0x72, 0x71, 0x18, 0x07, 0x20, 0x01, 0x28, 0x02, 0x52, 0x07, 0x73,
	0x6f, 0x66, 0x74, 0x69, 0x72, 0x71, 0x12, 0x14, 0x0a, 0x05, 0x73, 0x74, 0x65, 0x61, 0x6c, 0x18,
	0x08, 0x20, 0x01, 0x28, 0x02, 0x52, 0x05, 0x73, 0x74, 0x65, 0x61, 0x6c, 0x22, 0xb7, 0x03, 0x0a,
	0x08, 0x50, 0x72, 0x6f, 0x63, 0x53, 0x74, 0x61, 0x74, 0x12, 0x18, 0x0a, 0x07, 0x74, 0x68, 0x72,
	0x65, 0x61, 0x64, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x07, 0x74, 0x68, 0x72, 0x65,
	0x61, 0x64, 0x73, 0x12, 0x10, 0x0a, 0x03, 0x66, 0x64, 0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05,
	0x52, 0x03, 0x66, 0x64, 0x73, 0x12, 0x1f, 0x0a, 0x0b, 0x63, 0x70, 0x75, 0x5f, 0x70, 0x65, 0x72,
	0x63, 0x65, 0x6e, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x02, 0x52, 0x0a, 0x63, 0x70, 0x75, 0x50,
	0x65, 0x72, 0x63, 0x65, 0x6e, 0x74, 0x12, 0x22, 0x0a, 0x03, 0x6d, 0x65, 0x6d, 0x18, 0x04, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x10, 0x2e, 0x50, 0x72, 0x6f, 0x63, 0x53, 0x74, 0x61, 0x74, 0x2e, 0x4d,
	0x65, 0x6d, 0x6f, 0x72, 0x79, 0x52, 0x03, 0x6d, 0x65, 0x6d, 0x12, 0x1c, 0x0a, 0x02, 0x69, 0x6f,
	0x18, 0x05, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0c, 0x2e, 0x50, 0x72, 0x6f, 0x63, 0x53, 0x74, 0x61,
	0x74, 0x2e, 0x49, 0x4f, 0x52, 0x02, 0x69, 0x6f, 0x1a, 0x94, 0x01, 0x0a, 0x06, 0x4d, 0x65, 0x6d,
	0x6f, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x72, 0x73, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04,
	0x52, 0x03, 0x72, 0x73, 0x73, 0x12, 0x10, 0x0a, 0x03, 0x76, 0x6d, 0x73, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x04, 0x52, 0x03, 0x76, 0x6d, 0x73, 0x12, 0x10, 0x0a, 0x03, 0x68, 0x77, 0x6d, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x04, 0x52, 0x03, 0x68, 0x77, 0x6d, 0x12, 0x12, 0x0a, 0x04, 0x64, 0x61, 0x74,
	0x61, 0x18, 0x04, 0x20, 0x01, 0x28, 0x04, 0x52, 0x04, 0x64, 0x61, 0x74, 0x61, 0x12, 0x14, 0x0a,
	0x05, 0x73, 0x74, 0x61, 0x63, 0x6b, 0x18, 0x05, 0x20, 0x01, 0x28, 0x04, 0x52, 0x05, 0x73, 0x74,
	0x61, 0x63, 0x6b, 0x12, 0x16, 0x0a, 0x06, 0x6c, 0x6f, 0x63, 0x6b, 0x65, 0x64, 0x18, 0x06, 0x20,
	0x01, 0x28, 0x04, 0x52, 0x06, 0x6c, 0x6f, 0x63, 0x6b, 0x65, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x73,
	0x77, 0x61, 0x70, 0x18, 0x07, 0x20, 0x01, 0x28, 0x04, 0x52, 0x04, 0x73, 0x77, 0x61, 0x70, 0x1a,
	0x84, 0x01, 0x0a, 0x02, 0x49, 0x4f, 0x12, 0x1d, 0x0a, 0x0a, 0x72, 0x65, 0x61, 0x64, 0x5f, 0x63,
	0x6f, 0x75, 0x6e, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x09, 0x72, 0x65, 0x61, 0x64,
	0x43, 0x6f, 0x75, 0x6e, 0x74, 0x12, 0x1f, 0x0a, 0x0b, 0x77, 0x72, 0x69, 0x74, 0x65, 0x5f, 0x63,
	0x6f, 0x75, 0x6e, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x04, 0x52, 0x0a, 0x77, 0x72, 0x69, 0x74,
	0x65, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x12, 0x1d, 0x0a, 0x0a, 0x72, 0x65, 0x61, 0x64, 0x5f, 0x62,
	0x79, 0x74, 0x65, 0x73, 0x18, 0x03, 0x20, 0x01, 0x28, 0x04, 0x52, 0x09, 0x72, 0x65, 0x61, 0x64,
	0x42, 0x79, 0x74, 0x65, 0x73, 0x12, 0x1f, 0x0a, 0x0b, 0x77, 0x72, 0x69, 0x74, 0x65, 0x5f, 0x62,
	0x79, 0x74, 0x65, 0x73, 0x18, 0x04, 0x20, 0x01, 0x28, 0x04, 0x52, 0x0a, 0x77, 0x72, 0x69, 0x74,
	0x65, 0x42, 0x79, 0x74, 0x65, 0x73, 0x22, 0x41, 0x0a, 0x07, 0x53, 0x79, 0x73, 0x53, 0x74, 0x61,
	0x74, 0x12, 0x17, 0x0a, 0x02, 0x6f, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x07, 0x2e,
	0x4f, 0x53, 0x53, 0x74, 0x61, 0x74, 0x52, 0x02, 0x6f, 0x73, 0x12, 0x1d, 0x0a, 0x04, 0x70, 0x72,
	0x6f, 0x63, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x09, 0x2e, 0x50, 0x72, 0x6f, 0x63, 0x53,
	0x74, 0x61, 0x74, 0x52, 0x04, 0x70, 0x72, 0x6f, 0x63, 0x42, 0x0f, 0x5a, 0x0d, 0x63, 0x73, 0x70,
	0x61, 0x67, 0x65, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x33,
}

var (
	file_pkg_pb_sysinfo_proto_rawDescOnce sync.Once
	file_pkg_pb_sysinfo_proto_rawDescData = file_pkg_pb_sysinfo_proto_rawDesc
)

func file_pkg_pb_sysinfo_proto_rawDescGZIP() []byte {
	file_pkg_pb_sysinfo_proto_rawDescOnce.Do(func() {
		file_pkg_pb_sysinfo_proto_rawDescData = protoimpl.X.CompressGZIP(file_pkg_pb_sysinfo_proto_rawDescData)
	})
	return file_pkg_pb_sysinfo_proto_rawDescData
}

var file_pkg_pb_sysinfo_proto_msgTypes = make([]protoimpl.MessageInfo, 7)
var file_pkg_pb_sysinfo_proto_goTypes = []any{
	(*OSStat)(nil),          // 0: OSStat
	(*ProcStat)(nil),        // 1: ProcStat
	(*SysStat)(nil),         // 2: SysStat
	(*OSStat_Memory)(nil),   // 3: OSStat.Memory
	(*OSStat_CPU)(nil),      // 4: OSStat.CPU
	(*ProcStat_Memory)(nil), // 5: ProcStat.Memory
	(*ProcStat_IO)(nil),     // 6: ProcStat.IO
}
var file_pkg_pb_sysinfo_proto_depIdxs = []int32{
	3, // 0: OSStat.mem:type_name -> OSStat.Memory
	4, // 1: OSStat.cpu:type_name -> OSStat.CPU
	5, // 2: ProcStat.mem:type_name -> ProcStat.Memory
	6, // 3: ProcStat.io:type_name -> ProcStat.IO
	0, // 4: SysStat.os:type_name -> OSStat
	1, // 5: SysStat.proc:type_name -> ProcStat
	6, // [6:6] is the sub-list for method output_type
	6, // [6:6] is the sub-list for method input_type
	6, // [6:6] is the sub-list for extension type_name
	6, // [6:6] is the sub-list for extension extendee
	0, // [0:6] is the sub-list for field type_name
}

func init() { file_pkg_pb_sysinfo_proto_init() }
func file_pkg_pb_sysinfo_proto_init() {
	if File_pkg_pb_sysinfo_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_pkg_pb_sysinfo_proto_msgTypes[0].Exporter = func(v any, i int) any {
			switch v := v.(*OSStat); i {
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
		file_pkg_pb_sysinfo_proto_msgTypes[1].Exporter = func(v any, i int) any {
			switch v := v.(*ProcStat); i {
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
		file_pkg_pb_sysinfo_proto_msgTypes[2].Exporter = func(v any, i int) any {
			switch v := v.(*SysStat); i {
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
		file_pkg_pb_sysinfo_proto_msgTypes[3].Exporter = func(v any, i int) any {
			switch v := v.(*OSStat_Memory); i {
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
		file_pkg_pb_sysinfo_proto_msgTypes[4].Exporter = func(v any, i int) any {
			switch v := v.(*OSStat_CPU); i {
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
		file_pkg_pb_sysinfo_proto_msgTypes[5].Exporter = func(v any, i int) any {
			switch v := v.(*ProcStat_Memory); i {
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
		file_pkg_pb_sysinfo_proto_msgTypes[6].Exporter = func(v any, i int) any {
			switch v := v.(*ProcStat_IO); i {
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
			RawDescriptor: file_pkg_pb_sysinfo_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   7,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_pkg_pb_sysinfo_proto_goTypes,
		DependencyIndexes: file_pkg_pb_sysinfo_proto_depIdxs,
		MessageInfos:      file_pkg_pb_sysinfo_proto_msgTypes,
	}.Build()
	File_pkg_pb_sysinfo_proto = out.File
	file_pkg_pb_sysinfo_proto_rawDesc = nil
	file_pkg_pb_sysinfo_proto_goTypes = nil
	file_pkg_pb_sysinfo_proto_depIdxs = nil
}
