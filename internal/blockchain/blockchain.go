package blockchain

import (
	"bytes"
	"errors"
	"log"
	"math"
	"os"
	"runtime"

	"github.com/coocos/cryptocurrency/internal/keys"
)

// Blockchain represents a full blockchain
type Blockchain struct {
	chain []*Block
	pool  []Transaction
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
	if block.Number != previous.Number+1 || !bytes.Equal(block.PreviousHash, previous.Hash) {
		return errors.New("New block does not follow the last block in blockchain")
	}
	b.chain = append(b.chain, block)
	return nil
}

// AddTransaction adds transaction to the pool of available transactions to include in next block
func (b *Blockchain) AddTransaction(transaction Transaction) error {
	if !transaction.ValidSignature() {
		return errors.New("Transaction has invalid signature")
	}
	b.pool = append(b.pool, transaction)
	return nil
}

func (b *Blockchain) filterValidTransactions() []Transaction {
	validTransactions := make([]Transaction, 0)
	accounts := ReadAccounts(b)
	for _, transaction := range b.pool {
		if err := accounts.Subtract(transaction.Sender, transaction.Amount); err != nil {
			log.Printf("Rejecting transaction: %v\n", transaction)
			continue
		}
		accounts.Add(transaction.Receiver, transaction.Amount)
		validTransactions = append(validTransactions, transaction)
	}
	return validTransactions
}

func (b *Blockchain) transactionsForNextBlock() []Transaction {
	// Transaction to reward the miner and increase money supply
	privateKeyPath := os.Getenv("NODE_PRIVATE_KEY")
	publicKeyPath := os.Getenv("NODE_PUBLIC_KEY")
	keyPair, err := keys.LoadKeyPair(privateKeyPath, publicKeyPath)
	if err != nil {
		log.Fatalf("Failed to load key pair, unable to sign transactions: %v\n", err)
	}
	coinbaseTransaction := Transaction{
		Sender:   nil,
		Receiver: keyPair.EncodedPublicKey,
		Amount:   10,
	}
	coinbaseTransaction.Sign(keyPair.PrivateKey)
	return append([]Transaction{coinbaseTransaction}, b.filterValidTransactions()...)
}

// ProofOfWorkRequest is a request to mine a new block
type ProofOfWorkRequest struct {
	blockNumber       int
	previousBlockHash []byte
	blockTransactions []Transaction
}

// MineBlock mines a new valid block with transactions from the mempool
func (b *Blockchain) MineBlock() {
	nonces := make(chan int)
	validBlock := make(chan Block, runtime.NumCPU())

	// Create a worker per core to mine for a valid block
	for worker := 0; worker < runtime.NumCPU(); worker++ {
		go func(nonces <-chan int, validBlock chan<- Block, request ProofOfWorkRequest) {
			for {
				nonce, more := <-nonces
				if !more {
					return
				}
				block := NewBlock(request.blockNumber, request.previousBlockHash, request.blockTransactions, nonce)
				if block.IsValid() {
					validBlock <- *block
					return
				}
			}
		}(nonces, validBlock, ProofOfWorkRequest{b.LastBlock().Number + 1, b.LastBlock().Hash, b.transactionsForNextBlock()})
	}

	// Send incremental nonces to workers until a valid block is found
	for nonce := 0; nonce < math.MaxInt64; nonce++ {
		select {
		case block := <-validBlock:
			b.AddBlock(&block)
			log.Printf("Found valid block: %+v\n", block)
			for _, transaction := range block.Transactions {
				log.Println(transaction)
			}
			close(nonces)
			return
		case nonces <- nonce:
		}
	}
}
