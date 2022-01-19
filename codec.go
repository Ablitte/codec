package codec

type Codec interface {
	Marshal(router string, dataPtr interface{}, err error) ([]byte, error)
	Unmarshal([]byte) (int, *MsgPack, error)
	ToString(interface{}) string
}

var codecsList = make(map[string]Codec)

func RegisterCodec(name string, codec Codec) {
	if codec == nil {
		panic("codec: Register provide is nil")
	}
	if _, dup := codecsList[name]; dup {
		panic("codec: Register called twice for provide " + name)
	}
	codecsList[name] = codec
}

func GetCodec(name string) Codec {
	if v, ok := codecsList[name]; ok {
		return v
	}
	return nil
}
