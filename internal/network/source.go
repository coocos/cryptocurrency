package network

import (
	"github.com/coocos/cryptocurrency/internal/blockchain"
)

// BlockSource is used to read and submit blocks to the blockchain
type BlockSource interface {
	ReadBlock() blockchain.Block
	SubmitBlock(block blockchain.Block)
}
