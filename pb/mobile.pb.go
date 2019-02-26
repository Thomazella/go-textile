// Code generated by protoc-gen-go. DO NOT EDIT.
// source: mobile.proto

package pb

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type MobileEvent_Type int32

const (
	MobileEvent_NODE_START     MobileEvent_Type = 0
	MobileEvent_NODE_ONLINE    MobileEvent_Type = 1
	MobileEvent_NODE_STOP      MobileEvent_Type = 2
	MobileEvent_WALLET_UPDATE  MobileEvent_Type = 10
	MobileEvent_THREAD_UPDATE  MobileEvent_Type = 11
	MobileEvent_NOTIFICATION   MobileEvent_Type = 12
	MobileEvent_QUERY_RESPONSE MobileEvent_Type = 20
)

var MobileEvent_Type_name = map[int32]string{
	0:  "NODE_START",
	1:  "NODE_ONLINE",
	2:  "NODE_STOP",
	10: "WALLET_UPDATE",
	11: "THREAD_UPDATE",
	12: "NOTIFICATION",
	20: "QUERY_RESPONSE",
}
var MobileEvent_Type_value = map[string]int32{
	"NODE_START":     0,
	"NODE_ONLINE":    1,
	"NODE_STOP":      2,
	"WALLET_UPDATE":  10,
	"THREAD_UPDATE":  11,
	"NOTIFICATION":   12,
	"QUERY_RESPONSE": 20,
}

func (x MobileEvent_Type) String() string {
	return proto.EnumName(MobileEvent_Type_name, int32(x))
}
func (MobileEvent_Type) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_mobile_b3486e87e4932f04, []int{3, 0}
}

type MobileQueryEvent_Type int32

const (
	MobileQueryEvent_DATA  MobileQueryEvent_Type = 0
	MobileQueryEvent_DONE  MobileQueryEvent_Type = 1
	MobileQueryEvent_ERROR MobileQueryEvent_Type = 2
)

var MobileQueryEvent_Type_name = map[int32]string{
	0: "DATA",
	1: "DONE",
	2: "ERROR",
}
var MobileQueryEvent_Type_value = map[string]int32{
	"DATA":  0,
	"DONE":  1,
	"ERROR": 2,
}

func (x MobileQueryEvent_Type) String() string {
	return proto.EnumName(MobileQueryEvent_Type_name, int32(x))
}
func (MobileQueryEvent_Type) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_mobile_b3486e87e4932f04, []int{4, 0}
}

type MobileWalletAccount struct {
	Seed                 string   `protobuf:"bytes,1,opt,name=seed,proto3" json:"seed,omitempty"`
	Address              string   `protobuf:"bytes,2,opt,name=address,proto3" json:"address,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *MobileWalletAccount) Reset()         { *m = MobileWalletAccount{} }
func (m *MobileWalletAccount) String() string { return proto.CompactTextString(m) }
func (*MobileWalletAccount) ProtoMessage()    {}
func (*MobileWalletAccount) Descriptor() ([]byte, []int) {
	return fileDescriptor_mobile_b3486e87e4932f04, []int{0}
}
func (m *MobileWalletAccount) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_MobileWalletAccount.Unmarshal(m, b)
}
func (m *MobileWalletAccount) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_MobileWalletAccount.Marshal(b, m, deterministic)
}
func (dst *MobileWalletAccount) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MobileWalletAccount.Merge(dst, src)
}
func (m *MobileWalletAccount) XXX_Size() int {
	return xxx_messageInfo_MobileWalletAccount.Size(m)
}
func (m *MobileWalletAccount) XXX_DiscardUnknown() {
	xxx_messageInfo_MobileWalletAccount.DiscardUnknown(m)
}

var xxx_messageInfo_MobileWalletAccount proto.InternalMessageInfo

func (m *MobileWalletAccount) GetSeed() string {
	if m != nil {
		return m.Seed
	}
	return ""
}

func (m *MobileWalletAccount) GetAddress() string {
	if m != nil {
		return m.Address
	}
	return ""
}

type MobilePreparedFiles struct {
	Dir                  *Directory        `protobuf:"bytes,1,opt,name=dir,proto3" json:"dir,omitempty"`
	Pin                  map[string]string `protobuf:"bytes,2,rep,name=pin,proto3" json:"pin,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	XXX_NoUnkeyedLiteral struct{}          `json:"-"`
	XXX_unrecognized     []byte            `json:"-"`
	XXX_sizecache        int32             `json:"-"`
}

