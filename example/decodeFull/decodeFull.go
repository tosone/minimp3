package main

import (
	"io/ioutil"
	"log"
	"time"

	"github.com/hajimehoshi/oto"
	"github.com/tosone/minimp3"
)

func main() {
	var err error

	var dec *minimp3.Decoder
	var data, file []byte
	var player *oto.Player

	if file, err = ioutil.ReadFile("../test.mp3"); err != nil {
		log.Fatal(err)
	}

	if dec, data, err = minimp3.DecodeFull(file); err != nil {
		log.Fatal(err)
	}

	if player, err = oto.NewPlayer(dec.SampleRate, dec.Channels, 2, 1024); err != nil {
		log.Fatal(err)
	}
	player.Write(data)

	<-time.After(time.Second)
	dec.Close()
	player.Close()
}
