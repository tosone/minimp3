package main

import (
	"io/ioutil"

	"github.com/hajimehoshi/oto"
	"github.com/tosone/minimp3"
)

func main() {
	var file, _ = ioutil.ReadFile("../test.mp3")
	dec, data, _ := minimp3.DecodeFull(file)

	player, _ := oto.NewPlayer(dec.SampleRate, dec.Channels, 2, 1024)
	player.Write(data)
}
