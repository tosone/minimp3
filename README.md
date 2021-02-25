# minimp3

Decode mp3 base on <https://github.com/lieff/minimp3>

[![Build Status](https://travis-ci.org/tosone/minimp3.svg?branch=master)](https://travis-ci.org/tosone/minimp3) [![Coverage Status](https://coveralls.io/repos/github/tosone/minimp3/badge.svg)](https://coveralls.io/github/tosone/minimp3) [![GoDoc](https://godoc.org/github.com/tosone/minimp3?status.svg)](https://godoc.org/github.com/tosone/minimp3)

See examples in example directory. `make` and `make test` test the example.

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
