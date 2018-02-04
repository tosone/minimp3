package minimp3

import (
	"io"
	"io/ioutil"
	"log"
	"os"
	"sync"
	"testing"
	"time"
)

func TestDecoder(t *testing.T) {
	var err error
	var file, pcmFile *os.File
	var dec *Decoder

	if file, err = os.Open("./example/test.mp3"); err != nil {
		t.Error(err)
	}
	if dec, err = NewDecoder(file); err != nil {
		t.Error(err)
	}
	started := dec.Started()
	<-started

	log.Printf("Convert audio sample rate: %d, channels: %d\n", dec.SampleRate, dec.Channels)

	if pcmFile, err = os.Create("test1.pcm"); err != nil {
		t.Error(err)
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

			if _, err = pcmFile.Write(data); err != nil {
				t.Error(err)
			}
		}
		log.Println("over play.")
		waitForPlayOver.Done()
	}()
	waitForPlayOver.Wait()

	<-time.After(time.Second)
	dec.Close()
	if err = pcmFile.Close(); err != nil {
		t.Error(err)
	}
}

func TestDecodeFull(t *testing.T) {
	var err error

	var dec *Decoder
	var data, file []byte

	if file, err = ioutil.ReadFile("./example/test.mp3"); err != nil {
		t.Error(err)
	}

	if dec, data, err = DecodeFull(file); err != nil {
		t.Error(err)
	}

	ioutil.WriteFile("test2.pcm", data, 0644)

	<-time.After(time.Second)
	dec.Close()
}

func TestDecoder_CloseEarly1(t *testing.T) {
	var err error
	var file *os.File
	var dec *Decoder

	if file, err = os.Open("./example/test.mp3"); err != nil {
		t.Error(err)
	}
	if dec, err = NewDecoder(file); err != nil {
		t.Error(err)
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
	var file, pcmFile *os.File
	var dec *Decoder

	if file, err = os.Open("./example/test.mp3"); err != nil {
		t.Error(err)
	}
	if dec, err = NewDecoder(file); err != nil {
		t.Error(err)
	}

	started := dec.Started()
	<-started

	dec.Close()

	log.Printf("Convert audio sample rate: %d, channels: %d\n", dec.SampleRate, dec.Channels)

	if pcmFile, err = os.Create("test4.pcm"); err != nil {
		t.Error(err)
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
			if _, err = pcmFile.Write(data); err != nil {
				t.Error(err)
			}
		}
		log.Println("over play.")
		waitForPlayOver.Done()
	}()
	waitForPlayOver.Wait()

	<-time.After(time.Second)
	dec.Close()
	if err = pcmFile.Close(); err != nil {
		t.Error(err)
	}
}
