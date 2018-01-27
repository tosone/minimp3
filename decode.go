package minimp3

/*
#define MINIMP3_IMPLEMENTATION

#include "minimp3.h"
#include <stdlib.h>
#include <stdio.h>

int decode(mp3dec_t *dec, mp3dec_frame_info_t *info, unsigned char *data, int *length, unsigned char *decoded, int *decoded_length) {
	int samples;
    short pcm[MINIMP3_MAX_SAMPLES_PER_FRAME];
    samples = mp3dec_decode_frame(dec, data, *length, pcm, info);
    *decoded_length = samples * info->channels;
    *length -= info->frame_bytes;
	unsigned char buffer[samples*info->channels*2];
	memcpy(buffer, (unsigned char*)&(pcm), sizeof(short) * samples * info->channels);
	memcpy(decoded, buffer, sizeof(short) * samples * info->channels);
    return info->frame_bytes;
}
*/
import "C"
import (
	"context"
	"io"
	"sync"
	"time"
	"unsafe"
)

const maxSamplesPerFrame = 1152 * 2

type Decoder struct {
	sync.Mutex
	data          []byte
	readIndex     int
	decodedData   []byte
	decodeIndex   int
	decode        C.mp3dec_t
	info          C.mp3dec_frame_info_t
	context       context.Context
	contextCancel context.CancelFunc
	SampleRate    int
	Channels      int
	Kbps          int
	Layer         int
}

func NewDecoder(reader io.Reader) (dec *Decoder, err error) {
	dec = new(Decoder)
	dec.context, dec.contextCancel = context.WithCancel(context.Background())
	dec.decode = C.mp3dec_t{}
	C.mp3dec_init(&dec.decode)
	dec.info = C.mp3dec_frame_info_t{}
	go func() {
		for {
			select {
			case <-dec.context.Done():
				break
			default:
			}
			var data = make([]byte, 1024)
			var n int
			n, err = reader.Read(data)
			if err == io.EOF {
				return
			}
			if err != nil {
				return
			}
			dec.Lock()
			dec.data = append(dec.data, data[:n]...)
			dec.Unlock()
		}
	}()
	go func() {
		for {
			select {
			case <-dec.context.Done():
				break
			default:
			}
			var decoded = [maxSamplesPerFrame * 2]byte{}
			var decodedLength = C.int(0)
			var length = C.int(len(dec.data))
			if len(dec.data) == 0 {
				continue
			}
			frameSize := C.decode(&dec.decode, &dec.info,
				(*C.uchar)(unsafe.Pointer(&dec.data[0])),
				&length, (*C.uchar)(unsafe.Pointer(&decoded[0])),
				&decodedLength)
			if int(frameSize) == 0 {
				<-time.After(time.Millisecond * 100)
				continue
			}
			dec.SampleRate = int(dec.info.hz)
			dec.Channels = int(dec.info.channels)
			dec.Kbps = int(dec.info.bitrate_kbps)
			dec.Layer = int(dec.info.layer)
			dec.decodeIndex += int(frameSize)
			dec.Lock()
			dec.decodedData = append(dec.decodedData, decoded[:decodedLength*2]...)
			if int(frameSize) < len(dec.data) {
				dec.data = dec.data[int(frameSize):]
			}
			dec.Unlock()
		}
	}()
	return
}

func (dec *Decoder) Started() (channel chan bool) {
	channel = make(chan bool)
	go func() {
		for {
			select {
			case <-dec.context.Done():
				channel <- false
			default:
			}
			if len(dec.decodedData) != 0 {
				channel <- true
			} else {
				<-time.After(time.Millisecond * 100)
			}
		}
	}()
	return
}

func (dec *Decoder) Read(data []byte) (n int, err error) {
	dec.Lock()
	defer dec.Unlock()
	if dec.readIndex > len(dec.decodedData) {
		err = io.EOF
		return
	}
	n = copy(data, dec.decodedData[dec.readIndex:])
	dec.readIndex += n
	return
}

func (dec *Decoder) Close() {
	if dec.contextCancel != nil {
		dec.contextCancel()
	}
}

func DecodeFull(mp3 []byte) (dec *Decoder, decodedData []byte, err error) {
	dec = new(Decoder)
	dec.decode = C.mp3dec_t{}
	C.mp3dec_init(&dec.decode)
	info := C.mp3dec_frame_info_t{}
	var length = C.int(len(mp3))
	for {
		var decoded = [maxSamplesPerFrame * 2]byte{}
		var decodedLength = C.int(0)
		frameSize := C.decode(&dec.decode,
			&info, (*C.uchar)(unsafe.Pointer(&mp3[0])),
			&length, (*C.uchar)(unsafe.Pointer(&decoded[0])),
			&decodedLength)
		if int(frameSize) == 0 {
			break
		}
		decodedData = append(decodedData, decoded[:decodedLength*2]...)
		if int(frameSize) < len(mp3) {
			mp3 = mp3[int(frameSize):]
		}
		dec.SampleRate = int(info.hz)
		dec.Channels = int(info.channels)
		dec.Kbps = int(info.bitrate_kbps)
		dec.Layer = int(info.layer)
	}
	return
}
