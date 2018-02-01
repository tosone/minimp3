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

//type safeData struct {
//	locker *sync.Mutex
//	data   []byte
//}

//var waitForPlay = safeData{}

func main() {
	//fmt.Println(audio.AudioInit())
	//waitForPlay.locker = new(sync.Mutex)
	var response, _ = http.Get("http://zhangmenshiting.qianqian.com/data2/music/248543902/248543902.mp3?xcode=68992518c55697fac70098e904ce32a9")

	dec, _ := minimp3.NewDecoder(response.Body)
	started := dec.Started()
	<-started
	fmt.Println(dec.SampleRate)
	player, _ := oto.NewPlayer(dec.SampleRate, dec.Channels, 2, 10240)
	player.SetUnderrunCallback(func() {

	})
	go func() {
		defer response.Body.Close()
		for {
			var data = make([]byte, 8192)
			n, err := dec.Read(data)
			if err == io.EOF {
				break
			}
			if err != nil {
				break
			}
			if n == 0 {
				break
			}
			//waitForPlay.locker.Lock()
			//waitForPlay.data = append(waitForPlay.data, data...)
			//waitForPlay.locker.Unlock()
			fmt.Println(len(data))
			player.Write(data)
			//if channelEnd, err := audio.PlaySE(data); err != nil {
			//	fmt.Println(err)
			//} else {
			//	<-channelEnd
			//}
		}
		fmt.Println("over")
	}()

	//go func() {
	//	var audioCVT = new(sdl.AudioCVT)
	//	sdl.BuildAudioCVT(audioCVT, sdl.AUDIO_S16, uint8(dec.Channels), dec.SampleRate, sdl.AUDIO_S16, 1, 16000)
	//	audioCVT.AllocBuf(uintptr(8192))
	//	for {
	//		var data []byte
	//		if len(waitForPlay.data) == 0 {
	//			<-time.After(time.Millisecond * 100)
	//			continue
	//		}
	//		if len(waitForPlay.data) < 8192 {
	//			data = waitForPlay.data
	//		} else {
	//			data = waitForPlay.data[:8192]
	//		}
	//		audioCVT.Len = int32(len(data))
	//		audioCVT.Buf = unsafe.Pointer(&waitForPlay.data[0])
	//		sdl.ConvertAudio(audioCVT)
	//		var converted []byte
	//		sliceHeader := (*reflect.SliceHeader)(unsafe.Pointer(&converted))
	//		sliceHeader.Len = int(audioCVT.LenCVT)
	//		sliceHeader.Cap = int(audioCVT.Len * audioCVT.LenMult)
	//		sliceHeader.Data = uintptr(unsafe.Pointer(audioCVT.Buf))
	//
	//		audio.PlaySE(converted)
	//	}
	//	audioCVT.FreeBuf()
	//}()

	signalChanel := make(chan os.Signal, 1)
	signal.Notify(signalChanel, os.Interrupt)
	for {
		select {
		case <-signalChanel:
			return
		}
	}
}
