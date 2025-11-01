package unmarshal_test

import (
	// "github.com/stretchr/testify/assert"

	// "github.com/stretchr/testify/require"
	"fmt"
	"testing"

	"github.com/bencode/unmarshal"
	"github.com/stretchr/testify/assert"
	// "github.com/bencode-parser/marshal"
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

func TestInt(t *testing.T) {
	var i float64
	val := []byte("i-3.1e")
	unmarshal.UnMarshaler(val, &i)
	assert.Equal(t, -3.1, i)
}

func TestString(t *testing.T) {
	var i string
	val := []byte("4:four")
	unmarshal.UnMarshaler(val, &i)
	assert.Equal(t, "four", i)
}

func TestArr(t *testing.T) {
	var arr []int
	val := []byte("li2ei4ee")
	unmarshal.UnMarshaler(val, &arr)
	assert.Equal(t, []int{2, 4}, arr)
	fmt.Println("res", arr)
	var arrNest [][]int
	val = []byte("lli2eee")
	unmarshal.UnMarshaler(val, &arrNest)
	assert.Equal(t, [][]int{{2}}, arrNest)
	fmt.Println("res", arrNest)
	var arrStr []string
	valStr := []byte("l5:three4:foure")
	unmarshal.UnMarshaler(valStr, &arrStr)
	fmt.Println("res", arrStr)
	assert.Equal(t, []string{"three", "four"}, arrStr)
}

func TestStruct(t *testing.T) {
	p := Packet{}
	val := []byte("d6:lengthi4e4:data4:eggs4:metald4:info4:INFO6:hashesl5:three4:foureeee")
	err := unmarshal.UnMarshaler(val, &p)
	assert.ErrorIs(t, err, nil)
	fmt.Println("res struct", p)
	assert.Equal(
		t,
		Packet{
			Length: 4,
			Data:   "eggs",
			Meta: []Metadata{
				{Info: "INFO", Hashes: []string{"three", "four"}},
			}},
		p)
}
