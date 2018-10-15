// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: teststruct.proto

package benchmark

import proto "github.com/gogo/protobuf/proto"
import fmt "fmt"
import math "math"

import io "io"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion2 // please upgrade the proto package

type TestStruct2 struct {
	Id      uint64 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Key     string `protobuf:"bytes,2,opt,name=key,proto3" json:"key,omitempty"`
	Msg     string `protobuf:"bytes,3,opt,name=msg,proto3" json:"msg,omitempty"`
	Counter int32  `protobuf:"varint,4,opt,name=counter,proto3" json:"counter,omitempty"`
	Flag    bool   `protobuf:"varint,5,opt,name=flag,proto3" json:"flag,omitempty"`
}

func (m *TestStruct2) Reset()         { *m = TestStruct2{} }
func (m *TestStruct2) String() string { return proto.CompactTextString(m) }
func (*TestStruct2) ProtoMessage()    {}
func (*TestStruct2) Descriptor() ([]byte, []int) {
	return fileDescriptor_teststruct_be2f434f12a989fc, []int{0}
}
func (m *TestStruct2) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *TestStruct2) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_TestStruct2.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalTo(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (dst *TestStruct2) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TestStruct2.Merge(dst, src)
}
func (m *TestStruct2) XXX_Size() int {
	return m.Size()
}
func (m *TestStruct2) XXX_DiscardUnknown() {
	xxx_messageInfo_TestStruct2.DiscardUnknown(m)
}

var xxx_messageInfo_TestStruct2 proto.InternalMessageInfo

func (m *TestStruct2) GetId() uint64 {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *TestStruct2) GetKey() string {
	if m != nil {
		return m.Key
	}
	return ""
}

func (m *TestStruct2) GetMsg() string {
	if m != nil {
		return m.Msg
	}
	return ""
}

func (m *TestStruct2) GetCounter() int32 {
	if m != nil {
		return m.Counter
	}
	return 0
}

func (m *TestStruct2) GetFlag() bool {
	if m != nil {
		return m.Flag
	}
	return false
}

func init() {
	proto.RegisterType((*TestStruct2)(nil), "benchmark.TestStruct2")
}
func (m *TestStruct2) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *TestStruct2) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if m.Id != 0 {
		dAtA[i] = 0x8
		i++
		i = encodeVarintTeststruct(dAtA, i, uint64(m.Id))
	}
	if len(m.Key) > 0 {
		dAtA[i] = 0x12
		i++
		i = encodeVarintTeststruct(dAtA, i, uint64(len(m.Key)))
		i += copy(dAtA[i:], m.Key)
	}
	if len(m.Msg) > 0 {
		dAtA[i] = 0x1a
		i++
		i = encodeVarintTeststruct(dAtA, i, uint64(len(m.Msg)))
		i += copy(dAtA[i:], m.Msg)
	}
	if m.Counter != 0 {
		dAtA[i] = 0x20
		i++
		i = encodeVarintTeststruct(dAtA, i, uint64(m.Counter))
	}
	if m.Flag {
		dAtA[i] = 0x28
		i++
		if m.Flag {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i++
	}
	return i, nil
}

func encodeVarintTeststruct(dAtA []byte, offset int, v uint64) int {
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return offset + 1
}
func (m *TestStruct2) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Id != 0 {
		n += 1 + sovTeststruct(uint64(m.Id))
	}
	l = len(m.Key)
	if l > 0 {
		n += 1 + l + sovTeststruct(uint64(l))
	}
	l = len(m.Msg)
	if l > 0 {
		n += 1 + l + sovTeststruct(uint64(l))
	}
	if m.Counter != 0 {
		n += 1 + sovTeststruct(uint64(m.Counter))
	}
	if m.Flag {
		n += 2
	}
	return n
}

