package protobuf

import (
	"errors"
	"fmt"
	"github.com/greyworlds/codec"
	"github.com/greyworlds/codec/protobuf/baseproto"

	"github.com/golang/protobuf/proto"
)

type protobufCodec struct{}

func (*protobufCodec) Marshal(router string, dataPtr interface{}, retErr error) ([]byte, error) {
	if router == "" {
		return nil, fmt.Errorf("marshal: empty router")
	}
	if dataPtr == nil && retErr == nil {
		return nil, fmt.Errorf("marshal: empty data")
	}
	ack := &baseproto.TransPack{
		Id: router,
	}
	if dataPtr != nil {
		pbMsg, ok := dataPtr.(proto.Message)
		if !ok {
			return nil, fmt.Errorf("marshal: dataptr only support proto.Message type. router:%s dt:%T ",
				router, dataPtr)
		}
		data, err := proto.Marshal(pbMsg)
		if err != nil {
			return nil, fmt.Errorf("marshal:protocol buffer marshal failed. router:%s dt:%T err:%v",
				router, dataPtr, err)
		}
		ack.Data = data
	} else {
		ack.Error = retErr.Error()
	}
	ackByte, err := proto.Marshal(ack)
	if err != nil {
		return nil, fmt.Errorf("marshal:protocol buffer marshal failed. router:%s dt:%T err:%v",
			router, ack, err)
	}
	return ackByte, nil
}

func (*protobufCodec) Unmarshal(msg []byte) (int, *codec.MsgPack, error) {
	var len = len(msg)
	req := &baseproto.TransPack{}
	err := proto.Unmarshal(msg, req)
	if err != nil {
		return len, nil, errors.New("unmarshal split message id failed.")
	}
	var router = req.Id
	msgPack := &codec.MsgPack{Router: router}
	dt := codec.GetMessage(router)
	if dt == nil {
		return len, nil, fmt.Errorf("unmarshal message not registed. router:%s",
			router)
	}
	if req.Data != nil {
		err = proto.Unmarshal(req.Data, dt.(proto.Message))
		if err != nil {
			return len, nil, fmt.Errorf("unmarshal failed. router:%s", router)
		}
	}
	msgPack.DataPtr = dt
	if req.Error != "" {
		msgPack.Err = errors.New(req.Error)
	}
	return len, msgPack, nil
}

func (*protobufCodec) ToString(data interface{}) string {
	pbMsg, ok := data.(proto.Message)
	if !ok {
		return fmt.Sprintf("invalid type %T", data)
	}
	return pbMsg.String()
}

func init() {
	codec.RegisterCodec("protobuf_codec", new(protobufCodec))
}
