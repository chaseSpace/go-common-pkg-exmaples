// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.5.1
// source: protocol/sub1/item.proto

package sub1

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

type ItemDetailStatus int32

const (
	ItemDetail_DEFAULT  ItemDetailStatus = 0
	ItemDetail_ACTIVE   ItemDetailStatus = 1
	ItemDetail_INACTIVE ItemDetailStatus = 2
	ItemDetail_OFF      ItemDetailStatus = 3
)

// Enum value maps for ItemDetailStatus.
var (
	ItemDetailStatus_name = map[int32]string{
		0: "DEFAULT",
		1: "ACTIVE",
		2: "INACTIVE",
		3: "OFF",
	}
	ItemDetailStatus_value = map[string]int32{
		"DEFAULT":  0,
		"ACTIVE":   1,
		"INACTIVE": 2,
		"OFF":      3,
	}
)

func (x ItemDetailStatus) Enum() *ItemDetailStatus {
	p := new(ItemDetailStatus)
	*p = x
	return p
}

func (x ItemDetailStatus) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (ItemDetailStatus) Descriptor() protoreflect.EnumDescriptor {
	return file_protocol_sub1_item_proto_enumTypes[0].Descriptor()
}

func (ItemDetailStatus) Type() protoreflect.EnumType {
	return &file_protocol_sub1_item_proto_enumTypes[0]
}

func (x ItemDetailStatus) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use ItemDetailStatus.Descriptor instead.
func (ItemDetailStatus) EnumDescriptor() ([]byte, []int) {
	return file_protocol_sub1_item_proto_rawDescGZIP(), []int{0, 0}
}

type ItemDetail struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name   string           `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Price  int32            `protobuf:"varint,2,opt,name=price,proto3" json:"price,omitempty"`
	Desc   string           `protobuf:"bytes,3,opt,name=desc,proto3" json:"desc,omitempty"`
	Status ItemDetailStatus `protobuf:"varint,4,opt,name=Status,proto3,enum=sub1.ItemDetailStatus" json:"Status,omitempty"`
}

func (x *ItemDetail) Reset() {
	*x = ItemDetail{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protocol_sub1_item_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ItemDetail) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ItemDetail) ProtoMessage() {}

func (x *ItemDetail) ProtoReflect() protoreflect.Message {
	mi := &file_protocol_sub1_item_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ItemDetail.ProtoReflect.Descriptor instead.
func (*ItemDetail) Descriptor() ([]byte, []int) {
	return file_protocol_sub1_item_proto_rawDescGZIP(), []int{0}
}

func (x *ItemDetail) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *ItemDetail) GetPrice() int32 {
	if x != nil {
		return x.Price
	}
	return 0
}

func (x *ItemDetail) GetDesc() string {
	if x != nil {
		return x.Desc
	}
	return ""
}

func (x *ItemDetail) GetStatus() ItemDetailStatus {
	if x != nil {
		return x.Status
	}
	return ItemDetail_DEFAULT
}

var File_protocol_sub1_item_proto protoreflect.FileDescriptor

var file_protocol_sub1_item_proto_rawDesc = []byte{
	0x0a, 0x18, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6c, 0x2f, 0x73, 0x75, 0x62, 0x31, 0x2f,
	0x69, 0x74, 0x65, 0x6d, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x04, 0x73, 0x75, 0x62, 0x31,
	0x22, 0xb5, 0x01, 0x0a, 0x0a, 0x49, 0x74, 0x65, 0x6d, 0x44, 0x65, 0x74, 0x61, 0x69, 0x6c, 0x12,
	0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e,
	0x61, 0x6d, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x70, 0x72, 0x69, 0x63, 0x65, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x05, 0x52, 0x05, 0x70, 0x72, 0x69, 0x63, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x64, 0x65, 0x73,
	0x63, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x64, 0x65, 0x73, 0x63, 0x12, 0x2f, 0x0a,
	0x06, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x17, 0x2e,
	0x73, 0x75, 0x62, 0x31, 0x2e, 0x49, 0x74, 0x65, 0x6d, 0x44, 0x65, 0x74, 0x61, 0x69, 0x6c, 0x2e,
	0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x52, 0x06, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x22, 0x38,
	0x0a, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x0b, 0x0a, 0x07, 0x44, 0x45, 0x46, 0x41,
	0x55, 0x4c, 0x54, 0x10, 0x00, 0x12, 0x0a, 0x0a, 0x06, 0x41, 0x43, 0x54, 0x49, 0x56, 0x45, 0x10,
	0x01, 0x12, 0x0c, 0x0a, 0x08, 0x49, 0x4e, 0x41, 0x43, 0x54, 0x49, 0x56, 0x45, 0x10, 0x02, 0x12,
	0x07, 0x0a, 0x03, 0x4f, 0x46, 0x46, 0x10, 0x03, 0x42, 0x19, 0x5a, 0x17, 0x65, 0x78, 0x61, 0x6d,
	0x70, 0x6c, 0x65, 0x5f, 0x67, 0x65, 0x6e, 0x5f, 0x70, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x73,
	0x75, 0x62, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_protocol_sub1_item_proto_rawDescOnce sync.Once
	file_protocol_sub1_item_proto_rawDescData = file_protocol_sub1_item_proto_rawDesc
)

func file_protocol_sub1_item_proto_rawDescGZIP() []byte {
	file_protocol_sub1_item_proto_rawDescOnce.Do(func() {
		file_protocol_sub1_item_proto_rawDescData = protoimpl.X.CompressGZIP(file_protocol_sub1_item_proto_rawDescData)
	})
	return file_protocol_sub1_item_proto_rawDescData
}

var file_protocol_sub1_item_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_protocol_sub1_item_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_protocol_sub1_item_proto_goTypes = []interface{}{
	(ItemDetailStatus)(0), // 0: sub1.ItemDetail.status
	(*ItemDetail)(nil),    // 1: sub1.ItemDetail
}
var file_protocol_sub1_item_proto_depIdxs = []int32{
	0, // 0: sub1.ItemDetail.Status:type_name -> sub1.ItemDetail.status
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_protocol_sub1_item_proto_init() }
func file_protocol_sub1_item_proto_init() {
	if File_protocol_sub1_item_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_protocol_sub1_item_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ItemDetail); i {
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
			RawDescriptor: file_protocol_sub1_item_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_protocol_sub1_item_proto_goTypes,
		DependencyIndexes: file_protocol_sub1_item_proto_depIdxs,
		EnumInfos:         file_protocol_sub1_item_proto_enumTypes,
		MessageInfos:      file_protocol_sub1_item_proto_msgTypes,
	}.Build()
	File_protocol_sub1_item_proto = out.File
	file_protocol_sub1_item_proto_rawDesc = nil
	file_protocol_sub1_item_proto_goTypes = nil
	file_protocol_sub1_item_proto_depIdxs = nil
}
