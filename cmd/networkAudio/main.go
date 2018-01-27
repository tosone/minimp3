package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"

	"github.com/hajimehoshi/oto"
	"github.com/tosone/minimp3"
)

func main() {
	var response, _ = http.Get("http://kuwo-moxuanriji.oss-cn-qingdao.aliyuncs.com/mo.mp3")

	dec, _ := minimp3.NewDecoder(response.Body)
	started := dec.Started()
	<-started
	//fmt.Println(dec)
	player, _ := oto.NewPlayer(dec.SampleRate, dec.Channels, 2, 8192)
	go func() {
		defer response.Body.Close()
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
		fmt.Println("over")
	}()

	signalChanel := make(chan os.Signal, 1)
	signal.Notify(signalChanel, os.Interrupt)
	for {
		select {
		case <-signalChanel:
			return
		}
	}
}
