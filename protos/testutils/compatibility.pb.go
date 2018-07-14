// Code generated by protoc-gen-go. DO NOT EDIT.
// source: testutils/compatibility.proto

/*
Package testutils is a generated protocol buffer package.

It is generated from these files:
	testutils/compatibility.proto

It has these top-level messages:
	String
	Bytes
*/
package testutils

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

type String struct {
	Str string `protobuf:"bytes,1,opt,name=str" json:"str,omitempty"`
}

func (m *String) Reset()                    { *m = String{} }
func (m *String) String() string            { return proto.CompactTextString(m) }
func (*String) ProtoMessage()               {}
func (*String) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *String) GetStr() string {
	if m != nil {
		return m.Str
	}
	return ""
}

type Bytes struct {
	B []byte `protobuf:"bytes,1,opt,name=b,proto3" json:"b,omitempty"`
}

func (m *Bytes) Reset()                    { *m = Bytes{} }
func (m *Bytes) String() string            { return proto.CompactTextString(m) }
func (*Bytes) ProtoMessage()               {}
func (*Bytes) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *Bytes) GetB() []byte {
	if m != nil {
		return m.B
	}
	return nil
}

func init() {
	proto.RegisterType((*String)(nil), "String")
	proto.RegisterType((*Bytes)(nil), "Bytes")
}

func init() { proto.RegisterFile("testutils/compatibility.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 162 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x92, 0x2d, 0x49, 0x2d, 0x2e,
	0x29, 0x2d, 0xc9, 0xcc, 0x29, 0xd6, 0x4f, 0xce, 0xcf, 0x2d, 0x48, 0x2c, 0xc9, 0x4c, 0xca, 0xcc,
	0xc9, 0x2c, 0xa9, 0xd4, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x57, 0x92, 0xe2, 0x62, 0x0b, 0x2e, 0x29,
	0xca, 0xcc, 0x4b, 0x17, 0x12, 0xe0, 0x62, 0x2e, 0x2e, 0x29, 0x92, 0x60, 0x54, 0x60, 0xd4, 0xe0,
	0x0c, 0x02, 0x31, 0x95, 0x44, 0xb9, 0x58, 0x9d, 0x2a, 0x4b, 0x52, 0x8b, 0x85, 0x78, 0xb8, 0x18,
	0x93, 0xc0, 0x12, 0x3c, 0x41, 0x8c, 0x49, 0x4e, 0x91, 0x5c, 0xea, 0xf9, 0x45, 0xe9, 0x7a, 0x19,
	0x95, 0x05, 0xa9, 0x45, 0x39, 0xa9, 0x29, 0xe9, 0xa9, 0x45, 0x7a, 0x69, 0x89, 0x49, 0x45, 0x99,
	0xc9, 0x10, 0x23, 0x8b, 0xf5, 0xe0, 0x36, 0x46, 0xe9, 0xa5, 0x67, 0x96, 0x64, 0x94, 0x26, 0xe9,
	0x25, 0xe7, 0xe7, 0xea, 0x23, 0xa9, 0xd7, 0x87, 0xa8, 0xd7, 0x87, 0xa8, 0xd7, 0x87, 0xab, 0x4f,
	0x62, 0x03, 0x8b, 0x18, 0x03, 0x02, 0x00, 0x00, 0xff, 0xff, 0xc5, 0x2b, 0x90, 0x48, 0xb5, 0x00,
	0x00, 0x00,
}