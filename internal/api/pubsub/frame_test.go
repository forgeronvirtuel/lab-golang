package pubsub

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"testing"
)

func TestReadFrameHeader_Success(t *testing.T) {
	tests := []struct {
		name          string
		channel       string
		dataLen       uint32
		maxChannelLen uint8
		maxDataLen    uint32
	}{
		{
			name:          "simple frame",
			channel:       "test",
			dataLen:       100,
			maxChannelLen: 255,
			maxDataLen:    1024,
		},
		{
			name:          "single char channel",
			channel:       "a",
			dataLen:       1,
			maxChannelLen: 10,
			maxDataLen:    1000,
		},
		{
			name:          "max channel length",
			channel:       string(make([]byte, 255)),
			dataLen:       0,
			maxChannelLen: 255,
			maxDataLen:    1000,
		},
		{
			name:          "max data length",
			channel:       "channel",
			dataLen:       999,
			maxChannelLen: 255,
			maxDataLen:    999,
		},
		{
			name:          "large data length",
			channel:       "events",
			dataLen:       1048576,
			maxChannelLen: 100,
			maxDataLen:    10485760,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create frame data
			buf := new(bytes.Buffer)
			chLen := uint8(len(tt.channel))
			binary.Write(buf, binary.BigEndian, chLen)
			buf.WriteString(tt.channel)
			binary.Write(buf, binary.BigEndian, tt.dataLen)

			// Read frame
			frame, reader, err := ReadFrameHeader(buf, tt.maxChannelLen, tt.maxDataLen)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			// Verify frame
			if frame.ChannelName != tt.channel {
				t.Errorf("channel = %q, want %q", frame.ChannelName, tt.channel)
			}
			if frame.DataLen != tt.dataLen {
				t.Errorf("dataLen = %d, want %d", frame.DataLen, tt.dataLen)
			}
			if reader == nil {
				t.Error("reader is nil")
			}
		})
	}
}

func TestReadFrameHeader_ChannelTooLarge(t *testing.T) {
	tests := []struct {
		name          string
		chLen         uint8
		maxChannelLen uint8
		wantErr       error
	}{
		{
			name:          "channel length is zero",
			chLen:         0,
			maxChannelLen: 10,
			wantErr:       ErrChannelTooLarge,
		},
		{
			name:          "channel length exceeds max",
			chLen:         11,
			maxChannelLen: 10,
			wantErr:       ErrChannelTooLarge,
		},
		{
			name:          "channel length equals max plus one",
			chLen:         101,
			maxChannelLen: 100,
			wantErr:       ErrChannelTooLarge,
		},
		{
			name:          "channel length is 255 but max is 254",
			chLen:         255,
			maxChannelLen: 254,
			wantErr:       ErrChannelTooLarge,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := new(bytes.Buffer)
			binary.Write(buf, binary.BigEndian, tt.chLen)

			_, _, err := ReadFrameHeader(buf, tt.maxChannelLen, 1000)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("error = %v, want %v", err, tt.wantErr)
			}
		})
	}
}

func TestReadFrameHeader_DataTooLarge(t *testing.T) {
	tests := []struct {
		name       string
		dataLen    uint32
		maxDataLen uint32
	}{
		{
			name:       "data length exceeds max by 1",
			dataLen:    1001,
			maxDataLen: 1000,
		},
		{
			name:       "data length much larger than max",
			dataLen:    1000000,
			maxDataLen: 1000,
		},
		{
			name:       "max uint32 data length",
			dataLen:    ^uint32(0),
			maxDataLen: 1000,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := new(bytes.Buffer)
			chLen := uint8(4)
			binary.Write(buf, binary.BigEndian, chLen)
			buf.WriteString("test")
			binary.Write(buf, binary.BigEndian, tt.dataLen)

			_, _, err := ReadFrameHeader(buf, 255, tt.maxDataLen)
			if !errors.Is(err, ErrDataTooLarge) {
				t.Errorf("error = %v, want %v", err, ErrDataTooLarge)
			}
		})
	}
}

