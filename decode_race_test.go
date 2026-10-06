package minimp3

import (
	"os"
	"sync"
	"testing"
	"time"
)

// TestNewDecoderStarted is a regression test for the data races reported when
// calling NewDecoder, waiting for Started and then reading the decoder
// metadata. It mirrors the reproduction from the issue and is meant to be run
// with the race detector enabled: `go test -race`.
func TestNewDecoderStarted(t *testing.T) {
	file, err := os.Open("./test.mp3")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	dec, err := NewDecoder(file)
	if err != nil {
		t.Fatal(err)
	}
	defer dec.Close()

	select {
	case started := <-dec.Started():
		if !started {
			t.Fatal("decoder was closed before it started")
		}
	case <-time.After(5 * time.Second):
		t.Fatal("timed out waiting for the decoder to start")
	}

	if dec.SampleRate <= 0 {
		t.Errorf("SampleRate = %d, want a positive value", dec.SampleRate)
	}
	if dec.Channels <= 0 {
		t.Errorf("Channels = %d, want a positive value", dec.Channels)
	}
	if dec.Layer <= 0 {
		t.Errorf("Layer = %d, want a positive value", dec.Layer)
	}
}

// TestConcurrentRead exercises concurrent calls to Read. The decoded data and
// the read cursor are shared state, so this fails under `-race` if Read is not
// properly synchronized.
func TestConcurrentRead(t *testing.T) {
	file, err := os.Open("./test.mp3")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	dec, err := NewDecoder(file)
	if err != nil {
		t.Fatal(err)
	}
	defer dec.Close()

	const readers = 4
	var (
		wg    sync.WaitGroup
		mu    sync.Mutex
		total int
	)
	for i := 0; i < readers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			buf := make([]byte, 1024)
			for {
				n, err := dec.Read(buf)
				if n > 0 {
					mu.Lock()
					total += n
					mu.Unlock()
				}
				if err != nil {
					return
				}
			}
		}()
	}
	wg.Wait()

	// test.mp3 decodes to 44928 bytes of PCM data (see issue18).
	if total != 44928 {
		t.Errorf("decoded bytes = %d, want 44928", total)
	}
}
