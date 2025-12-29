# bencode

bencode is a small Go library providing encoding and decoding (marshal/unmarshal) support for the Bencode format used by BitTorrent. It implements parsing and serialization of the four bencode types: integers, byte strings, lists and dictionaries. 

## Overview

Bencode is the serialization format used by the BitTorrent protocol. It encodes:
- integers as `i<integer>e`
- byte strings as `<length>:<data>`
- lists as `l<items>e`
- dictionaries as `d<pairs>e` (dictionary keys are byte strings and are typically sorted)

This package provides:
- Decoding (unmarshal) of bencoded data into Go types
- Encoding (marshal) of Go values into bencoded data
- Tests covering typical bencode inputs and edge cases

## Installation

Install the package with:

	go get github.com/shrawankc11/bencode@latest

Then import into your code:

	import "github.com/shrawankc11/bencode"

(Adjust import path/module version as appropriate for your module setup.)

## Quick start

Decode a bencoded text:

```go

type User struct {
	Email   string `bencode:"email"`
	Name    string `bencode:"name"`
    Address string `bencode:"address"`
}

var user User
bencode := []byte("d:email:user@example.com:name:examplee") //no entry for address
reader := bytes.NewReader(bencode)

bencode.Unmarshal(reader, &user)
fmt.Println("#%v", user)

// Expected output (approx): {"email":"user@example", "name":"example", "address": ""} //default string value for address
//Unmarshalling requires a proper declaration of a struct. Partially provided properties are skipped. 
```

Encode from a variable:

```go
type User struct {
	Email   string `bencode:"email"`
	Name    string `bencode:"name"`
    Address string `bencode:"address"`
}

user := User{
    Email: "user@example.com",
    Name: "example",
}

out, err := bencode.Marshal(u)
if err != nil {
	log.Fatal(err)
}
fmt.Printf("%s\n", string(out))
// Expected bencoded string: d4:name7:example5:email16:user@example.come
```

## Testing

Run the unit tests included with the package:

	go test ./...

The repo includes tests for both marshal and unmarshal behavior. Review `marshal_test.go` and `unmarshal_test.go` for examples of expected behavior and edge cases covered by the test suite.

## License
