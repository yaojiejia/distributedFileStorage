package main

import (
	"log"

	"github.com/yaojiejia/distributedfilestorage/p2p"
)

func main() {

	tcpOpts := p2p.TCPTransportOps{
		ListenAddr:    ":3000",
		HandshakeFunc: p2p.NOPHandshakefunc,
		Decoder:       p2p.DefaultDecoder{},
	}
	tr := p2p.NewTCPTransport(tcpOpts)

	if err := tr.ListenAndAccept(); err != nil {
		log.Fatal(err)
	}

	select {}
}
