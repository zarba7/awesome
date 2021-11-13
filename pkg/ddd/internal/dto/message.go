package dto

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
)

type Message struct {
	Code int32
	TrxID   int64
	AggID   int64
	Tid        uint32
	Content     []byte
	Result interface{}
}
const messageHeaderLen = 4 + 8 + 8 + 4
func (msg *Message) Reset()  {
	msg.Content = msg.Content[:]
	msg.Result = nil
	msg.Code = 0
}
func (msg *Message) UnmarshalContent(val interface{}) error {
	return json.Unmarshal(msg.Content, val)
}
func (msg *Message) MarshalContent(val interface{}) (err error) {
	msg.Content, err = json.Marshal(val)
	return
}

func (msg *Message) ResponseMarshal(buf *bytes.Buffer) error {
	if err := json.NewEncoder(buf).Encode(msg.Result); err != nil{
		return err
	}
	return write(buf, msg.TrxID, msg.AggID, msg.Tid, msg.Code)
}
func (msg *Message) ResponseUnmarshal(data []byte) error {
	hl := len(data)-messageHeaderLen
	if err := read(data[hl:], &msg.TrxID, &msg.AggID, &msg.Tid, &msg.Code); err != nil{
		return err
	}
	return json.NewDecoder(bytes.NewBuffer(data[:hl])).Decode(msg.Result)
}
const messageHeaderLen2 = 8 + 8 + 4
func (msg *Message) RequestMarshal(buf *bytes.Buffer) error {
	if _, err := buf.Write(msg.Content); err != nil{
		return err
	}
	return write(buf, msg.TrxID, msg.AggID, msg.Tid)
}
func (msg *Message) RequestUnmarshal(data []byte) error {
	hl := len(data)-messageHeaderLen2
	msg.Content = data[:hl]
	return read(data[hl:], &msg.TrxID, &msg.AggID, &msg.Tid)
}

func write(buf *bytes.Buffer, values... interface{}) error {
	for _, v := range values {
		if err := binary.Write(buf, binary.BigEndian, v); err != nil {
			return err
		}
	}
	return nil
}

func read(data []byte, valuesPtr... interface{}) error {
	buf := bytes.NewBuffer(data)
	for _, v := range valuesPtr {
		if err := binary.Read(buf, binary.BigEndian, v); err != nil {
			return err
		}
	}
	return nil
}
