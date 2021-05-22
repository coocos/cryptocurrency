package blockchain

import (
	"bytes"
	"encoding/base64"
	"errors"
	"log"
	"math"
	"runtime"

	"github.com/coocos/cryptocurrency/internal/keys"
)

// Blockchain represents a full blockchain
type Blockchain struct {
	chain        []*Block
	pool         map[string]Transaction
	keyPair      *keys.KeyPair
	eventEmitter EventEmitter
}

// NewBlockchain returns a new blockchain with a genesis block
func NewBlockchain(keyPair *keys.KeyPair, eventEmitter EventEmitter) *Blockchain {
	if keyPair == nil {
		log.Println("No key pair given - generating a new one")
		keyPair = keys.NewKeyPair()
	}
	if eventEmitter == nil {
		eventEmitter = &DummyEventEmitter{}
	}
	blockchain := Blockchain{
		keyPair:      keyPair,
		pool:         make(map[string]Transaction),
		eventEmitter: eventEmitter,
	}
	blockchain.addBlock(GenesisBlock())
	return &blockchain
}

// LastBlock returns the last block in the blockchain
func (b *Blockchain) LastBlock() *Block {
	if len(b.chain) > 0 {
		return b.chain[len(b.chain)-1]
	}
	return nil
}

func (b *Blockchain) addBlock(block *Block) error {
	previous := b.LastBlock()
	if previous == nil {
		log.Printf("Adding genesis block: %+v\n", block)
		b.chain = append(b.chain, block)
		return nil
	}
	if !block.IsValid() {
		return errors.New("New block is not valid")
	}
	if block.Number != previous.Number+1 || !bytes.Equal(block.PreviousHash, previous.Hash) {
		return errors.New("New block does not follow the last block in blockchain")
	}
	b.chain = append(b.chain, block)
	b.eventEmitter.EmitBlock(*block)
	return nil
}

// AddTransaction adds transaction to the pool of available transactions to include in next block
func (b *Blockchain) AddTransaction(transaction Transaction) error {
	if !transaction.ValidSignature() {
		return errors.New("Transaction has invalid signature")
	}
	b.pool[base64.StdEncoding.EncodeToString(transaction.Signature)] = transaction
	b.eventEmitter.EmitTransaction(transaction)
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

func (b *Blockchain) clearSpentTransactions() {
	for _, transaction := range b.LastBlock().Transactions {
		delete(b.pool, base64.StdEncoding.EncodeToString(transaction.Signature))
	}
}

func (b *Blockchain) transactionsForNextBlock() []Transaction {
	coinbaseTransaction := Transaction{
		Sender:   nil,
		Receiver: b.keyPair.PublicKey,
		Amount:   10,
	}
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
			b.addBlock(&block)
			log.Printf("🎉 Found valid block: %+v\n", block)
			for _, transaction := range block.Transactions {
				log.Printf("💰 %v\n", transaction)
			}
			close(nonces)
			b.clearSpentTransactions()
			return
		case nonces <- nonce:
		}
	}
}
