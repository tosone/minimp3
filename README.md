# minimp3

[![Go Reference](https://pkg.go.dev/badge/github.com/tosone/minimp3.svg)](https://pkg.go.dev/github.com/tosone/minimp3) [![Builder](https://github.com/tosone/minimp3/actions/workflows/ci.yaml/badge.svg?branch=main)](https://github.com/tosone/minimp3/actions/workflows/ci.yaml) [![codecov](https://codecov.io/gh/tosone/minimp3/branch/main/graph/badge.svg?token=LUIF0jZw6E)](https://codecov.io/gh/tosone/minimp3)

Decode mp3 base on <https://github.com/lieff/minimp3>

## Installation

1. The first need Go installed (version 1.15+ is required), then you can use the below Go command to install minimp3.

``` bash
$ go get -u github.com/tosone/minimp3
```

2. Import it in your code:

``` bash
import "github.com/tosone/minimp3"
```

## Examples are here

<details>
  <summary>Example1: Decode the whole mp3 and play.</summary>

``` golang
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

	var file []byte
	if file, err = ioutil.ReadFile("test.mp3"); err != nil {
		log.Fatal(err)
	}

	var dec *minimp3.Decoder
	var data []byte
	if dec, data, err = minimp3.DecodeFull(file); err != nil {
		log.Fatal(err)
	}

	var context *oto.Context
	if context, err = oto.NewContext(dec.SampleRate, dec.Channels, 2, 1024); err != nil {
		log.Fatal(err)
	}

	var player = context.NewPlayer()
	player.Write(data)

	<-time.After(time.Second)

	dec.Close()
	if err = player.Close(); err != nil {
		log.Fatal(err)
	}
}
```

</details>

<details>
  <summary>Example2: Decode and play.</summary>

``` go
package main

import (
	"log"
	"os"
	"time"

	"github.com/ebitengine/oto/v3"
	"github.com/tosone/minimp3"
)

func main() {
	file, err := os.Open("input.mp3")
	if err != nil {
		log.Fatal(err)
	}
	dec, err := minimp3.NewDecoder(file)
	if err != nil {
		log.Fatal(err)
	}
	<-dec.Started()
	log.Printf("Convert audio sample rate: %d, channels: %d\n", dec.SampleRate, dec.Channels)
	op := &oto.NewContextOptions{
		SampleRate:   dec.SampleRate,
		ChannelCount: dec.Channels,
		Format:       oto.FormatSignedInt16LE,
	}
	context, readyChan, err := oto.NewContext(op)
	if err != nil {
		log.Fatal(err)
	}
	<-readyChan
	player := context.NewPlayer(dec)
	player.Play()
	for player.IsPlaying() {
		time.Sleep(time.Millisecond)
	}
	if err := player.Close(); err != nil {
		log.Fatal(err)
	}
}
```

</details>

<details>
  <summary>Example3: Play the network audio.</summary>

``` go
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

	var args = os.Args
	if len(args) != 2 {
		log.Fatal("Run test like this:\n\n\t./networkAudio.test [mp3url]\n\n")
	}

	var response *http.Response
	if response, err = http.Get(args[1]); err != nil {
		log.Fatal(err)
	}

	var dec *minimp3.Decoder
	if dec, err = minimp3.NewDecoder(response.Body); err != nil {
		log.Fatal(err)
	}
	<-dec.Started()

	log.Printf("Convert audio sample rate: %d, channels: %d\n", dec.SampleRate, dec.Channels)

	var context *oto.Context
	if context, err = oto.NewContext(dec.SampleRate, dec.Channels, 2, 4096); err != nil {
		log.Fatal(err)
	}

	var waitForPlayOver = new(sync.WaitGroup)
	waitForPlayOver.Add(1)

	var player = context.NewPlayer()

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
```

</details>
