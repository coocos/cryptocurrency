package main

import (
	"flag"
	"log"

	"github.com/coocos/cryptocurrency/internal/keys"
	"github.com/coocos/cryptocurrency/internal/network"
)

// Options passed as CLI flags
type Options struct {
	privateKey string
	publicKey  string
}

func parseArgs() Options {
	options := Options{}
	flag.StringVar(&options.privateKey, "private", "private.key", "private key path")
	flag.StringVar(&options.publicKey, "public", "public.key", "public key path")
	flag.Parse()
	return options
}

func loadMinerKeyPair(options Options) *keys.KeyPair {
	keyPair, err := keys.LoadKeyPair(options.privateKey)
	if err != nil {
		log.Fatalf("Failed to load key pair, unable to sign transactions: %v\n", err)

	}
	return keyPair
}

func main() {
	keyOptions := parseArgs()
	node := network.NewNode(loadMinerKeyPair(keyOptions))
	node.Start()
}
