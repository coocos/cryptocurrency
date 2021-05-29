package network

import (
	"github.com/coocos/cryptocurrency/internal/blockchain"
	"github.com/coocos/cryptocurrency/internal/keys"
)

// Node represents the node running the blockchain
type Node struct {
	chain *blockchain.Blockchain
	api   *Api
}

type BlockchainSource struct {
	chain *blockchain.Blockchain
}

// ReadBlock returns the last block in the blockchain
func (b *BlockchainSource) ReadBlock() blockchain.Block {
	return *b.chain.LastBlock()
}

// SubmitBlock submits a new block to the blockchain
func (b *BlockchainSource) SubmitBlock(block blockchain.Block) {
	b.chain.SubmitExternalBlock(&block)
}

// NewNode returns a new node which mines blocks using the given key pair
func NewNode(keyPair *keys.KeyPair) *Node {
	chain := blockchain.NewBlockchain(keyPair)
	proxy := BlockchainSource{chain}
	server := NewApi(&proxy)
	return &Node{
		chain,
		server,
	}
}

// Start starts the node
func (n *Node) Start() {
	n.mine()
	n.api.Serve()
}

func (n *Node) mine() {
	go func() {
		for {
			n.chain.MineBlock()
		}
	}()
}