func (m *MobilePreparedFiles) Reset()         { *m = MobilePreparedFiles{} }
func (m *MobilePreparedFiles) String() string { return proto.CompactTextString(m) }
func (*MobilePreparedFiles) ProtoMessage()    {}
func (*MobilePreparedFiles) Descriptor() ([]byte, []int) {
	return fileDescriptor_mobile_b3486e87e4932f04, []int{1}
}
func (m *MobilePreparedFiles) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_MobilePreparedFiles.Unmarshal(m, b)
}
func (m *MobilePreparedFiles) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_MobilePreparedFiles.Marshal(b, m, deterministic)
}
func (dst *MobilePreparedFiles) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MobilePreparedFiles.Merge(dst, src)
}
func (m *MobilePreparedFiles) XXX_Size() int {
	return xxx_messageInfo_MobilePreparedFiles.Size(m)
}
func (m *MobilePreparedFiles) XXX_DiscardUnknown() {
	xxx_messageInfo_MobilePreparedFiles.DiscardUnknown(m)
}

var xxx_messageInfo_MobilePreparedFiles proto.InternalMessageInfo

func (m *MobilePreparedFiles) GetDir() *Directory {
	if m != nil {
		return m.Dir
	}
	return nil
}

func (m *MobilePreparedFiles) GetPin() map[string]string {
	if m != nil {
		return m.Pin
	}
	return nil
}

type MobileFileData struct {
	Url                  string   `protobuf:"bytes,1,opt,name=url,proto3" json:"url,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *MobileFileData) Reset()         { *m = MobileFileData{} }
func (m *MobileFileData) String() string { return proto.CompactTextString(m) }
func (*MobileFileData) ProtoMessage()    {}
func (*MobileFileData) Descriptor() ([]byte, []int) {
	return fileDescriptor_mobile_b3486e87e4932f04, []int{2}
}
func (m *MobileFileData) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_MobileFileData.Unmarshal(m, b)
}
func (m *MobileFileData) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_MobileFileData.Marshal(b, m, deterministic)
}
func (dst *MobileFileData) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MobileFileData.Merge(dst, src)
}
func (m *MobileFileData) XXX_Size() int {
	return xxx_messageInfo_MobileFileData.Size(m)
}
func (m *MobileFileData) XXX_DiscardUnknown() {
	xxx_messageInfo_MobileFileData.DiscardUnknown(m)
}

var xxx_messageInfo_MobileFileData proto.InternalMessageInfo

func (m *MobileFileData) GetUrl() string {
	if m != nil {
		return m.Url
	}
	return ""
}

type MobileEvent struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *MobileEvent) Reset()         { *m = MobileEvent{} }
func (m *MobileEvent) String() string { return proto.CompactTextString(m) }
func (*MobileEvent) ProtoMessage()    {}
func (*MobileEvent) Descriptor() ([]byte, []int) {
	return fileDescriptor_mobile_b3486e87e4932f04, []int{3}
}
func (m *MobileEvent) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_MobileEvent.Unmarshal(m, b)
}
func (m *MobileEvent) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_MobileEvent.Marshal(b, m, deterministic)
}
func (dst *MobileEvent) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MobileEvent.Merge(dst, src)
}
func (m *MobileEvent) XXX_Size() int {
	return xxx_messageInfo_MobileEvent.Size(m)
}
func (m *MobileEvent) XXX_DiscardUnknown() {
	xxx_messageInfo_MobileEvent.DiscardUnknown(m)
}

var xxx_messageInfo_MobileEvent proto.InternalMessageInfo

type MobileQueryEvent struct {
	Id                   string                `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Type                 MobileQueryEvent_Type `protobuf:"varint,2,opt,name=type,proto3,enum=MobileQueryEvent_Type" json:"type,omitempty"`
	Data                 *QueryResult          `protobuf:"bytes,3,opt,name=data,proto3" json:"data,omitempty"`
	Error                *Error                `protobuf:"bytes,4,opt,name=error,proto3" json:"error,omitempty"`
	XXX_NoUnkeyedLiteral struct{}              `json:"-"`
	XXX_unrecognized     []byte                `json:"-"`
	XXX_sizecache        int32                 `json:"-"`
}

