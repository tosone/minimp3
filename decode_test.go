package minimp3

import (
	"io"
	"io/ioutil"
	"log"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/hajimehoshi/oto"
)

func TestDecoder(t *testing.T) {
	var err error
	var file *os.File
	var dec *Decoder
	var player *oto.Player

	if file, err = os.Open("./example/test.mp3"); err != nil {
		t.Fatal(err)
	}
	if dec, err = NewDecoder(file); err != nil {
		t.Fatal(err)
	}
	started := dec.Started()
	<-started

	log.Printf("Convert audio sample rate: %d, channels: %d\n", dec.SampleRate, dec.Channels)

	if player, err = oto.NewPlayer(dec.SampleRate, dec.Channels, 2, 1024); err != nil {
		t.Fatal(err)
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
			if _, err = player.Write(data); err != nil {
				t.Fatal(err)
			}
		}
		log.Println("over play.")
		waitForPlayOver.Done()
	}()
	waitForPlayOver.Wait()

	<-time.After(time.Second)
	dec.Close()
	if err = player.Close(); err != nil {
		t.Fatal(err)
	}
}

func TestDecodeFull(t *testing.T) {
	var err error

	var dec *Decoder
	var data, file []byte
	var player *oto.Player

	if file, err = ioutil.ReadFile("./example/test.mp3"); err != nil {
		t.Fatal(err)
	}

	if dec, data, err = DecodeFull(file); err != nil {
		t.Fatal(err)
	}

	if player, err = oto.NewPlayer(dec.SampleRate, dec.Channels, 2, 1024); err != nil {
		t.Fatal(err)
	}
	if _, err = player.Write(data); err != nil {
		t.Fatal(err)
	}

	<-time.After(time.Second)
	dec.Close()
	if err = player.Close(); err != nil {
		t.Fatal(err)
	}
}

func TestDecoder_CloseEarly1(t *testing.T) {
	var err error
	var file *os.File
	var dec *Decoder

	if file, err = os.Open("./example/test.mp3"); err != nil {
		t.Fatal(err)
	}
	if dec, err = NewDecoder(file); err != nil {
		t.Fatal(err)
	}
	dec.Close()
	started := dec.Started()
	<-started
	if dec.SampleRate != 0 || dec.Channels != 0 {
		t.Error("minimp3 decoder cannot be closed correctly.")
	}
}

func TestDecoder_CloseEarly2(t *testing.T) {
	var err error
	var file *os.File
	var dec *Decoder
	var player *oto.Player

	if file, err = os.Open("./example/test.mp3"); err != nil {
		t.Fatal(err)
	}
	if dec, err = NewDecoder(file); err != nil {
		t.Fatal(err)
	}

	started := dec.Started()
	<-started

	dec.Close()

	log.Printf("Convert audio sample rate: %d, channels: %d\n", dec.SampleRate, dec.Channels)

	if player, err = oto.NewPlayer(dec.SampleRate, dec.Channels, 2, 1024); err != nil {
		t.Fatal(err)
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
			if _, err = player.Write(data); err != nil {
				t.Fatal(err)
			}
		}
		log.Println("over play.")
		waitForPlayOver.Done()
	}()
	waitForPlayOver.Wait()

	<-time.After(time.Second)
	if err = player.Close(); err != nil {
		t.Fatal(err)
	}
}
