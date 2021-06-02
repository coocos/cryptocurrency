package network

import "github.com/coocos/cryptocurrency/internal/blockchain"

// NewBlockEvent indicates a peer has mined a new block
type NewBlock struct {
	Block blockchain.Block `json:"block"`
}

// NewTransaction indicates a peer has received a new transaction
type NewTransaction struct {
	Transaction blockchain.Transaction `json:"transaction"`
}

// NewPeer indicates a new peer has been discovered
type NewPeer struct {
	Address string `json:"peerAddress"`
}
