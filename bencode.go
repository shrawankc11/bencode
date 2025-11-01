package bencode

import (
	"github.com/bencode/unmarshal"
	"github.com/bencode/marshal"
)


func UnMarshal(e []byte, val any) (err error) {
	return unmarshal.UnMarshaler(e, val)
}

func Marshal(val any) ([]byte, error) {
	return marshal.Marshaler(val)
}