func sovTeststruct(x uint64) (n int) {
	for {
		n++
		x >>= 7
		if x == 0 {
			break
		}
	}
	return n
}
func sozTeststruct(x uint64) (n int) {
	return sovTeststruct(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *TestStruct2) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTeststruct
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: TestStruct2: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: TestStruct2: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Id", wireType)
			}
			m.Id = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTeststruct
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Id |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Key", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTeststruct
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthTeststruct
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Key = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Msg", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTeststruct
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthTeststruct
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Msg = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 4:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Counter", wireType)
			}
			m.Counter = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTeststruct
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Counter |= (int32(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 5:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Flag", wireType)
			}
			var v int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTeststruct
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				v |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			m.Flag = bool(v != 0)
		default:
			iNdEx = preIndex
			skippy, err := skipTeststruct(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthTeststruct
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipTeststruct(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowTeststruct
			}
			if iNdEx >= l {
				return 0, io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		wireType := int(wire & 0x7)
		switch wireType {
		case 0:
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowTeststruct
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
			return iNdEx, nil
		case 1:
			iNdEx += 8
			return iNdEx, nil
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowTeststruct
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				length |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			iNdEx += length
			if length < 0 {
				return 0, ErrInvalidLengthTeststruct
			}
			return iNdEx, nil
		case 3:
			for {
				var innerWire uint64
				var start int = iNdEx
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return 0, ErrIntOverflowTeststruct
					}
					if iNdEx >= l {
						return 0, io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					innerWire |= (uint64(b) & 0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				innerWireType := int(innerWire & 0x7)
				if innerWireType == 4 {
					break
				}
				next, err := skipTeststruct(dAtA[start:])
				if err != nil {
					return 0, err
				}
				iNdEx = start + next
			}
			return iNdEx, nil
		case 4:
			return iNdEx, nil
		case 5:
			iNdEx += 4
			return iNdEx, nil
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
	}
	panic("unreachable")
}

var (
	ErrInvalidLengthTeststruct = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowTeststruct   = fmt.Errorf("proto: integer overflow")
)

func init() { proto.RegisterFile("teststruct.proto", fileDescriptor_teststruct_be2f434f12a989fc) }

var fileDescriptor_teststruct_be2f434f12a989fc = []byte{
	// 178 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x12, 0x28, 0x49, 0x2d, 0x2e,
	0x29, 0x2e, 0x29, 0x2a, 0x4d, 0x2e, 0xd1, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0xe2, 0x4c, 0x4a,
	0xcd, 0x4b, 0xce, 0xc8, 0x4d, 0x2c, 0xca, 0x56, 0xca, 0xe7, 0xe2, 0x0e, 0x49, 0x2d, 0x2e, 0x09,
	0x06, 0x4b, 0x1b, 0x09, 0xf1, 0x71, 0x31, 0x65, 0xa6, 0x48, 0x30, 0x2a, 0x30, 0x6a, 0xb0, 0x04,
	0x31, 0x65, 0xa6, 0x08, 0x09, 0x70, 0x31, 0x67, 0xa7, 0x56, 0x4a, 0x30, 0x29, 0x30, 0x6a, 0x70,
	0x06, 0x81, 0x98, 0x20, 0x91, 0xdc, 0xe2, 0x74, 0x09, 0x66, 0x88, 0x48, 0x6e, 0x71, 0xba, 0x90,
	0x04, 0x17, 0x7b, 0x72, 0x7e, 0x69, 0x5e, 0x49, 0x6a, 0x91, 0x04, 0x8b, 0x02, 0xa3, 0x06, 0x6b,
	0x10, 0x8c, 0x2b, 0x24, 0xc4, 0xc5, 0x92, 0x96, 0x93, 0x98, 0x2e, 0xc1, 0xaa, 0xc0, 0xa8, 0xc1,
	0x11, 0x04, 0x66, 0x3b, 0x49, 0x9c, 0x78, 0x24, 0xc7, 0x78, 0xe1, 0x91, 0x1c, 0xe3, 0x83, 0x47,
	0x72, 0x8c, 0x13, 0x1e, 0xcb, 0x31, 0x5c, 0x78, 0x2c, 0xc7, 0x70, 0xe3, 0xb1, 0x1c, 0x43, 0x12,
	0x1b, 0xd8, 0x71, 0xc6, 0x80, 0x00, 0x00, 0x00, 0xff, 0xff, 0x2e, 0x15, 0x29, 0xf2, 0xb0, 0x00,
	0x00, 0x00,
}