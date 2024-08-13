package p2p

import "errors"

// ErrInvalidHandshake is returned if the handshake between
// local and remote node could nto be established
var ErrInvalidHandshake = errors.New("invalid handshake")

// HandshakeFunc
type HandshakeFunc func(Peer) error

func NOPHandshakefunc(Peer) error { return nil }
