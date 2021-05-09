package main

import (
	"log"

	"github.com/coocos/cryptocurrency/internal/blockchain"
)

func main() {

	chain := blockchain.NewBlockchain()
	chain.Mine()
	log.Printf("%+v\n", chain.LastBlock())

}
