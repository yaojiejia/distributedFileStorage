package p2p

// Peer is an interface that representes the remote node
type Peer interface {
}

// Transport anything that handles communication
// between the nodes in the network. This can be of the
// form (TCP, UDP, WebSockets, ...)
type Transport interface {
	ListenAndAccept() error
}
