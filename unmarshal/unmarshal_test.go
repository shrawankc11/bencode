package unmarshal_test

import (
	// "github.com/stretchr/testify/assert"

	// "github.com/stretchr/testify/require"
	"fmt"
	"testing"

	"github.com/bencode-parser/unmarshal"
	"github.com/stretchr/testify/assert"
	// "github.com/bencode-parser/marshal"
)

type Metadata struct {
	Info   string   `bencode:"info"`
	Hashes []string `bencode:"hashes"`
}

type Packet struct {
	Length int      `bencode:"length"`
	Data   string   `bencode:"data"`
	Meta   Metadata `bencode:"meta"`
}

func TestInt(t *testing.T) {
	var i float64
	val := []byte("i-3.1e")
	unmarshal.UnMarshal(val, &i)
	assert.Equal(t, -3.1, i)
}

func TestString(t *testing.T) {
	var i string
	val := []byte("4:four")
	unmarshal.UnMarshal(val, &i)
	assert.Equal(t, "four", i)
}

func TestArr(t *testing.T) {
	var arr []int
	val := []byte("li2ei4ee")
	unmarshal.UnMarshal(val, &arr)
	assert.Equal(t, []int{2, 4}, arr)
	var arrStr []string
	valStr := []byte("l5:three4:foure")
	_, res := unmarshal.UnMarshal(valStr, &arrStr)
	fmt.Println(res)
	assert.Equal(t, []string{"three", "four"}, arrStr)
}

func TestStruct(t *testing.T) {
	var err error
	p := Packet{}
	val := []byte("d6:lengthi4e4:data4:eggs4:metad4:info4:INFO6:hashesl5:three4:foureee")
	// val := []byte("d6:lengthi4e4:data4:eggse")
	err, res := unmarshal.UnMarshal(val, &p)
	assert.Equal(t, Packet{Length: 4, Data: "eggs", Meta: Metadata{Info: "INFO", Hashes: []string{"three", "four"}}}, res)
	fmt.Println(err)
	fmt.Println(res)
}
