package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/hajimehoshi/oto"
	"github.com/tosone/minimp3"
)

func main() {
	var err error
	var response *http.Response
	var dec *minimp3.Decoder
	var player *oto.Player

	args := os.Args
	if len(args) != 2 {
		log.Fatal("Run test like this:\n\n\t./networkAudio.test [mp3url]\n\n")
	}
	if response, err = http.Get(args[1]); err != nil {
		log.Fatal(err)
	}

	if dec, err = minimp3.NewDecoder(response.Body); err != nil {
		log.Fatal(err)
	}
	<-dec.Started()

	log.Printf("Convert audio sample rate: %d, channels: %d\n", dec.SampleRate, dec.Channels)

	if player, err = oto.NewPlayer(dec.SampleRate, dec.Channels, 2, 4096); err != nil {
		log.Fatal(err)
	}

	var waitForPlayOver = new(sync.WaitGroup)
	waitForPlayOver.Add(1)

	go func() {
		defer response.Body.Close()
		for {
			var data = make([]byte, 512)
			_, err = dec.Read(data)
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatal(err)
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
