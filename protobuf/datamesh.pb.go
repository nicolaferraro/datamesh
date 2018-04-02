// Code generated by protoc-gen-go. DO NOT EDIT.
// source: protobuf/datamesh.proto

/*
Package protobuf is a generated protocol buffer package.

It is generated from these files:
	protobuf/datamesh.proto

It has these top-level messages:
	Transaction
	Operation
	ReadOperation
	UpsertOperation
	DeleteOperation
	GenerateEventOperation
	ApplicationFailure
	Event
	Path
	Data
*/
package protobuf

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import google_protobuf "github.com/golang/protobuf/ptypes/struct"
import google_protobuf1 "github.com/golang/protobuf/ptypes/empty"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

// A transaction is a sequence of operation triggered by a event.
type Transaction struct {
	// The trigger
	Event *Event `protobuf:"bytes,1,opt,name=event" json:"event,omitempty"`
	// The operations
	Operations []*Operation `protobuf:"bytes,2,rep,name=operations" json:"operations,omitempty"`
}

func (m *Transaction) Reset()                    { *m = Transaction{} }
func (m *Transaction) String() string            { return proto.CompactTextString(m) }
func (*Transaction) ProtoMessage()               {}
func (*Transaction) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *Transaction) GetEvent() *Event {
	if m != nil {
		return m.Event
	}
	return nil
}

func (m *Transaction) GetOperations() []*Operation {
	if m != nil {
		return m.Operations
	}
	return nil
}

type Operation struct {
	// Types that are valid to be assigned to Kind:
	//	*Operation_Read
	//	*Operation_Upsert
	//	*Operation_Delete
	//	*Operation_Generate
	//	*Operation_Failure
	Kind isOperation_Kind `protobuf_oneof:"kind"`
}

func (m *Operation) Reset()                    { *m = Operation{} }
func (m *Operation) String() string            { return proto.CompactTextString(m) }
func (*Operation) ProtoMessage()               {}
func (*Operation) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

type isOperation_Kind interface {
	isOperation_Kind()
}

type Operation_Read struct {
	Read *ReadOperation `protobuf:"bytes,1,opt,name=read,oneof"`
}
type Operation_Upsert struct {
	Upsert *UpsertOperation `protobuf:"bytes,2,opt,name=upsert,oneof"`
}
type Operation_Delete struct {
	Delete *DeleteOperation `protobuf:"bytes,3,opt,name=delete,oneof"`
}
type Operation_Generate struct {
	Generate *GenerateEventOperation `protobuf:"bytes,4,opt,name=generate,oneof"`
}
type Operation_Failure struct {
	Failure *ApplicationFailure `protobuf:"bytes,5,opt,name=failure,oneof"`
}

func (*Operation_Read) isOperation_Kind()     {}
func (*Operation_Upsert) isOperation_Kind()   {}
func (*Operation_Delete) isOperation_Kind()   {}
func (*Operation_Generate) isOperation_Kind() {}
func (*Operation_Failure) isOperation_Kind()  {}

func (m *Operation) GetKind() isOperation_Kind {
	if m != nil {
		return m.Kind
	}
	return nil
}

func (m *Operation) GetRead() *ReadOperation {
	if x, ok := m.GetKind().(*Operation_Read); ok {
		return x.Read
	}
	return nil
}

func (m *Operation) GetUpsert() *UpsertOperation {
	if x, ok := m.GetKind().(*Operation_Upsert); ok {
		return x.Upsert
	}
	return nil
}

func (m *Operation) GetDelete() *DeleteOperation {
	if x, ok := m.GetKind().(*Operation_Delete); ok {
		return x.Delete
	}
	return nil
}

func (m *Operation) GetGenerate() *GenerateEventOperation {
	if x, ok := m.GetKind().(*Operation_Generate); ok {
		return x.Generate
	}
	return nil
}

