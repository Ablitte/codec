package gob

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
	"github.com/greyworlds/codec"
)

type gobPacket struct {
	Router string
	Data   []byte
	Err    string
}

type gobCodec struct{}

func (*gobCodec) Marshal(router string, dataPtr interface{}, retErr error) ([]byte, error) {
	if router == "" {
		return nil, fmt.Errorf("marshal: empty router")
	}
	if dataPtr == nil && retErr == nil {
		return nil, fmt.Errorf("marshal: empty data")
	}
	ack := &gobPacket{
		Router: router,
	}
	if dataPtr != nil {
		data := bytes.NewBuffer(nil)
		encoder1 := gob.NewEncoder(data)
		err := encoder1.Encode(dataPtr)
		if err != nil {
			return nil, fmt.Errorf("marshal:protocol buffer marshal failed. router:%s dt:%T err:%v",
				router, dataPtr, err)
		}
		ack.Data = data.Bytes()
	}
	if retErr != nil {
		ack.Err = retErr.Error()
	}
	ackByte := bytes.NewBuffer(nil)
	encoder2 := gob.NewEncoder(ackByte)
	err := encoder2.Encode(&ack)
	if err != nil {
		return nil, fmt.Errorf("marshal:gob marshal failed. router:%s dt:%T err:%v",
			router, ack, err)
	}
	return ackByte.Bytes(), nil
}
func (*gobCodec) Unmarshal(msg []byte) (int, *codec.MsgPack, error) {
	var len = len(msg)
	req := &gobPacket{}
	err := gob.NewDecoder(bytes.NewReader(msg)).Decode(req)
	if err != nil {
		return len, nil, fmt.Errorf("unmarshal split message id failed.")
	}
	var router = req.Router
	msgPack := &codec.MsgPack{Router: router}
	dt := codec.GetMessage(router)
	if dt == nil {
		return len, nil, fmt.Errorf("unmarshal message not registed. router:%s",
			router)
	}
	if req.Data != nil {
		err = gob.NewDecoder(bytes.NewReader(req.Data)).Decode(dt)
		if err != nil {
			return len, nil, fmt.Errorf("unmarshal failed. router:%s", router)
		}
	}
	msgPack.DataPtr = dt
	if req.Err != "" {
		msgPack.Err = errors.New(req.Err)
	}
	return len, msgPack, nil
}

func (*gobCodec) ToString(data interface{}) string {
	return fmt.Sprintf("type:%T, data:%v", data, data)
}

func init() {
	codec.RegisterCodec("gob_codec", new(gobCodec))
}
