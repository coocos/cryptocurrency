package network

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/coocos/cryptocurrency/internal/blockchain"
)

// NodeClient is an HTTP client used to communicate with a node
type NodeClient struct {
	peerAddress string
}

func (c *NodeClient) apiUrl(resource string) string {
	return fmt.Sprintf("http://%s/api/v1%s", c.peerAddress, resource)
}

// GetBlocks requests all blocks from peer node
func (c *NodeClient) GetBlocks() ([]blockchain.Block, error) {
	response, err := http.Get(c.apiUrl("/blockchain/"))
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	var blocks []blockchain.Block
	if err := json.NewDecoder(response.Body).Decode(&blocks); err != nil {
		return nil, err
	}
	return blocks, nil
}

// Greet sends a greeting to peer node
func (c *NodeClient) Greet() error {
	greeting := NewPeer{os.Getenv("CRYPTO_NODE_HOST")}
	payload, err := json.Marshal(greeting)
	if err != nil {
		return err
	}
	response, err := http.Post(c.apiUrl("/peer/"), "application/json", bytes.NewBuffer(payload))
	if err != nil {
		return err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("Greeting response: %v", response.StatusCode)
	}
	return nil
}
