// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        (unknown)
// source: ory/keto/relation_tuples/v1alpha2/check_service.proto

package rts

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

// The request for a CheckService.Check RPC.
// Checks whether a specific subject is related to an object.
type CheckRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The namespace to evaluate the check.
	//
	// Note: If you use the expand-API and the check
	// evaluates a RelationTuple specifying a SubjectSet as
	// subject or due to a rewrite rule in a namespace config
	// this check request may involve other namespaces automatically.
	//
	// Deprecated: Do not use.
	Namespace string `protobuf:"bytes,1,opt,name=namespace,proto3" json:"namespace,omitempty"`
	// The related object in this check.
	//
	// Deprecated: Do not use.
	Object string `protobuf:"bytes,2,opt,name=object,proto3" json:"object,omitempty"`
	// The relation between the Object and the Subject.
	//
	// Deprecated: Do not use.
	Relation string `protobuf:"bytes,3,opt,name=relation,proto3" json:"relation,omitempty"`
	// The related subject in this check.
	//
	// Deprecated: Do not use.
	Subject *Subject       `protobuf:"bytes,4,opt,name=subject,proto3" json:"subject,omitempty"`
	Tuple   *RelationTuple `protobuf:"bytes,8,opt,name=tuple,proto3" json:"tuple,omitempty"`
	// This field is not implemented yet and has no effect.
	// <!--
	// Set this field to `true` in case your application
	// needs to authorize depending on up to date ACLs,
	// also called a "content-change check".
	//
	// If set to `true` the `snaptoken` field is ignored,
	// the check is evaluated at the latest snapshot
	// (globally consistent) and the response includes a
	// snaptoken for clients to store along with object
	// contents that can be used for subsequent checks
	// of the same content version.
	//
	// Example use case:
	//   - You need to authorize a user to modify/delete some resource
	//     and it is unacceptable that if the permission to do that had
	//     just been revoked some seconds ago so that the change had not
	//     yet been fully replicated to all availability zones.
	//
	// -->
	Latest bool `protobuf:"varint,5,opt,name=latest,proto3" json:"latest,omitempty"`
	// This field is not implemented yet and has no effect.
	// <!--
	// Optional. Like reads, a check is always evaluated at a
	// consistent snapshot no earlier than the given snaptoken.
	//
	// Leave this field blank if you want to evaluate the check
	// based on eventually consistent ACLs, benefiting from very
	// low latency, but possibly slightly stale results.
	//
	// If the specified token is too old and no longer known,
	// the server falls back as if no snaptoken had been specified.
	//
	// If not specified the server tries to evaluate the check
	// on the best snapshot version where it is very likely that
	// ACLs had already been replicated to all availability zones.
	// -->
	Snaptoken string `protobuf:"bytes,6,opt,name=snaptoken,proto3" json:"snaptoken,omitempty"`
	// The maximum depth to search for a relation.
	//
	// If the value is less than 1 or greater than the global
	// max-depth then the global max-depth will be used instead.
	MaxDepth int32 `protobuf:"varint,7,opt,name=max_depth,json=maxDepth,proto3" json:"max_depth,omitempty"`
}

func (x *CheckRequest) Reset() {
	*x = CheckRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_ory_keto_relation_tuples_v1alpha2_check_service_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CheckRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CheckRequest) ProtoMessage() {}

func (x *CheckRequest) ProtoReflect() protoreflect.Message {
	mi := &file_ory_keto_relation_tuples_v1alpha2_check_service_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CheckRequest.ProtoReflect.Descriptor instead.
func (*CheckRequest) Descriptor() ([]byte, []int) {
	return file_ory_keto_relation_tuples_v1alpha2_check_service_proto_rawDescGZIP(), []int{0}
}

// Deprecated: Do not use.
func (x *CheckRequest) GetNamespace() string {
	if x != nil {
		return x.Namespace
	}
	return ""
}

// Deprecated: Do not use.
func (x *CheckRequest) GetObject() string {
	if x != nil {
		return x.Object
	}
	return ""
}

// Deprecated: Do not use.
func (x *CheckRequest) GetRelation() string {
	if x != nil {
		return x.Relation
	}
	return ""
}

// Deprecated: Do not use.
func (x *CheckRequest) GetSubject() *Subject {
	if x != nil {
		return x.Subject
	}
	return nil
}

func (x *CheckRequest) GetTuple() *RelationTuple {
	if x != nil {
		return x.Tuple
	}
	return nil
}

func (x *CheckRequest) GetLatest() bool {
	if x != nil {
		return x.Latest
	}
	return false
}

func (x *CheckRequest) GetSnaptoken() string {
	if x != nil {
		return x.Snaptoken
	}
	return ""
}

func (x *CheckRequest) GetMaxDepth() int32 {
	if x != nil {
		return x.MaxDepth
	}
	return 0
}

// The response for a CheckService.Check rpc.
type CheckResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Whether the specified subject (id)
	// is related to the requested object.
	//
	// It is false by default if no ACL matches.
	Allowed bool `protobuf:"varint,1,opt,name=allowed,proto3" json:"allowed,omitempty"`
	// This field is not implemented yet and has no effect.
	// <!--
	// The last known snapshot token ONLY specified if
	// the request had not specified a snaptoken,
	// since this performed a "content-change request"
	// and consistently fetched the last known snapshot token.
	//
	// This field is not set if the request had specified a snaptoken!
	//
	// If set, clients should cache and use this token
	// for subsequent requests to have minimal latency,
	// but allow slightly stale responses (only some milliseconds or seconds).
	// -->
	Snaptoken string `protobuf:"bytes,2,opt,name=snaptoken,proto3" json:"snaptoken,omitempty"`
}

