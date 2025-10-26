package marshal_test

import (
	"fmt"

	"github.com/stretchr/testify/assert"

	// "github.com/stretchr/testify/require"
	"testing"

	"github.com/bencode-parser/marshal"
)

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
	res, _ = marshal.Marshal(user)
	fmt.Println(string(res))

	res, _ = marshal.Marshal("dinguss")
	assert.Equal(t, res, []byte("7:dinguss"))

	res, _ = marshal.Marshal(1)
	assert.Equal(t, res, []byte("i1e"))

	res, _ = marshal.Marshal([]int{1, 2, 3})
	assert.Equal(t, res, []byte("li1ei2ei3ee"))

	res, _ = marshal.Marshal([]string{"one", "two", "three"})
	assert.Equal(t, res, []byte("l3:one3:two5:threee"))
}
