package bencode

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

type Metadata struct {
	Info   string   `bencode:"info"`
	Hashes []string `bencode:"hashes"`
}

type Packet struct {
	Length  int        `bencode:"length"`
	Data    string     `bencode:"data"`
	Meta    []Metadata `bencode:"meta"`
	TwoDArr [][]string `bencode:"al"`
}

type TorrentInfo struct {
	Pieces      string `bencode:"pieces"`
	PieceLength int    `bencode:"piece length"`
	Name        string `bencode:"name"`
	Files       []File `bencode:"files"`
}

type File struct {
	Length int      `bencode:"length"`
	Path   []string `bencode:"path"`
}

type Torrent struct {
	Announce     string      `bencode:"announce"`
	AnnounceList [][]string  `bencode:"announce-list"`
	Info         TorrentInfo `bencode:"info"`
	Comment      string      `bencode:"comment"`
	CreatedBy    string      `bencode:"created by"`
	CreationDate int         `bencode:"creation date"`
}

func TestInt(t *testing.T) {
	var i float64
	file, _:= os.Open("file.torrent")
	UnMarshal(file, &i)
	assert.Equal(t, -3.1, i)
}

// func TestString(t *testing.T) {
// 	var i string
// 	val := []byte("4:four")
// 	UnMarshal(val, &i)
// 	assert.Equal(t, "four", i)
// }
//
// func TestArr(t *testing.T) {
// 	var arr []int
// 	val := []byte("li2ei4ee")
// 	UnMarshal(val, &arr)
// 	assert.Equal(t, []int{2, 4}, arr)
// 	var arrNest [][]int
// 	val = []byte("lli2eee")
// 	UnMarshal(val, &arrNest)
// 	assert.Equal(t, [][]int{{2}}, arrNest)
// 	var arrStr []string
// 	valStr := []byte("l5:three4:foure")
// 	UnMarshal(valStr, &arrStr)
// 	assert.Equal(t, []string{"three", "four"}, arrStr)
// }
//
// func TestUnMarshalStruct(t *testing.T) {
// 	p := Packet{}
// 	val := []byte("d6:lengthi4e4:data4:eggs4:metald4:info4:INFO6:hashesl5:three4:foureee2:alll3:helel2:loeee")
// 	err := UnMarshal(val, &p)
// 	assert.ErrorIs(t, err, nil)
// 	assert.Equal(
// 		t,
// 		Packet{
// 			Length: 4,
// 			Data:   "eggs",
// 			Meta: []Metadata{
// 				{Info: "INFO", Hashes: []string{"three", "four"}},
// 			},
// 			TwoDArr: [][]string{{"hel"}, {"lo"}},
// 		},
// 		p)
// 	data, _ := json.Marshal(p)
// 	fmt.Println(string(data))
// }
//
// func TestTorrent(t *testing.T) {
// 	torrent := Torrent{}
// 	data, _ := os.ReadFile("/home/shrawan/Downloads/52DBDF7C9DABD1EE320FF4D90A6A6F211A4E0113.torrent")
// 	UnMarshal(data, &torrent)
// 	fmt.Println(torrent)
// }
