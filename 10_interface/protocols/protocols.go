package protocols

import (
	"fmt"
	"os"
)

type Protocol interface {
	Send(data []byte) error
}

type ProtocolsA struct {
	// ...
}

func (p *ProtocolsA) Send(data []byte) error {
	return nil
}

type ProtocolsB struct {
	// ...
}

func (p *ProtocolsB) Send(data []byte) error {
	return nil
}

type ProtocolsC struct {
	// ...
}

func (p *ProtocolsC) Send(data []byte) error {
	return nil
}

func Send(protocol Protocol, data []byte) {
	switch protocol := protocol.(type) {
	case nil:
		// Treat nil case
	case *ProtocolsA:
		// Do specific stuff for protocol A
	case *ProtocolsB:
		// Do specific stuff for protocol A
	case *ProtocolsC:
		// Do specific stuff for protocol A
	}
	protocol.Send(data)
}

func main() {
	var i, j int
	var ui, uj, uk uint8
	var buf = make([]byte, 100)

	for i := 0; i < 10; i++ {
		if f, err := os.Open(fmt.Sprintf("file%d.txt", i)); err != nil {
			panic(err)
		} else {
			f.Read(buf)
		}
	}
}
