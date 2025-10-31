package bencode

import (
	"github.com/bencode-parser/unmarshal"
	"github.com/bencode-parser/marshal"
)


func UnMarshal(e []byte, val any) (err error) {
	return unmarshal.UnMarshaler(e, val)
}

func Marshal(val any) ([]byte, error) {
	return marshal.Marshaler(val)
}
