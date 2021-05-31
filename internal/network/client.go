package network

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/coocos/cryptocurrency/internal/blockchain"
)

// NodeClient is an HTTP client used to communicate with a node
type NodeClient struct {
	address string
}

// GetBlocks requests all blocks from the remote node
func (n *NodeClient) GetBlocks() ([]blockchain.Block, error) {
	response, err := http.Get(fmt.Sprintf("http://%s/api/v1/blockchain/", n.address))
	if err != nil {
		return nil, err
	}
	var blocks []blockchain.Block
	if err := json.NewDecoder(response.Body).Decode(&blocks); err != nil {
		return nil, err
	}
	return blocks, nil
}
