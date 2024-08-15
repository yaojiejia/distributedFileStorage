package p2p

import (
	"encoding/gob"
	"io"
)

// Decoder is an interface that have Decode function
type Decoder interface {
	Decode(io.Reader, *RPC) error
}

type GOBDecoder struct{}

// Decode the message, receiving from the source r and writing it to msg
func (dec GOBDecoder) Decode(r io.Reader, msg *RPC) error {
	return gob.NewDecoder(r).Decode(msg)
}

type DefaultDecoder struct{}

// Decode is a DefaultDecoder object that decode in 1028 bytes
func (dec DefaultDecoder) Decode(r io.Reader, msg *RPC) error {
	buf := make([]byte, 1028)
	n, err := r.Read(buf)
	if err != nil {
		return err
	}

	msg.Payload = buf[:n]
	return nil
}
