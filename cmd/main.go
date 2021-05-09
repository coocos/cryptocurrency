package main

import (
	"github.com/coocos/cryptocurrency/internal/blockchain"
)

func main() {

	chain := blockchain.NewBlockchain()
	chain.Mine()

}
