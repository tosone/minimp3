# minimp3

[![Go Reference](https://pkg.go.dev/badge/github.com/tosone/minimp3.svg)](https://pkg.go.dev/github.com/tosone/minimp3) [![Builder](https://github.com/tosone/minimp3/actions/workflows/ci.yaml/badge.svg?branch=main)](https://github.com/tosone/minimp3/actions/workflows/ci.yaml) [![Coverage Status](https://coveralls.io/repos/github/tosone/minimp3/badge.svg)](https://coveralls.io/github/tosone/minimp3)

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
	"bytes"
	"log"
	"os"
	"time"

	"github.com/ebitengine/oto/v3"
	"github.com/tosone/minimp3"
)

func main() {
	file, err := os.ReadFile("test.mp3")
	if err != nil {
		log.Fatal(err)
	}

	dec, data, err := minimp3.DecodeFull(file)
	if err != nil {
		log.Fatal(err)
	}

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

	player := context.NewPlayer(bytes.NewReader(data))
	player.Play()
	for player.IsPlaying() {
		time.Sleep(time.Millisecond)
	}

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
	"log"
	"net/http"
	"os"
	"time"

	"github.com/ebitengine/oto/v3"
	"github.com/tosone/minimp3"
)

func main() {
	var args = os.Args
	if len(args) != 2 {
		log.Fatal("Run test like this:\n\n\t./networkAudio.test [mp3url]\n\n")
	}

	response, err := http.Get(args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	dec, err := minimp3.NewDecoder(response.Body)
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
	dec.Close()
	if err := player.Close(); err != nil {
		log.Fatal(err)
	}
}
```

</details>
