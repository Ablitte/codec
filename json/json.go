package json

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/greywords/codec"
)

type jsonCodec struct{}

type jsonPack struct {
	Router string
	Data   []byte
	Err    string
}

func (*jsonCodec) Marshal(router string, dataPtr interface{}, retErr error) ([]byte, error) {

	if router == "" {
		return nil, fmt.Errorf("marshal: router is empty")
	}

	if dataPtr == nil && retErr == nil {
		return nil, fmt.Errorf("marshal data in package is nil. router:%s dt:%T",
			router, dataPtr)
	}
	ack := &jsonPack{
		Router: router,
	}
	if dataPtr != nil {
		data, err := json.Marshal(dataPtr)
		if err != nil {
			return nil, fmt.Errorf("marshal json marshal failed. routerr:%s dt:%T err:%v",
				router, dataPtr, err)
		}
		ack.Data = data
	}
	if retErr != nil {
		ack.Err = retErr.Error()
	}
	ackByte, err := json.Marshal(ack)
	if err != nil {
		return nil, fmt.Errorf("marshal json marshal failed. routerr:%s dt:%T err:%v",
			router, dataPtr, err)
	}
	return ackByte, nil
}
func (*jsonCodec) Unmarshal(msg []byte) (int, *codec.MsgPack, error) {
	fmt.Sprintln(string(msg))
	var len = len(msg)
	req := &jsonPack{}
	err := json.Unmarshal(msg, req)
	if err != nil {
		return len, nil, errors.New("unmarshal split message id failed.")
	}
	var router = req.Router
	msgPack := &codec.MsgPack{Router: router}
	dt := codec.GetMessage(router)
	if dt == nil {
		return len, nil, fmt.Errorf("unmarshal message not registed. router:%s",
			router)
	}
	if req.Data != nil {
		err = json.Unmarshal(req.Data, dt)
		if err != nil {
			return len, nil, fmt.Errorf("unmarshal json unmarshal failed. dt:%T msg:%s err:%v",
				dt, string(msg), err)
		}
	}
	msgPack.DataPtr = dt
	if req.Err != "" {
		msgPack.Err = errors.New(req.Err)
	}
	return len, msgPack, nil
}

func (*jsonCodec) ToString(data interface{}) string {
	ab, err := json.Marshal(data)
	if err != nil {
		return fmt.Sprintf("invalid type %T", data)
	}
	return string(ab)
}

func init() {
	codec.RegisterCodec("json_codec", new(jsonCodec))
}