func (m *Operation) GetFailure() *ApplicationFailure {
	if x, ok := m.GetKind().(*Operation_Failure); ok {
		return x.Failure
	}
	return nil
}

// XXX_OneofFuncs is for the internal use of the proto package.
func (*Operation) XXX_OneofFuncs() (func(msg proto.Message, b *proto.Buffer) error, func(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error), func(msg proto.Message) (n int), []interface{}) {
	return _Operation_OneofMarshaler, _Operation_OneofUnmarshaler, _Operation_OneofSizer, []interface{}{
		(*Operation_Read)(nil),
		(*Operation_Upsert)(nil),
		(*Operation_Delete)(nil),
		(*Operation_Generate)(nil),
		(*Operation_Failure)(nil),
	}
}

func _Operation_OneofMarshaler(msg proto.Message, b *proto.Buffer) error {
	m := msg.(*Operation)
	// kind
	switch x := m.Kind.(type) {
	case *Operation_Read:
		b.EncodeVarint(1<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.Read); err != nil {
			return err
		}
	case *Operation_Upsert:
		b.EncodeVarint(2<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.Upsert); err != nil {
			return err
		}
	case *Operation_Delete:
		b.EncodeVarint(3<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.Delete); err != nil {
			return err
		}
	case *Operation_Generate:
		b.EncodeVarint(4<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.Generate); err != nil {
			return err
		}
	case *Operation_Failure:
		b.EncodeVarint(5<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.Failure); err != nil {
			return err
		}
	case nil:
	default:
		return fmt.Errorf("Operation.Kind has unexpected type %T", x)
	}
	return nil
}

func _Operation_OneofUnmarshaler(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error) {
	m := msg.(*Operation)
	switch tag {
	case 1: // kind.read
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(ReadOperation)
		err := b.DecodeMessage(msg)
		m.Kind = &Operation_Read{msg}
		return true, err
	case 2: // kind.upsert
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(UpsertOperation)
		err := b.DecodeMessage(msg)
		m.Kind = &Operation_Upsert{msg}
		return true, err
	case 3: // kind.delete
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(DeleteOperation)
		err := b.DecodeMessage(msg)
		m.Kind = &Operation_Delete{msg}
		return true, err
	case 4: // kind.generate
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(GenerateEventOperation)
		err := b.DecodeMessage(msg)
		m.Kind = &Operation_Generate{msg}
		return true, err
	case 5: // kind.failure
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(ApplicationFailure)
		err := b.DecodeMessage(msg)
		m.Kind = &Operation_Failure{msg}
		return true, err
	default:
		return false, nil
	}
}

func _Operation_OneofSizer(msg proto.Message) (n int) {
	m := msg.(*Operation)
	// kind
	switch x := m.Kind.(type) {
	case *Operation_Read:
		s := proto.Size(x.Read)
		n += proto.SizeVarint(1<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(s))
		n += s
	case *Operation_Upsert:
		s := proto.Size(x.Upsert)
		n += proto.SizeVarint(2<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(s))
		n += s
	case *Operation_Delete:
		s := proto.Size(x.Delete)
		n += proto.SizeVarint(3<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(s))
		n += s
	case *Operation_Generate:
		s := proto.Size(x.Generate)
		n += proto.SizeVarint(4<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(s))
		n += s
	case *Operation_Failure:
		s := proto.Size(x.Failure)
		n += proto.SizeVarint(5<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(s))
		n += s
	case nil:
	default:
		panic(fmt.Sprintf("proto: unexpected type %T in oneof", x))
	}
	return n
}

type ReadOperation struct {
	Path *Path `protobuf:"bytes,1,opt,name=path" json:"path,omitempty"`
}

func (m *ReadOperation) Reset()                    { *m = ReadOperation{} }
func (m *ReadOperation) String() string            { return proto.CompactTextString(m) }
func (*ReadOperation) ProtoMessage()               {}
func (*ReadOperation) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *ReadOperation) GetPath() *Path {
	if m != nil {
		return m.Path
	}
	return nil
}

