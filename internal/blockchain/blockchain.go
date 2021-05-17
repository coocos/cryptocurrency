package blockchain

import (
	"errors"
	"log"
	"os"

	"github.com/coocos/cryptocurrency/internal/keys"
)

// Blockchain represents a full blockchain
type Blockchain struct {
	chain []*Block
	pool  []*Transaction
}

// NewBlockchain returns a new blockchain with a genesis block
func NewBlockchain() *Blockchain {
	blockchain := Blockchain{}
	blockchain.AddBlock(GenesisBlock())
	return &blockchain
}

// LastBlock returns the last block in the blockchain
func (b *Blockchain) LastBlock() *Block {
	if len(b.chain) > 0 {
		return b.chain[len(b.chain)-1]
	}
	return nil
}

// AddBlock adds block to blockchain or returns an error if the block is not valid
func (b *Blockchain) AddBlock(block *Block) error {
	previous := b.LastBlock()
	if previous == nil {
		log.Printf("Adding genesis block: %+v\n", block)
		b.chain = append(b.chain, block)
		return nil
	}
	if !block.IsValid() {
		return errors.New("New block not valid according to proof-of-work")
	}
	if block.Number != previous.Number+1 || block.PreviousHash != previous.Hash {
		return errors.New("New block does not follow the last block in blockchain")
	}
	b.chain = append(b.chain, block)
	return nil
}

// AddTransaction adds transaction to the pool of available transactions to include in next block
func (b *Blockchain) AddTransaction(transaction *Transaction) error {
	if !transaction.ValidSignature() {
		return errors.New("Transaction has invalid signature")
	}
	b.pool = append(b.pool, transaction)
	return nil
}

func (b *Blockchain) transactionsForNextBlock() []*Transaction {
	// Transaction to reward the miner and increase money supply
	privateKeyPath := os.Getenv("NODE_PRIVATE_KEY")
	publicKeyPath := os.Getenv("NODE_PUBLIC_KEY")
	keyPair, err := keys.LoadKeyPair(privateKeyPath, publicKeyPath)
	if err != nil {
		log.Fatalf("Failed to load key pair, unable to sign transactions: %v\n", err)
	}
	coinbaseTransaction := &Transaction{
		Sender:   nil,
		Receiver: keyPair.EncodedPublicKey,
		Amount:   10,
	}
	coinbaseTransaction.Sign(keyPair.PrivateKey)

	return append([]*Transaction{coinbaseTransaction}, b.pool...)
}

// Mine executes the proof-of-work algorithm to mine for a valid block
func (b *Blockchain) Mine() {
	block := NewBlock(b.LastBlock().Number+1, b.LastBlock().PreviousHash, b.transactionsForNextBlock())
	for !block.IsValid() {
		block = NewBlock(b.LastBlock().Number+1, b.LastBlock().PreviousHash, b.transactionsForNextBlock())
		if !block.IsValid() {
			continue
		}
		log.Printf("Found valid block: %+v\n", block)
		b.AddBlock(block)
		break
	}
}