func (x *CheckResponse) Reset() {
	*x = CheckResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_ory_keto_relation_tuples_v1alpha2_check_service_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CheckResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CheckResponse) ProtoMessage() {}

func (x *CheckResponse) ProtoReflect() protoreflect.Message {
	mi := &file_ory_keto_relation_tuples_v1alpha2_check_service_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CheckResponse.ProtoReflect.Descriptor instead.
func (*CheckResponse) Descriptor() ([]byte, []int) {
	return file_ory_keto_relation_tuples_v1alpha2_check_service_proto_rawDescGZIP(), []int{1}
}

func (x *CheckResponse) GetAllowed() bool {
	if x != nil {
		return x.Allowed
	}
	return false
}

func (x *CheckResponse) GetSnaptoken() string {
	if x != nil {
		return x.Snaptoken
	}
	return ""
}

var File_ory_keto_relation_tuples_v1alpha2_check_service_proto protoreflect.FileDescriptor

var file_ory_keto_relation_tuples_v1alpha2_check_service_proto_rawDesc = []byte{
	0x0a, 0x35, 0x6f, 0x72, 0x79, 0x2f, 0x6b, 0x65, 0x74, 0x6f, 0x2f, 0x72, 0x65, 0x6c, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x5f, 0x74, 0x75, 0x70, 0x6c, 0x65, 0x73, 0x2f, 0x76, 0x31, 0x61, 0x6c, 0x70,
	0x68, 0x61, 0x32, 0x2f, 0x63, 0x68, 0x65, 0x63, 0x6b, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63,
	0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x21, 0x6f, 0x72, 0x79, 0x2e, 0x6b, 0x65, 0x74,
	0x6f, 0x2e, 0x72, 0x65, 0x6c, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x74, 0x75, 0x70, 0x6c, 0x65,
	0x73, 0x2e, 0x76, 0x31, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x32, 0x1a, 0x37, 0x6f, 0x72, 0x79, 0x2f,
	0x6b, 0x65, 0x74, 0x6f, 0x2f, 0x72, 0x65, 0x6c, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x74, 0x75,
	0x70, 0x6c, 0x65, 0x73, 0x2f, 0x76, 0x31, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x32, 0x2f, 0x72, 0x65,
	0x6c, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x74, 0x75, 0x70, 0x6c, 0x65, 0x73, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x22, 0xd1, 0x02, 0x0a, 0x0c, 0x43, 0x68, 0x65, 0x63, 0x6b, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x12, 0x20, 0x0a, 0x09, 0x6e, 0x61, 0x6d, 0x65, 0x73, 0x70, 0x61, 0x63,
	0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x42, 0x02, 0x18, 0x01, 0x52, 0x09, 0x6e, 0x61, 0x6d,
	0x65, 0x73, 0x70, 0x61, 0x63, 0x65, 0x12, 0x1a, 0x0a, 0x06, 0x6f, 0x62, 0x6a, 0x65, 0x63, 0x74,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x42, 0x02, 0x18, 0x01, 0x52, 0x06, 0x6f, 0x62, 0x6a, 0x65,
	0x63, 0x74, 0x12, 0x1e, 0x0a, 0x08, 0x72, 0x65, 0x6c, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x09, 0x42, 0x02, 0x18, 0x01, 0x52, 0x08, 0x72, 0x65, 0x6c, 0x61, 0x74, 0x69,
	0x6f, 0x6e, 0x12, 0x48, 0x0a, 0x07, 0x73, 0x75, 0x62, 0x6a, 0x65, 0x63, 0x74, 0x18, 0x04, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x2a, 0x2e, 0x6f, 0x72, 0x79, 0x2e, 0x6b, 0x65, 0x74, 0x6f, 0x2e, 0x72,
	0x65, 0x6c, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x74, 0x75, 0x70, 0x6c, 0x65, 0x73, 0x2e, 0x76,
	0x31, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x32, 0x2e, 0x53, 0x75, 0x62, 0x6a, 0x65, 0x63, 0x74, 0x42,
	0x02, 0x18, 0x01, 0x52, 0x07, 0x73, 0x75, 0x62, 0x6a, 0x65, 0x63, 0x74, 0x12, 0x46, 0x0a, 0x05,
	0x74, 0x75, 0x70, 0x6c, 0x65, 0x18, 0x08, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x30, 0x2e, 0x6f, 0x72,
	0x79, 0x2e, 0x6b, 0x65, 0x74, 0x6f, 0x2e, 0x72, 0x65, 0x6c, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x5f,
	0x74, 0x75, 0x70, 0x6c, 0x65, 0x73, 0x2e, 0x76, 0x31, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x32, 0x2e,
	0x52, 0x65, 0x6c, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x54, 0x75, 0x70, 0x6c, 0x65, 0x52, 0x05, 0x74,
	0x75, 0x70, 0x6c, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x6c, 0x61, 0x74, 0x65, 0x73, 0x74, 0x18, 0x05,
	0x20, 0x01, 0x28, 0x08, 0x52, 0x06, 0x6c, 0x61, 0x74, 0x65, 0x73, 0x74, 0x12, 0x1c, 0x0a, 0x09,
	0x73, 0x6e, 0x61, 0x70, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x09, 0x73, 0x6e, 0x61, 0x70, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x12, 0x1b, 0x0a, 0x09, 0x6d, 0x61,
	0x78, 0x5f, 0x64, 0x65, 0x70, 0x74, 0x68, 0x18, 0x07, 0x20, 0x01, 0x28, 0x05, 0x52, 0x08, 0x6d,
	0x61, 0x78, 0x44, 0x65, 0x70, 0x74, 0x68, 0x22, 0x47, 0x0a, 0x0d, 0x43, 0x68, 0x65, 0x63, 0x6b,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x61, 0x6c, 0x6c, 0x6f,
	0x77, 0x65, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x07, 0x61, 0x6c, 0x6c, 0x6f, 0x77,
	0x65, 0x64, 0x12, 0x1c, 0x0a, 0x09, 0x73, 0x6e, 0x61, 0x70, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x73, 0x6e, 0x61, 0x70, 0x74, 0x6f, 0x6b, 0x65, 0x6e,
	0x32, 0x7a, 0x0a, 0x0c, 0x43, 0x68, 0x65, 0x63, 0x6b, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65,
	0x12, 0x6a, 0x0a, 0x05, 0x43, 0x68, 0x65, 0x63, 0x6b, 0x12, 0x2f, 0x2e, 0x6f, 0x72, 0x79, 0x2e,
	0x6b, 0x65, 0x74, 0x6f, 0x2e, 0x72, 0x65, 0x6c, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x74, 0x75,
	0x70, 0x6c, 0x65, 0x73, 0x2e, 0x76, 0x31, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x32, 0x2e, 0x43, 0x68,
	0x65, 0x63, 0x6b, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x30, 0x2e, 0x6f, 0x72, 0x79,
	0x2e, 0x6b, 0x65, 0x74, 0x6f, 0x2e, 0x72, 0x65, 0x6c, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x74,
	0x75, 0x70, 0x6c, 0x65, 0x73, 0x2e, 0x76, 0x31, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x32, 0x2e, 0x43,
	0x68, 0x65, 0x63, 0x6b, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0xc2, 0x01, 0x0a,
	0x24, 0x73, 0x68, 0x2e, 0x6f, 0x72, 0x79, 0x2e, 0x6b, 0x65, 0x74, 0x6f, 0x2e, 0x72, 0x65, 0x6c,
	0x61, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x74, 0x75, 0x70, 0x6c, 0x65, 0x73, 0x2e, 0x76, 0x31, 0x61,
	0x6c, 0x70, 0x68, 0x61, 0x32, 0x42, 0x11, 0x43, 0x68, 0x65, 0x63, 0x6b, 0x53, 0x65, 0x72, 0x76,
	0x69, 0x63, 0x65, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a, 0x3f, 0x67, 0x69, 0x74, 0x68,
	0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x6f, 0x72, 0x79, 0x2f, 0x6b, 0x65, 0x74, 0x6f, 0x2f,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x6f, 0x72, 0x79, 0x2f, 0x6b, 0x65, 0x74, 0x6f, 0x2f, 0x72,
	0x65, 0x6c, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x74, 0x75, 0x70, 0x6c, 0x65, 0x73, 0x2f, 0x76,
	0x31, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x32, 0x3b, 0x72, 0x74, 0x73, 0xaa, 0x02, 0x20, 0x4f, 0x72,
	0x79, 0x2e, 0x4b, 0x65, 0x74, 0x6f, 0x2e, 0x52, 0x65, 0x6c, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x54,
	0x75, 0x70, 0x6c, 0x65, 0x73, 0x2e, 0x76, 0x31, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x32, 0xca, 0x02,
	0x20, 0x4f, 0x72, 0x79, 0x5c, 0x4b, 0x65, 0x74, 0x6f, 0x5c, 0x52, 0x65, 0x6c, 0x61, 0x74, 0x69,
	0x6f, 0x6e, 0x54, 0x75, 0x70, 0x6c, 0x65, 0x73, 0x5c, 0x76, 0x31, 0x61, 0x6c, 0x70, 0x68, 0x61,
	0x32, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_ory_keto_relation_tuples_v1alpha2_check_service_proto_rawDescOnce sync.Once
	file_ory_keto_relation_tuples_v1alpha2_check_service_proto_rawDescData = file_ory_keto_relation_tuples_v1alpha2_check_service_proto_rawDesc
)

func file_ory_keto_relation_tuples_v1alpha2_check_service_proto_rawDescGZIP() []byte {
	file_ory_keto_relation_tuples_v1alpha2_check_service_proto_rawDescOnce.Do(func() {
		file_ory_keto_relation_tuples_v1alpha2_check_service_proto_rawDescData = protoimpl.X.CompressGZIP(file_ory_keto_relation_tuples_v1alpha2_check_service_proto_rawDescData)
	})
	return file_ory_keto_relation_tuples_v1alpha2_check_service_proto_rawDescData
}

var file_ory_keto_relation_tuples_v1alpha2_check_service_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_ory_keto_relation_tuples_v1alpha2_check_service_proto_goTypes = []interface{}{
	(*CheckRequest)(nil),  // 0: ory.keto.relation_tuples.v1alpha2.CheckRequest
	(*CheckResponse)(nil), // 1: ory.keto.relation_tuples.v1alpha2.CheckResponse
	(*Subject)(nil),       // 2: ory.keto.relation_tuples.v1alpha2.Subject
	(*RelationTuple)(nil), // 3: ory.keto.relation_tuples.v1alpha2.RelationTuple
}
var file_ory_keto_relation_tuples_v1alpha2_check_service_proto_depIdxs = []int32{
	2, // 0: ory.keto.relation_tuples.v1alpha2.CheckRequest.subject:type_name -> ory.keto.relation_tuples.v1alpha2.Subject
	3, // 1: ory.keto.relation_tuples.v1alpha2.CheckRequest.tuple:type_name -> ory.keto.relation_tuples.v1alpha2.RelationTuple
	0, // 2: ory.keto.relation_tuples.v1alpha2.CheckService.Check:input_type -> ory.keto.relation_tuples.v1alpha2.CheckRequest
	1, // 3: ory.keto.relation_tuples.v1alpha2.CheckService.Check:output_type -> ory.keto.relation_tuples.v1alpha2.CheckResponse
	3, // [3:4] is the sub-list for method output_type
	2, // [2:3] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_ory_keto_relation_tuples_v1alpha2_check_service_proto_init() }
func file_ory_keto_relation_tuples_v1alpha2_check_service_proto_init() {
	if File_ory_keto_relation_tuples_v1alpha2_check_service_proto != nil {
		return
	}
	file_ory_keto_relation_tuples_v1alpha2_relation_tuples_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_ory_keto_relation_tuples_v1alpha2_check_service_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CheckRequest); i {
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
		file_ory_keto_relation_tuples_v1alpha2_check_service_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CheckResponse); i {
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
			RawDescriptor: file_ory_keto_relation_tuples_v1alpha2_check_service_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_ory_keto_relation_tuples_v1alpha2_check_service_proto_goTypes,
		DependencyIndexes: file_ory_keto_relation_tuples_v1alpha2_check_service_proto_depIdxs,
		MessageInfos:      file_ory_keto_relation_tuples_v1alpha2_check_service_proto_msgTypes,
	}.Build()
	File_ory_keto_relation_tuples_v1alpha2_check_service_proto = out.File
	file_ory_keto_relation_tuples_v1alpha2_check_service_proto_rawDesc = nil
	file_ory_keto_relation_tuples_v1alpha2_check_service_proto_goTypes = nil
	file_ory_keto_relation_tuples_v1alpha2_check_service_proto_depIdxs = nil
}