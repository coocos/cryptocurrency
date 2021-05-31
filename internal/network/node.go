package network

import (
	"log"
	"os"

	"github.com/coocos/cryptocurrency/internal/blockchain"
	"github.com/coocos/cryptocurrency/internal/keys"
)

// Node represents the node running the blockchain
type Node struct {
	chain *blockchain.Blockchain
	api   *Api
	cache *BlockCache
}

// NewNode returns a new node which mines blocks using the given key pair
func NewNode(keyPair *keys.KeyPair) *Node {
	chain := blockchain.NewBlockchain(keyPair)
	cache := &BlockCache{}
	server := NewApi(cache, relayReceivedBlocks(chain))
	return &Node{
		chain,
		server,
		cache,
	}
}

// Start starts the node
func (n *Node) Start() {
	seedHost, defined := os.LookupEnv("CRYPTO_SEED_HOST")
	if defined {
		log.Println("Syncing blockchain via", seedHost)
		n.syncBlockchain(seedHost)
	}
	n.mine()
	n.api.Serve()
}

func (n *Node) syncBlockchain(seedHost string) {
	client := NodeClient{seedHost}
	blocks, err := client.GetBlocks()
	if err != nil {
		log.Fatalln("Failed to read blocks from seed host:", err)
	}
	for _, block := range blocks {
		n.chain.SubmitExternalBlock(&block)
	}
}

func relayReceivedBlocks(chain *blockchain.Blockchain) chan<- blockchain.Block {
	blocks := make(chan blockchain.Block)
	go func() {
		for block := range blocks {
			chain.SubmitExternalBlock(&block)
		}
	}()
	return blocks
}

func (n *Node) mine() {
	go func() {
		for {
			block := n.chain.MineBlock()
			n.cache.AddBlock(block)
		}
	}()
}
