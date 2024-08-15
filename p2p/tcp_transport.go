package p2p

import (
	"fmt"
	"net"
)

// TCPPeer represents the remote node over a tcp established connection.
type TCPPeer struct {
	//conn is the underlying connection of the peer
	conn net.Conn
	//dial to a remote node => outbound => true
	// accept from a dial => outbound => false
	outbound bool
}

// NewTCPPeer create a NewTCPPeer object with connection and outbound field
func NewTCPPeer(conn net.Conn, outbound bool) *TCPPeer {
	return &TCPPeer{
		conn:     conn,
		outbound: outbound,
	}
}

// close implements the peer interface
func (p *TCPPeer) Close() error {
	return p.conn.Close()
}

// TCPTransportOps is a object that can edit the options for TCPTransport
type TCPTransportOps struct {
	ListenAddr    string
	HandshakeFunc HandshakeFunc
	Decoder       Decoder
	OnPeer        func(Peer) error
}

// TCPTransport take in the TCPTransportOps, listener, and the RPC channel for data transmitting
type TCPTransport struct {
	TCPTransportOps
	listener net.Listener
	rpcch    chan RPC
}

// NewTCPTransport create a newtcptransport with default options
func NewTCPTransport(opts TCPTransportOps) *TCPTransport {
	return &TCPTransport{
		TCPTransportOps: opts,
		rpcch:           make(chan RPC),
	}
}

// Consume implements the transport interface, which will return read only channel
// for reading the incoming messages received from another peer in the network
func (t *TCPTransport) Consume() <-chan RPC {
	return t.rpcch
}

// ListenAndAccept listens from the tcp and accept the loop
func (t *TCPTransport) ListenAndAccept() error {

	var err error

	t.listener, err = net.Listen("tcp", t.ListenAddr)
	if err != nil {
		return err
	}

	go t.startAcceptLoop()

	return nil

}

// startAcceptLoop accept the listening and call the handleConn
func (t *TCPTransport) startAcceptLoop() {
	for {
		conn, err := t.listener.Accept()
		if err != nil {
			fmt.Printf("TCP accept error: %s\n", err)
		}

		fmt.Printf("New incoming tcp connetion: %+v\n", conn)
		go t.handleConn(conn)

	}

}

// handleConn handles the connection, create a newTCPPeer, shakehands, and read loop
// from RPC
func (t *TCPTransport) handleConn(conn net.Conn) {

	var err error

	defer func() {
		fmt.Printf("dropping peer connection: %s", err)
		conn.Close()
	}()

	peer := NewTCPPeer(conn, true)

	if err = t.HandshakeFunc(peer); err != nil {
		return
	}

	if t.OnPeer != nil {
		if err = t.OnPeer(peer); err != nil {
			return
		}
	}
	//Read loop
	rpc := RPC{}
	for {
		err := t.Decoder.Decode(conn, &rpc)
		if err != nil {
			return
		}

		rpc.From = conn.RemoteAddr()
		t.rpcch <- rpc

	}

}
