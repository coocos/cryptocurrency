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
	peers map[string]bool
}

// NewNode returns a new node which mines blocks using the given key pair
func NewNode(keyPair *keys.KeyPair) *Node {
	chain := blockchain.NewBlockchain(keyPair)
	api := NewApi(eventBus(chain))
	return &Node{
		chain: chain,
		api:   api,
		peers: make(map[string]bool),
	}
}

// Start starts the node
func (n *Node) Start() {
	if seedHost, ok := os.LookupEnv("CRYPTO_SEED_HOST"); ok {
		log.Println("Syncing blockchain via", seedHost)
		n.greetPeer(seedHost)
		n.syncBlockchain(seedHost)
	}
	n.mine()
	n.api.Serve()
}

func (n *Node) greetPeer(peerAddress string) {
	client := NodeClient{peerAddress}
	if err := client.Greet(); err != nil {
		log.Printf("Failed to greet %s: %v\n", client, err)
		return
	}
	n.peers[peerAddress] = true
}

func (n *Node) syncBlockchain(peer string) {
	client := NodeClient{peer}
	blocks, err := client.GetBlocks()
	if err != nil {
		log.Fatalln("Failed to read blocks from peer:", err)
	}
	for _, block := range blocks {
		n.chain.SubmitExternalBlock(&block)
	}
}

func eventBus(chain *blockchain.Blockchain) chan<- interface{} {
	events := make(chan interface{})
	go func() {
		for event := range events {
			switch e := event.(type) {
			case NewBlock:
				chain.SubmitExternalBlock(&e.Block)
			case NewPeer:
				log.Println("Notified of new peer", e.Address)
			default:
				log.Fatalf("Received an unknown event: %v\n", event)

			}
		}
	}()
	return events
}

func (n *Node) mine() {
	go func() {
		for {
			block := n.chain.MineBlock()
			n.api.updateCache(block)
		}
	}()
}
