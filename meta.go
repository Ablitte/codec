package codec

import (
	"fmt"
	"reflect"
	"sync"
)

type MsgPack struct {
	Router  interface{}
	DataPtr interface{}
	Err     error
}

var modelMap = make(map[interface{}]reflect.Type)
var modelMapLock sync.RWMutex

func RegisterMessage(router interface{}, datePtr interface{}) {
	modelMapLock.Lock()
	defer modelMapLock.Unlock()
	if _, ok := modelMap[router]; ok {
		fmt.Println(fmt.Sprintf("codec: repeat registration. router:%s ", router))
		return
	}
	if t, ok := datePtr.(reflect.Type); ok {
		modelMap[router] = t.Elem()
	} else {
		t := reflect.TypeOf(datePtr)
		if t.Kind() != reflect.Ptr {
			panic(fmt.Errorf("codec: cannot use non-ptr message struct `%s`", t))
		}
		modelMap[router] = t.Elem()
	}
}

func GetMessage(router interface{}) interface{} {
	modelMapLock.RLock()
	defer modelMapLock.RUnlock()
	if ptr, ok := modelMap[router]; ok {
		return reflect.New(ptr).Interface()
	}
	return nil
}
