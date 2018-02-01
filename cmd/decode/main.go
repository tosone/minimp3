package main

import (
	"io"
	"os"
	"os/signal"

	"fmt"

	"github.com/hajimehoshi/oto"
	"github.com/tosone/minimp3"
)

func main() {
	var file, _ = os.Open("./test.mp3")
	dec, _ := minimp3.NewDecoder(file)
	started := dec.Started()
	<-started
	fmt.Println(dec.SampleRate, dec.Channels)
	player, _ := oto.NewPlayer(dec.SampleRate, dec.Channels, 2, 512)
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
	}()

	signalChanel := make(chan os.Signal, 1)
	signal.Notify(signalChanel, os.Interrupt)
	for {
		select {
		case <-signalChanel:
			file.Close()
			return
		}
	}
}
