package decoder

//
// v1.0.0
//
type Message struct {
	Ver                  int32    `protobuf:"varint,1,opt,name=ver,proto3" json:"ver,omitempty"`
	Op                   int32    `protobuf:"varint,2,opt,name=op,proto3" json:"op,omitempty"`
	Seq                  int32    `protobuf:"varint,3,opt,name=seq,proto3" json:"seq,omitempty"`
	Body                 []byte   `protobuf:"bytes,4,opt,name=body,proto3" json:"body,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}