type UpsertOperation struct {
	Data *Data `protobuf:"bytes,1,opt,name=data" json:"data,omitempty"`
}

func (m *UpsertOperation) Reset()                    { *m = UpsertOperation{} }
func (m *UpsertOperation) String() string            { return proto.CompactTextString(m) }
func (*UpsertOperation) ProtoMessage()               {}
func (*UpsertOperation) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

func (m *UpsertOperation) GetData() *Data {
	if m != nil {
		return m.Data
	}
	return nil
}

type DeleteOperation struct {
	Path *Path `protobuf:"bytes,1,opt,name=path" json:"path,omitempty"`
}

func (m *DeleteOperation) Reset()                    { *m = DeleteOperation{} }
func (m *DeleteOperation) String() string            { return proto.CompactTextString(m) }
func (*DeleteOperation) ProtoMessage()               {}
func (*DeleteOperation) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{4} }

func (m *DeleteOperation) GetPath() *Path {
	if m != nil {
		return m.Path
	}
	return nil
}

type GenerateEventOperation struct {
	Event *Event `protobuf:"bytes,1,opt,name=event" json:"event,omitempty"`
}

func (m *GenerateEventOperation) Reset()                    { *m = GenerateEventOperation{} }
func (m *GenerateEventOperation) String() string            { return proto.CompactTextString(m) }
func (*GenerateEventOperation) ProtoMessage()               {}
func (*GenerateEventOperation) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{5} }

func (m *GenerateEventOperation) GetEvent() *Event {
	if m != nil {
		return m.Event
	}
	return nil
}

// To be used automatically in case of runtime errors (usually bugs)
type ApplicationFailure struct {
	Reason string `protobuf:"bytes,1,opt,name=reason" json:"reason,omitempty"`
}

func (m *ApplicationFailure) Reset()                    { *m = ApplicationFailure{} }
func (m *ApplicationFailure) String() string            { return proto.CompactTextString(m) }
func (*ApplicationFailure) ProtoMessage()               {}
func (*ApplicationFailure) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{6} }

func (m *ApplicationFailure) GetReason() string {
	if m != nil {
		return m.Reason
	}
	return ""
}

// A Event object may model a command (action to executed) or a proper event (action happened in the past)
type Event struct {
	Group   string `protobuf:"bytes,1,opt,name=group" json:"group,omitempty"`
	Name    string `protobuf:"bytes,2,opt,name=name" json:"name,omitempty"`
	Payload []byte `protobuf:"bytes,3,opt,name=payload,proto3" json:"payload,omitempty"`
	// Client identifier is used to match the logged Event with a Transaction in case of fast-path processing
	ClientIdentifier string `protobuf:"bytes,4,opt,name=client_identifier,json=clientIdentifier" json:"client_identifier,omitempty"`
	// Client version should be made visible to client API
	ClientVersion string `protobuf:"bytes,5,opt,name=client_version,json=clientVersion" json:"client_version,omitempty"`
	// Version is meaningful only when event is stored (0 before)
	Version uint64 `protobuf:"varint,6,opt,name=version" json:"version,omitempty"`
}

func (m *Event) Reset()                    { *m = Event{} }
func (m *Event) String() string            { return proto.CompactTextString(m) }
func (*Event) ProtoMessage()               {}
func (*Event) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{7} }

func (m *Event) GetGroup() string {
	if m != nil {
		return m.Group
	}
	return ""
}

func (m *Event) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *Event) GetPayload() []byte {
	if m != nil {
		return m.Payload
	}
	return nil
}

func (m *Event) GetClientIdentifier() string {
	if m != nil {
		return m.ClientIdentifier
	}
	return ""
}

func (m *Event) GetClientVersion() string {
	if m != nil {
		return m.ClientVersion
	}
	return ""
}