func (m *MobileQueryEvent) Reset()         { *m = MobileQueryEvent{} }
func (m *MobileQueryEvent) String() string { return proto.CompactTextString(m) }
func (*MobileQueryEvent) ProtoMessage()    {}
func (*MobileQueryEvent) Descriptor() ([]byte, []int) {
	return fileDescriptor_mobile_b3486e87e4932f04, []int{4}
}
func (m *MobileQueryEvent) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_MobileQueryEvent.Unmarshal(m, b)
}
func (m *MobileQueryEvent) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_MobileQueryEvent.Marshal(b, m, deterministic)
}
func (dst *MobileQueryEvent) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MobileQueryEvent.Merge(dst, src)
}
func (m *MobileQueryEvent) XXX_Size() int {
	return xxx_messageInfo_MobileQueryEvent.Size(m)
}
func (m *MobileQueryEvent) XXX_DiscardUnknown() {
	xxx_messageInfo_MobileQueryEvent.DiscardUnknown(m)
}

var xxx_messageInfo_MobileQueryEvent proto.InternalMessageInfo

func (m *MobileQueryEvent) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *MobileQueryEvent) GetType() MobileQueryEvent_Type {
	if m != nil {
		return m.Type
	}
	return MobileQueryEvent_DATA
}

func (m *MobileQueryEvent) GetData() *QueryResult {
	if m != nil {
		return m.Data
	}
	return nil
}

func (m *MobileQueryEvent) GetError() *Error {
	if m != nil {
		return m.Error
	}
	return nil
}

func init() {
	proto.RegisterType((*MobileWalletAccount)(nil), "MobileWalletAccount")
	proto.RegisterType((*MobilePreparedFiles)(nil), "MobilePreparedFiles")
	proto.RegisterMapType((map[string]string)(nil), "MobilePreparedFiles.PinEntry")
	proto.RegisterType((*MobileFileData)(nil), "MobileFileData")
	proto.RegisterType((*MobileEvent)(nil), "MobileEvent")
	proto.RegisterType((*MobileQueryEvent)(nil), "MobileQueryEvent")
	proto.RegisterEnum("MobileEvent_Type", MobileEvent_Type_name, MobileEvent_Type_value)
	proto.RegisterEnum("MobileQueryEvent_Type", MobileQueryEvent_Type_name, MobileQueryEvent_Type_value)
}

func init() { proto.RegisterFile("mobile.proto", fileDescriptor_mobile_b3486e87e4932f04) }

