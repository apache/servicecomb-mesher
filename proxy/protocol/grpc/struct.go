package grpc

import (
	"github.com/gogo/protobuf/proto"
	"google.golang.org/grpc/encoding"
)

type frame struct {
	payload []byte
}

type codec struct {
	protoCodec encoding.Codec
}

//Marshal encode struct into bytes
func (c *codec) Marshal(f interface{}) ([]byte, error) {
	b, ok := f.(*frame)
	if !ok {
		return c.protoCodec.Marshal(f)
	}
	return b.payload, nil

}

//Unmarshal decode bytes into struct
func (c *codec) Unmarshal(d []byte, f interface{}) error {
	b, ok := f.(*frame)
	if !ok {
		return c.protoCodec.Unmarshal(d, f)
	}
	b.payload = d
	return nil
}

//String return name
func (c *codec) String() string {
	return "proxy_codec"
}

type protobufCodec struct {
}

//Marshal encode struct into bytes
func (c *protobufCodec) Marshal(f interface{}) ([]byte, error) {
	return proto.Marshal(f.(proto.Message))

}

//Unmarshal decode bytes into struct
func (c *protobufCodec) Unmarshal(d []byte, f interface{}) error {
	return proto.UnmarshalMerge(d, f.(proto.Message))
}

//Name return name
func (c *protobufCodec) Name() string {
	return "protobuf_codec"
}