func TestReadFrameHeader_IOErrors(t *testing.T) {
	t.Run("empty reader", func(t *testing.T) {
		buf := new(bytes.Buffer)
		_, _, err := ReadFrameHeader(buf, 10, 1000)
		if err == nil {
			t.Error("expected error for empty reader")
		}
		if !errors.Is(err, io.EOF) {
			t.Errorf("error = %v, want io.EOF", err)
		}
	})

	t.Run("truncated channel length", func(t *testing.T) {
		buf := new(bytes.Buffer)
		// Empty buffer will cause EOF when reading channel length
		_, _, err := ReadFrameHeader(buf, 10, 1000)
		if err == nil {
			t.Error("expected error for truncated data")
		}
	})

	t.Run("truncated channel data", func(t *testing.T) {
		buf := new(bytes.Buffer)
		chLen := uint8(10)
		binary.Write(buf, binary.BigEndian, chLen)
		buf.WriteString("short") // only 5 bytes, but chLen says 10

		_, _, err := ReadFrameHeader(buf, 255, 1000)
		if err == nil {
			t.Error("expected error for truncated channel data")
		}
		if !errors.Is(err, io.ErrUnexpectedEOF) && !errors.Is(err, io.EOF) {
			t.Errorf("error = %v, want io.ErrUnexpectedEOF or io.EOF", err)
		}
	})

	t.Run("truncated data length", func(t *testing.T) {
		buf := new(bytes.Buffer)
		chLen := uint8(4)
		binary.Write(buf, binary.BigEndian, chLen)
		buf.WriteString("test")
		// Don't write dataLen - will cause EOF

		_, _, err := ReadFrameHeader(buf, 255, 1000)
		if err == nil {
			t.Error("expected error for truncated data length")
		}
		if !errors.Is(err, io.EOF) && !errors.Is(err, io.ErrUnexpectedEOF) {
			t.Errorf("error = %v, want io.EOF or io.ErrUnexpectedEOF", err)
		}
	})

	t.Run("partial data length", func(t *testing.T) {
		buf := new(bytes.Buffer)
		chLen := uint8(4)
		binary.Write(buf, binary.BigEndian, chLen)
		buf.WriteString("test")
		buf.Write([]byte{0x00, 0x01}) // only 2 bytes of 4-byte uint32

		_, _, err := ReadFrameHeader(buf, 255, 1000)
		if err == nil {
			t.Error("expected error for partial data length")
		}
	})
}

func TestReadFrameHeader_ReaderReturned(t *testing.T) {
	// Test that the returned reader can be used to read remaining data
	buf := new(bytes.Buffer)
	chLen := uint8(4)
	binary.Write(buf, binary.BigEndian, chLen)
	buf.WriteString("test")
	dataLen := uint32(12)
	binary.Write(buf, binary.BigEndian, dataLen)
	expectedData := []byte("hello world!")
	buf.Write(expectedData)

	frame, reader, err := ReadFrameHeader(buf, 255, 1000)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if frame.DataLen != dataLen {
		t.Errorf("dataLen = %d, want %d", frame.DataLen, dataLen)
	}

	// Read the actual data using the returned reader
	actualData := make([]byte, dataLen)
	n, err := io.ReadFull(reader, actualData)
	if err != nil {
		t.Fatalf("failed to read data: %v", err)
	}
	if uint32(n) != dataLen {
		t.Errorf("read %d bytes, want %d", n, dataLen)
	}
	if !bytes.Equal(actualData, expectedData) {
		t.Errorf("data = %q, want %q", actualData, expectedData)
	}
}

func TestReadFrameHeader_BoundaryValues(t *testing.T) {
	t.Run("channel length at boundary", func(t *testing.T) {
		buf := new(bytes.Buffer)
		chLen := uint8(10)
		binary.Write(buf, binary.BigEndian, chLen)
		buf.WriteString("0123456789")
		binary.Write(buf, binary.BigEndian, uint32(100))

		frame, _, err := ReadFrameHeader(buf, 10, 1000)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if frame.ChannelName != "0123456789" {
			t.Errorf("channel = %q, want %q", frame.ChannelName, "0123456789")
		}
	})

	t.Run("data length at boundary", func(t *testing.T) {
		buf := new(bytes.Buffer)
		chLen := uint8(4)
		binary.Write(buf, binary.BigEndian, chLen)
		buf.WriteString("test")
		dataLen := uint32(1000)
		binary.Write(buf, binary.BigEndian, dataLen)

		frame, _, err := ReadFrameHeader(buf, 255, 1000)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if frame.DataLen != dataLen {
			t.Errorf("dataLen = %d, want %d", frame.DataLen, dataLen)
		}
	})

	t.Run("zero data length", func(t *testing.T) {
		buf := new(bytes.Buffer)
		chLen := uint8(5)
		binary.Write(buf, binary.BigEndian, chLen)
		buf.WriteString("empty")
		binary.Write(buf, binary.BigEndian, uint32(0))

		frame, _, err := ReadFrameHeader(buf, 255, 1000)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if frame.DataLen != 0 {
			t.Errorf("dataLen = %d, want 0", frame.DataLen)
		}
	})
}
