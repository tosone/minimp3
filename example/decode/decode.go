package main

import (
	"io"
	"log"
	"os"
	"sync"
	"time"

	"github.com/hajimehoshi/oto"
	"github.com/tosone/minimp3"
)

func main() {
	var err error
	var file *os.File
	var dec *minimp3.Decoder
	var player *oto.Player

	if file, err = os.Open("../test.mp3"); err != nil {
		log.Fatal(err)
	}
	if dec, err = minimp3.NewDecoder(file); err != nil {
		log.Fatal(err)
	}
	started := dec.Started()
	<-started

	log.Printf("Convert audio sample rate: %d, channels: %d\n", dec.SampleRate, dec.Channels)

	if player, err = oto.NewPlayer(dec.SampleRate, dec.Channels, 2, 1024); err != nil {
		log.Fatal(err)
	}

	var waitForPlayOver = new(sync.WaitGroup)
	waitForPlayOver.Add(1)

	go func() {
		for {
			var data = make([]byte, 1024)
			_, err := dec.Read(data)
			if err == io.EOF {
				break
			}
			if err != nil {
				break
			}
			player.Write(data)
		}
		log.Println("over play.")
		waitForPlayOver.Done()
	}()
	waitForPlayOver.Wait()

	<-time.After(time.Second)
	dec.Close()
	player.Close()
}