func (m *Event) GetVersion() uint64 {
	if m != nil {
		return m.Version
	}
	return 0
}

// A path in the projection store
type Path struct {
	Path    string `protobuf:"bytes,1,opt,name=path" json:"path,omitempty"`
	Version uint64 `protobuf:"varint,2,opt,name=version" json:"version,omitempty"`
}

func (m *Path) Reset()                    { *m = Path{} }
func (m *Path) String() string            { return proto.CompactTextString(m) }
func (*Path) ProtoMessage()               {}
func (*Path) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{8} }

func (m *Path) GetPath() string {
	if m != nil {
		return m.Path
	}
	return ""
}

func (m *Path) GetVersion() uint64 {
	if m != nil {
		return m.Version
	}
	return 0
}

// The object contained in a specific Path
type Data struct {
	Path    *Path                   `protobuf:"bytes,1,opt,name=path" json:"path,omitempty"`
	Content *google_protobuf.Struct `protobuf:"bytes,3,opt,name=content" json:"content,omitempty"`
}

func (m *Data) Reset()                    { *m = Data{} }
func (m *Data) String() string            { return proto.CompactTextString(m) }
func (*Data) ProtoMessage()               {}
func (*Data) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{9} }

func (m *Data) GetPath() *Path {
	if m != nil {
		return m.Path
	}
	return nil
}

func (m *Data) GetContent() *google_protobuf.Struct {
	if m != nil {
		return m.Content
	}
	return nil
}

func init() {
	proto.RegisterType((*Transaction)(nil), "protobuf.Transaction")
	proto.RegisterType((*Operation)(nil), "protobuf.Operation")
	proto.RegisterType((*ReadOperation)(nil), "protobuf.ReadOperation")
	proto.RegisterType((*UpsertOperation)(nil), "protobuf.UpsertOperation")
	proto.RegisterType((*DeleteOperation)(nil), "protobuf.DeleteOperation")
	proto.RegisterType((*GenerateEventOperation)(nil), "protobuf.GenerateEventOperation")
	proto.RegisterType((*ApplicationFailure)(nil), "protobuf.ApplicationFailure")
	proto.RegisterType((*Event)(nil), "protobuf.Event")
	proto.RegisterType((*Path)(nil), "protobuf.Path")
	proto.RegisterType((*Data)(nil), "protobuf.Data")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for DataMesh service

type DataMeshClient interface {
	// Used to push a event that will be stored on the event log.
	Push(ctx context.Context, in *Event, opts ...grpc.CallOption) (*google_protobuf1.Empty, error)
	// Allows to pass a transaction related to a event before it's explicitly requested.
	FastProcess(ctx context.Context, in *Transaction, opts ...grpc.CallOption) (*google_protobuf1.Empty, error)
	// The server sends events that need to be processed. The client reply with the corresponding transactions.
	Process(ctx context.Context, opts ...grpc.CallOption) (DataMesh_ProcessClient, error)
	// Used by the client to query the projections.
	Read(ctx context.Context, in *Path, opts ...grpc.CallOption) (*Data, error)
}

type dataMeshClient struct {
	cc *grpc.ClientConn
}

func NewDataMeshClient(cc *grpc.ClientConn) DataMeshClient {
	return &dataMeshClient{cc}
}

func (c *dataMeshClient) Push(ctx context.Context, in *Event, opts ...grpc.CallOption) (*google_protobuf1.Empty, error) {
	out := new(google_protobuf1.Empty)
	err := grpc.Invoke(ctx, "/protobuf.DataMesh/Push", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *dataMeshClient) FastProcess(ctx context.Context, in *Transaction, opts ...grpc.CallOption) (*google_protobuf1.Empty, error) {
	out := new(google_protobuf1.Empty)
	err := grpc.Invoke(ctx, "/protobuf.DataMesh/FastProcess", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *dataMeshClient) Process(ctx context.Context, opts ...grpc.CallOption) (DataMesh_ProcessClient, error) {
	stream, err := grpc.NewClientStream(ctx, &_DataMesh_serviceDesc.Streams[0], c.cc, "/protobuf.DataMesh/Process", opts...)
	if err != nil {
		return nil, err
	}
	x := &dataMeshProcessClient{stream}
	return x, nil
}

type DataMesh_ProcessClient interface {
	Send(*Transaction) error
	Recv() (*Event, error)
	grpc.ClientStream
}

type dataMeshProcessClient struct {
	grpc.ClientStream
}

func (x *dataMeshProcessClient) Send(m *Transaction) error {
	return x.ClientStream.SendMsg(m)
}

func (x *dataMeshProcessClient) Recv() (*Event, error) {
	m := new(Event)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *dataMeshClient) Read(ctx context.Context, in *Path, opts ...grpc.CallOption) (*Data, error) {
	out := new(Data)
	err := grpc.Invoke(ctx, "/protobuf.DataMesh/Read", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for DataMesh service

type DataMeshServer interface {
	// Used to push a event that will be stored on the event log.
	Push(context.Context, *Event) (*google_protobuf1.Empty, error)
	// Allows to pass a transaction related to a event before it's explicitly requested.
	FastProcess(context.Context, *Transaction) (*google_protobuf1.Empty, error)
	// The server sends events that need to be processed. The client reply with the corresponding transactions.
	Process(DataMesh_ProcessServer) error
	// Used by the client to query the projections.
	Read(context.Context, *Path) (*Data, error)
}

func RegisterDataMeshServer(s *grpc.Server, srv DataMeshServer) {
	s.RegisterService(&_DataMesh_serviceDesc, srv)
}

func _DataMesh_Push_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Event)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DataMeshServer).Push(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/protobuf.DataMesh/Push",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DataMeshServer).Push(ctx, req.(*Event))
	}
	return interceptor(ctx, in, info, handler)
}

func _DataMesh_FastProcess_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Transaction)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DataMeshServer).FastProcess(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/protobuf.DataMesh/FastProcess",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DataMeshServer).FastProcess(ctx, req.(*Transaction))
	}
	return interceptor(ctx, in, info, handler)
}

