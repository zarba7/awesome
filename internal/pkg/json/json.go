package json

import (
	"errors"
	jsonIter "github.com/json-iterator/go"
)

var js = jsonIter.ConfigFastest

var (
	MarshalToString     = js.MarshalToString
	Marshal             = js.Marshal
	MarshalIndent       = js.MarshalIndent
	UnmarshalFromString = js.UnmarshalFromString
	Unmarshal           = js.Unmarshal
	Get                 = js.Get
	NewEncoder          = js.NewEncoder
	NewDecoder          = js.NewDecoder
	Valid               = js.Valid
	RegisterExtension   = js.RegisterExtension
	DecoderOf           = js.DecoderOf
	EncoderOf           = js.EncoderOf
)

func MarshalPanic(val interface{}) []byte {
	data, err := js.Marshal(val)
	if !errors.Is(err, nil) {
		panic(err)
	}
	return data
}

func UnmarshalPanic(data []byte, v interface{}) {
	if 0 == len(data) {
		return
	}
	err := js.Unmarshal(data, v)
	if !errors.Is(err, nil) {
		panic(err)
	}
}
