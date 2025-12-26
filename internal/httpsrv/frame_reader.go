package httpsrv

import (
	"bufio"
	"encoding/binary"
	"errors"
	"io"
)

var (
	ErrChannelTooLarge = errors.New("channel too large")
	ErrDataTooLarge    = errors.New("data too large")
)

type Frame struct {
	Channel string
	DataLen uint32
}

// ReadFrameHeader reads channel and data length from the provided reader. Does not read the actual data.
// A frame header consists of:
// - 1 byte: channel length (N)
// - N bytes: channel name (string)
// - 4 bytes: data length (M)
// Returns a Frame struct and a reader positioned after the header.
func ReadFrameHeader(r io.Reader, maxChannelLen uint8, maxDataLen uint32) (Frame, io.Reader, error) {
	// Reasonable buffer size for HTTP bodies
	// TODO: use a buffer pool to reduce allocations
	br := bufio.NewReaderSize(r, 32*1024)

	// Reads channel length
	var chLen uint8
	if err := binary.Read(br, binary.BigEndian, &chLen); err != nil {
		return Frame{}, br, err
	}
	if chLen == 0 || chLen > maxChannelLen {
		return Frame{}, br, ErrChannelTooLarge
	}

	// Reads channel name
	chBytes := make([]byte, chLen)
	if _, err := io.ReadFull(br, chBytes); err != nil {
		return Frame{}, br, err
	}

	// Reads data length
	var dataLen uint32
	if err := binary.Read(br, binary.BigEndian, &dataLen); err != nil {
		return Frame{}, br, err
	}
	if dataLen > maxDataLen {
		return Frame{}, br, ErrDataTooLarge
	}

	// Return frame
	return Frame{
		Channel: string(chBytes),
		DataLen: dataLen,
	}, br, nil
}
