package minimp3

import (
	"bytes"
	"io"
	"os"
	"testing"
	"time"
)

// TestIssue18DelayPutData ...
func TestIssue18DelayPutData(t *testing.T) {
	reader, writer := io.Pipe()
	dec, err := NewDecoder(reader)
	if err != nil {
		t.Errorf("NewDecoder failed: %v", err)
	}
	go func() {
		time.Sleep(3 * time.Second)
		file, err := os.ReadFile("test.mp3")
		if err != nil {
			t.Errorf("open file failed: %v", err)
		}
		_, err = io.Copy(writer, bytes.NewReader(file))
		if err != nil {
			t.Errorf("copy mp3 data to pipe failed: %v", err)
		}
		writer.Close() // nolint: errcheck
	}()
	// seems like the 'dec.Started' is unnecessary here
	data, err := io.ReadAll(dec)
	if err != nil {
		t.Errorf("read the whole decoded data failed: %v", err)
	}
	if len(data) != 44928 {
		t.Errorf("decode mp3 file failed, real is 44928, but got %d", len(data))
	}
}

// TestIssue18GracefulExit ...
func TestIssue18GracefulExit(t *testing.T) {
	reader, writer := io.Pipe()
	dec, err := NewDecoder(reader)
	if err != nil {
		t.Errorf("NewDecoder failed: %v", err)
	}
	go func() {
		time.Sleep(3 * time.Second)
		writer.Close() // nolint: errcheck
	}()
	// seems like the 'dec.Started' is necessary here
	// if the Decoder input reader is closed, then the Decoder.Read will be returned
	data, err := io.ReadAll(dec)
	if err != nil {
		t.Errorf("read the whole decoded data failed: %v", err)
	}
	if len(data) != 0 {
		t.Errorf("graceful exit read something data")
	}
}
