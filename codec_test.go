package codec_test

import (
	"errors"
	"github.com/greyworlds/codec"
	"github.com/greyworlds/codec/protobuf/baseproto"
	"log"
	"testing"
)

type TestMsg struct {
	Id      int
	Content string
}

func TestWsJsonCodec_Unmarshal(t *testing.T) {
	c := codec.GetCodec("protobuf_codec")
	codec.RegisterMessage("1", (*baseproto.PingPang)(nil))
	//req := baseproto.PingPang{
	//	Timestamp: 123,
	//}
	rdata, err := c.Marshal("1", &baseproto.PingPang{}, nil)
	if err != nil {
		t.Fatal(err)
	}
	_, pack, err := c.Unmarshal(rdata)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(pack.DataPtr, pack.Err)
}

func TestGobCodec(t *testing.T) {
	c := codec.GetCodec("gob_codec")
	codec.RegisterMessage("2", (*TestMsg)(nil))
	//req := TestMsg{
	//	Id:      123,
	//	Content: "测试测试",
	//}
	rdata, err := c.Marshal("2", nil, errors.New("测试错误"))
	if err != nil {
		t.Fatal(err)
	}
	_, pack, err := c.Unmarshal(rdata)
	if err != nil {
		t.Fatal(err)
	}
	log.Println(c.ToString(pack.DataPtr), pack.Err)
}
