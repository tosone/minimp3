# minimp3

Decode mp3 base on https://github.com/lieff/minimp3

[![Build Status](https://travis-ci.org/tosone/minimp3.svg?branch=master)](https://travis-ci.org/tosone/minimp3) [![Coverage Status](https://coveralls.io/repos/github/tosone/minimp3/badge.svg)](https://coveralls.io/github/tosone/minimp3) [![GoDoc](https://godoc.org/github.com/tosone/minimp3?status.svg)](https://godoc.org/github.com/tosone/minimp3)

See examples in example directory. `make` and `make test` test the example.

``` golang
package main

import (
	"io/ioutil"

	"github.com/hajimehoshi/oto"
	"github.com/tosone/minimp3"
)

func main() {
	var file, _ = ioutil.ReadFile("test.mp3")
	dec, data, _ := minimp3.DecodeFull(file)

	player, _ := oto.NewPlayer(dec.SampleRate, dec.Channels, 2, 1024)
	player.Write(data)
}
```