func _DataMesh_Process_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(DataMeshServer).Process(&dataMeshProcessServer{stream})
}

type DataMesh_ProcessServer interface {
	Send(*Event) error
	Recv() (*Transaction, error)
	grpc.ServerStream
}

type dataMeshProcessServer struct {
	grpc.ServerStream
}

func (x *dataMeshProcessServer) Send(m *Event) error {
	return x.ServerStream.SendMsg(m)
}

func (x *dataMeshProcessServer) Recv() (*Transaction, error) {
	m := new(Transaction)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _DataMesh_Read_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Path)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DataMeshServer).Read(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/protobuf.DataMesh/Read",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DataMeshServer).Read(ctx, req.(*Path))
	}
	return interceptor(ctx, in, info, handler)
}

var _DataMesh_serviceDesc = grpc.ServiceDesc{
	ServiceName: "protobuf.DataMesh",
	HandlerType: (*DataMeshServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Push",
			Handler:    _DataMesh_Push_Handler,
		},
		{
			MethodName: "FastProcess",
			Handler:    _DataMesh_FastProcess_Handler,
		},
		{
			MethodName: "Read",
			Handler:    _DataMesh_Read_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Process",
			Handler:       _DataMesh_Process_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "protobuf/datamesh.proto",
}

func init() { proto.RegisterFile("protobuf/datamesh.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 580 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x94, 0xdd, 0x6e, 0xd3, 0x30,
	0x14, 0xc7, 0x9b, 0x36, 0xfd, 0xc8, 0x29, 0xdb, 0xc0, 0x40, 0x17, 0x2a, 0x2e, 0x8a, 0xa5, 0x49,
	0x95, 0x80, 0x16, 0x5a, 0x10, 0x5c, 0x0d, 0x81, 0xb6, 0x31, 0x2e, 0x10, 0x95, 0xf9, 0xb8, 0x43,
	0xc8, 0x4b, 0xdc, 0x36, 0x5a, 0x1a, 0x47, 0xb6, 0x33, 0x69, 0x2f, 0xc4, 0x3b, 0xf0, 0x36, 0x3c,
	0x0a, 0xb2, 0xf3, 0xd5, 0xa6, 0x13, 0xeb, 0x5d, 0x7c, 0xfe, 0xff, 0x9f, 0x1d, 0xfd, 0xcf, 0xb1,
	0xe1, 0x30, 0x16, 0x5c, 0xf1, 0x8b, 0x64, 0x3e, 0xf6, 0xa9, 0xa2, 0x2b, 0x26, 0x97, 0x23, 0x53,
	0x41, 0x9d, 0x5c, 0xe8, 0xe3, 0x25, 0xf5, 0x2e, 0xe5, 0x78, 0xc1, 0xf9, 0x22, 0x64, 0xe3, 0xc2,
	0x2f, 0x95, 0x48, 0x3c, 0x95, 0xba, 0xfb, 0x4f, 0x6e, 0xf6, 0xb0, 0x55, 0xac, 0xae, 0x53, 0x0b,
	0x0e, 0xa0, 0xfb, 0x4d, 0xd0, 0x48, 0x52, 0x4f, 0x05, 0x3c, 0x42, 0x47, 0xd0, 0x64, 0x57, 0x2c,
	0x52, 0xae, 0x35, 0xb0, 0x86, 0xdd, 0xc9, 0xc1, 0x28, 0x87, 0x46, 0xa7, 0xba, 0x4c, 0x52, 0x15,
	0x4d, 0x01, 0x78, 0xcc, 0x04, 0xd5, 0x8c, 0x74, 0xeb, 0x83, 0xc6, 0xb0, 0x3b, 0xb9, 0x5f, 0x7a,
	0xbf, 0xe4, 0x1a, 0x59, 0xb3, 0xe1, 0xdf, 0x75, 0x70, 0x0a, 0x05, 0x3d, 0x07, 0x5b, 0x30, 0xea,
	0x67, 0x07, 0x1d, 0x96, 0x30, 0x61, 0xd4, 0x2f, 0x6c, 0xe7, 0x35, 0x62, 0x6c, 0x68, 0x0a, 0xad,
	0x24, 0x96, 0x4c, 0x28, 0xb7, 0x6e, 0x80, 0x47, 0x25, 0xf0, 0xdd, 0xd4, 0xd7, 0x91, 0xcc, 0xaa,
	0x21, 0x9f, 0x85, 0x4c, 0x31, 0xb7, 0x51, 0x85, 0x4e, 0x4c, 0x7d, 0x03, 0x4a, 0xad, 0xe8, 0x18,
	0x3a, 0x0b, 0x16, 0xe9, 0x3a, 0x73, 0x6d, 0x83, 0x0d, 0x4a, 0xec, 0x63, 0xa6, 0x98, 0x34, 0xd6,
	0xe9, 0x82, 0x41, 0x6f, 0xa1, 0x3d, 0xa7, 0x41, 0x98, 0x08, 0xe6, 0x36, 0x0d, 0xfe, 0xb8, 0xc4,
	0xdf, 0xc7, 0x71, 0x18, 0x78, 0x86, 0x39, 0x4b, 0x3d, 0xe7, 0x35, 0x92, 0xdb, 0x3f, 0xb4, 0xc0,
	0xbe, 0x0c, 0x22, 0x1f, 0x4f, 0x61, 0x6f, 0x23, 0x04, 0x84, 0xc1, 0x8e, 0xa9, 0x5a, 0x66, 0x59,
	0xed, 0x97, 0xfb, 0xcd, 0xa8, 0x5a, 0x12, 0xa3, 0xe1, 0xd7, 0x70, 0x50, 0x09, 0x42, 0x63, 0x7a,
	0x7c, 0xb6, 0xb1, 0x13, 0xaa, 0x28, 0x31, 0x9a, 0xc6, 0x2a, 0x51, 0xec, 0x74, 0xda, 0x3b, 0xe8,
	0xdd, 0x1c, 0xc5, 0x8e, 0x13, 0x84, 0x9f, 0x01, 0xda, 0x0e, 0x03, 0xf5, 0xa0, 0x25, 0x18, 0x95,
	0x3c, 0x32, 0xb4, 0x43, 0xb2, 0x15, 0xfe, 0x63, 0x41, 0xd3, 0xe0, 0xe8, 0x01, 0x34, 0x17, 0x82,
	0x27, 0x71, 0x66, 0x48, 0x17, 0x08, 0x81, 0x1d, 0xd1, 0x15, 0x33, 0xb3, 0xe1, 0x10, 0xf3, 0x8d,
	0x5c, 0x68, 0xc7, 0xf4, 0x3a, 0xe4, 0xd4, 0x37, 0xdd, 0xbf, 0x43, 0xf2, 0x25, 0x7a, 0x0a, 0xf7,
	0xbc, 0x30, 0x60, 0x91, 0xfa, 0x15, 0xf8, 0x2c, 0x52, 0xc1, 0x3c, 0x60, 0xc2, 0xb4, 0xda, 0x21,
	0x77, 0x53, 0xe1, 0x53, 0x51, 0x47, 0x47, 0xb0, 0x9f, 0x99, 0xaf, 0x98, 0x90, 0x01, 0x8f, 0x4c,
	0x57, 0x1d, 0xb2, 0x97, 0x56, 0x7f, 0xa4, 0x45, 0x7d, 0x5a, 0xae, 0xb7, 0x06, 0xd6, 0xd0, 0x26,
	0xf9, 0x12, 0xbf, 0x02, 0x5b, 0x07, 0xa7, 0xff, 0xb1, 0x88, 0xd5, 0x49, 0x63, 0x5c, 0xa7, 0xea,
	0x9b, 0xd4, 0x4f, 0xb0, 0x75, 0x97, 0x76, 0x69, 0x06, 0x7a, 0x09, 0x6d, 0x8f, 0x47, 0x4a, 0x87,
	0xde, 0xc8, 0x6e, 0x53, 0x7a, 0xe5, 0x4b, 0xf7, 0x57, 0xf3, 0x2c, 0x90, 0xdc, 0x37, 0xf9, 0x6b,
	0x41, 0x47, 0xef, 0xff, 0x99, 0x49, 0xcd, 0xdb, 0xb3, 0x44, 0x2e, 0x51, 0xb5, 0x57, 0xfd, 0xde,
	0xd6, 0x3e, 0xa7, 0xfa, 0xe9, 0xc0, 0x35, 0x74, 0x0c, 0xdd, 0x33, 0x2a, 0xd5, 0x4c, 0x70, 0x8f,
	0x49, 0x89, 0x1e, 0x96, 0x8e, 0xb5, 0xd7, 0xe4, 0x3f, 0xfc, 0x1b, 0x68, 0xdf, 0xc2, 0x56, 0x7f,
	0x06, 0xd7, 0x86, 0xd6, 0x0b, 0x0b, 0x0d, 0xc1, 0xd6, 0x77, 0x03, 0x55, 0x92, 0xe8, 0x57, 0xa6,
	0x1b, 0xd7, 0x2e, 0x5a, 0xa6, 0x30, 0xfd, 0x17, 0x00, 0x00, 0xff, 0xff, 0x3a, 0xaa, 0xbd, 0x5b,
	0x4c, 0x05, 0x00, 0x00,
}
