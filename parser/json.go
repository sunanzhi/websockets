package parser

import (
	"encoding/json"
	"fmt"
	"reflect"
	"websockets"
)

type JsonProtocol struct {
	types   map[int64]reflect.Type
	names   map[reflect.Type]int64
	handler map[int64]websockets.ParserHandler
}

func Json() websockets.Protocol {
	return &JsonProtocol{
		types: make(map[int64]reflect.Type),
		names: make(map[reflect.Type]int64),
		handler: make(map[int64]websockets.ParserHandler),
	}
}

func (j *JsonProtocol) Register(protocol int64, t interface{}, handler websockets.ParserHandler) {
	rt := reflect.TypeOf(t)
	if rt.Kind() == reflect.Ptr {
		rt = rt.Elem()
	}
	j.types[protocol] = rt
	j.names[rt] = protocol
	j.handler[protocol] = handler
}

type jsonIn struct {
	Protocol int64
	Body     *json.RawMessage
}

type jsonOut struct {
	Protocol int64
	Body     interface{}
}

func (j *JsonProtocol) Receive(c *websockets.Client, msg []byte) error {
	var in jsonIn
	err := json.Unmarshal(msg, &in)
	if err != nil {
		return err
	}

	if t, ok := j.types[in.Protocol]; ok {
		body := reflect.New(t).Interface()
		fmt.Println("in.body:", string(*in.Body))
		err = json.Unmarshal(*in.Body, &body)
		if err != nil {
			return err
		}
		if handler, ok := j.handler[in.Protocol]; ok && handler != nil {
			handler(c, body)
		}
	}

	return nil
}

func (j *JsonProtocol) Send(c *websockets.Client, msg interface{}) error {
	var out jsonOut
	t := reflect.TypeOf(msg)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if protocol, exists := j.names[t]; exists {
		out.Protocol = protocol
	}
	out.Body = msg

	res, err := json.Marshal(out)
	if err != nil {
		return err
	}
	c.Send <- res

	return nil
}
