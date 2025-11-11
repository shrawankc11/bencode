package bencode

import (
	"github.com/stretchr/testify/assert"

	// "github.com/stretchr/testify/require"
	"testing"
)

type Address struct {
	Country string  `bencode:"country"`
	Zip     string  `bencode:"zip"`
	Lat     float32 `bencode:"lat"`
	Long    float64 `bencode:"long"`
}

type User struct {
	Name    string  `bencode:"name"`
	Address Address `bencode:"address"`
}

func TestSimple(t *testing.T) {
	var res []byte
	res, _ = Marshal("dinguss")
	assert.Equal(t, res, []byte("7:dinguss"))

	res, _ = Marshal(1)
	assert.Equal(t, res, []byte("i1e"))

	res, _ = Marshal([]int{1, 2, 3})
	assert.Equal(t, res, []byte("li1ei2ei3ee"))

	res, _ = Marshal([]string{"one", "two", "three"})
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
	res, _ := Marshal(p)
	val := []byte("d6:lengthi4e4:data4:eggs4:metald4:info4:INFO6:hashesl5:three4:foureeee")
	assert.Equal(t, res, val)
}
