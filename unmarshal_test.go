package bencode_test

import (
	"github.com/shrawankc11/bencode"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
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
	bencode.UnMarshal(val, &i)
	assert.Equal(t, -3.1, i)
}

func TestString(t *testing.T) {
	var i string
	val := []byte("4:four")
	bencode.UnMarshal(val, &i)
	assert.Equal(t, "four", i)
}

func TestArr(t *testing.T) {
	var arr []int
	val := []byte("li2ei4ee")
	bencode.UnMarshal(val, &arr)
	assert.Equal(t, []int{2, 4}, arr)
	fmt.Println("res", arr)
	var arrNest [][]int
	val = []byte("lli2eee")
	bencode.UnMarshal(val, &arrNest)
	assert.Equal(t, [][]int{{2}}, arrNest)
	fmt.Println("res", arrNest)
	var arrStr []string
	valStr := []byte("l5:three4:foure")
	bencode.UnMarshal(valStr, &arrStr)
	fmt.Println("res", arrStr)
	assert.Equal(t, []string{"three", "four"}, arrStr)
}

func TestUnMarshalStruct(t *testing.T) {
	p := Packet{}
	val := []byte("d6:lengthi4e4:data4:eggs4:metald4:info4:INFO6:hashesl5:three4:foureeee")
	err := bencode.UnMarshal(val, &p)
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
