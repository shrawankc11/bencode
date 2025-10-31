package marshal_test

import (

	"github.com/stretchr/testify/assert"

	// "github.com/stretchr/testify/require"
	"testing"

	"github.com/bencode-parser/marshal"
)

type Metadata struct {
	Info   string   `bencode:"info"`
	Hashes []string `bencode:"hashes"`
}

type Packet struct {
	Length int        `bencode:"length"`
	Data   string     `bencode:"data"`
	Meta   []Metadata `bencode:"meta"`
}

type Address struct {
	Country string
	Zip     string
	Lat     float32
	Long    float64
}

type User struct {
	Name    string `bencode:"name"`
	Address Address
}

func TestReflect(t *testing.T) {
	user := User{Name: "Dingus", Address: Address{Zip: "1234", Country: "Nepal", Lat: 13.232, Long: 12.12312}}
	var res []byte
	res, _ = marshal.Marshaler(user)

	res, _ = marshal.Marshaler("dinguss")
	assert.Equal(t, res, []byte("7:dinguss"))

	res, _ = marshal.Marshaler(1)
	assert.Equal(t, res, []byte("i1e"))

	res, _ = marshal.Marshaler([]int{1, 2, 3})
	assert.Equal(t, res, []byte("li1ei2ei3ee"))

	res, _ = marshal.Marshaler([]string{"one", "two", "three"})
	assert.Equal(t, res, []byte("l3:one3:two5:threee"))
}

func TestWhole(t *testing.T) {
	p := Packet{
			Length: 4,
			Data:   "eggs",
			Meta: []Metadata{
				{Info: "INFO", Hashes: []string{"three", "four"}},
			},
	}
	res, _ := marshal.Marshaler(p)
	val := []byte("d6:lengthi4e4:data4:eggs4:metald4:info4:INFO6:hashesl5:three4:foureeee")
	assert.Equal(t, res, val)
}