var fileDescriptor_mobile_b3486e87e4932f04 = []byte{
	// 460 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x6c, 0x52, 0x41, 0x6b, 0xd4, 0x40,
	0x18, 0x6d, 0xb2, 0xd9, 0xda, 0xfd, 0xb2, 0x1b, 0xc7, 0xb1, 0x48, 0x28, 0x15, 0x96, 0x80, 0x50,
	0x3c, 0x44, 0x58, 0x41, 0xc4, 0x5b, 0x6c, 0xa6, 0xb8, 0xb0, 0x26, 0xe9, 0x6c, 0x4a, 0xd1, 0xcb,
	0x92, 0xdd, 0x0c, 0x32, 0x98, 0x26, 0x71, 0x32, 0x59, 0xc9, 0xd5, 0x9f, 0xe0, 0xd9, 0xbf, 0xe1,
	0xff, 0x93, 0x99, 0x24, 0x17, 0xe9, 0xed, 0x7b, 0x2f, 0xef, 0x3d, 0xde, 0x7c, 0x5f, 0x60, 0xfe,
	0x50, 0xed, 0x79, 0xc1, 0xfc, 0x5a, 0x54, 0xb2, 0xba, 0x80, 0x23, 0x67, 0x3f, 0x87, 0xd9, 0xfe,
	0xd1, 0x32, 0xd1, 0x0d, 0x60, 0xf1, 0xc0, 0x9a, 0x26, 0xfb, 0x36, 0xe8, 0xbc, 0x6b, 0x78, 0xfe,
	0x59, 0xfb, 0xee, 0xb3, 0xa2, 0x60, 0x32, 0x38, 0x1c, 0xaa, 0xb6, 0x94, 0x18, 0x83, 0xd5, 0x30,
	0x96, 0xbb, 0xc6, 0xd2, 0xb8, 0x9a, 0x51, 0x3d, 0x63, 0x17, 0x9e, 0x64, 0x79, 0x2e, 0x58, 0xd3,
	0xb8, 0xa6, 0xa6, 0x47, 0xe8, 0xfd, 0x31, 0xc6, 0x94, 0x44, 0xb0, 0x3a, 0x13, 0x2c, 0xbf, 0xe1,
	0x05, 0x6b, 0xf0, 0x25, 0x4c, 0x72, 0x2e, 0x74, 0x88, 0xbd, 0x02, 0x3f, 0xe4, 0x82, 0x1d, 0x64,
	0x25, 0x3a, 0xaa, 0x68, 0xfc, 0x06, 0x26, 0x35, 0x2f, 0x5d, 0x73, 0x39, 0xb9, 0xb2, 0x57, 0x2f,
	0xfd, 0x47, 0x02, 0xfc, 0x84, 0x97, 0xa4, 0x94, 0xca, 0x50, 0xf3, 0xf2, 0xe2, 0x1d, 0x9c, 0x8d,
	0x04, 0x46, 0x30, 0xf9, 0xce, 0xba, 0xa1, 0x9f, 0x1a, 0xf1, 0x39, 0x4c, 0x8f, 0x59, 0xd1, 0xb2,
	0xa1, 0x5c, 0x0f, 0x3e, 0x98, 0xef, 0x0d, 0xcf, 0x03, 0xa7, 0x0f, 0x57, 0xa1, 0x61, 0x26, 0x33,
	0xe5, 0x6e, 0x45, 0x31, 0xba, 0x5b, 0x51, 0x78, 0xbf, 0x0d, 0xb0, 0x7b, 0x11, 0x39, 0xb2, 0x52,
	0x7a, 0xbf, 0x0c, 0xb0, 0xd2, 0xae, 0x66, 0xd8, 0x01, 0x88, 0xe2, 0x90, 0xec, 0xb6, 0x69, 0x40,
	0x53, 0x74, 0x82, 0x9f, 0x82, 0xad, 0x71, 0x1c, 0x6d, 0xd6, 0x11, 0x41, 0x06, 0x5e, 0xc0, 0x6c,
	0x10, 0xc4, 0x09, 0x32, 0xf1, 0x33, 0x58, 0xdc, 0x07, 0x9b, 0x0d, 0x49, 0x77, 0x77, 0x49, 0x18,
	0xa4, 0x04, 0x81, 0xa2, 0xd2, 0x4f, 0x94, 0x04, 0xe1, 0x48, 0xd9, 0x18, 0xc1, 0x3c, 0x8a, 0xd3,
	0xf5, 0xcd, 0xfa, 0x3a, 0x48, 0xd7, 0x71, 0x84, 0xe6, 0x18, 0x83, 0x73, 0x7b, 0x47, 0xe8, 0x97,
	0x1d, 0x25, 0xdb, 0x24, 0x8e, 0xb6, 0x04, 0x9d, 0x7b, 0x7f, 0x0d, 0x40, 0x7d, 0xa9, 0x5b, 0x75,
	0x41, 0xdd, 0x0c, 0x3b, 0x60, 0xf2, 0xf1, 0x30, 0x26, 0xcf, 0xf1, 0x6b, 0xb0, 0x64, 0x57, 0xf7,
	0xcf, 0x76, 0x56, 0x2f, 0xfc, 0xff, 0x0d, 0xbe, 0x7a, 0x06, 0xd5, 0x1a, 0xbc, 0x04, 0x2b, 0xcf,
	0x64, 0xe6, 0x4e, 0xf4, 0x45, 0xe6, 0xbe, 0x56, 0x51, 0xd6, 0xb4, 0x85, 0xa4, 0xfa, 0x0b, 0xbe,
	0x84, 0x29, 0x13, 0xa2, 0x12, 0xae, 0xa5, 0x25, 0xa7, 0x3e, 0x51, 0x88, 0xf6, 0xa4, 0xf7, 0x6a,
	0x58, 0xca, 0x19, 0x58, 0x61, 0x90, 0x06, 0xe8, 0x44, 0x4f, 0xb1, 0xde, 0xc3, 0x0c, 0xa6, 0x84,
	0xd2, 0x98, 0x22, 0xf3, 0xa3, 0xf5, 0xd5, 0xac, 0xf7, 0xfb, 0x53, 0xfd, 0x87, 0xbd, 0xfd, 0x17,
	0x00, 0x00, 0xff, 0xff, 0xf0, 0x78, 0x3a, 0xbe, 0x99, 0x02, 0x00, 0x00,
}
