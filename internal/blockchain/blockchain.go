package blockchain

import (
	"encoding/base64"
	"errors"
	"log"
	"math"
	"runtime"

	"github.com/coocos/cryptocurrency/internal/keys"
)

// Blockchain represents a full blockchain
type Blockchain struct {
	blocks         []*Block
	pool           map[string]Transaction
	keyPair        *keys.KeyPair
	externalBlocks chan Block
}

// NewBlockchain returns a new blockchain with a genesis block
func NewBlockchain(keyPair *keys.KeyPair) *Blockchain {
	if keyPair == nil {
		log.Println("No key pair given - generating a new one")
		keyPair = keys.NewKeyPair()
	}
	blockchain := Blockchain{
		keyPair:        keyPair,
		pool:           make(map[string]Transaction),
		externalBlocks: make(chan Block, 128),
	}
	blockchain.addBlock(GenesisBlock())
	return &blockchain
}

// LastBlock returns the last block in the blockchain
func (b *Blockchain) LastBlock() *Block {
	if len(b.blocks) > 0 {
		return b.blocks[len(b.blocks)-1]
	}
	return nil
}

// SubmitExternalBlock sends an externally received block to the blockchain
func (b *Blockchain) SubmitExternalBlock(block *Block) {
	b.externalBlocks <- *block
}

func (b *Blockchain) addBlock(block *Block) error {
	previous := b.LastBlock()
	if previous == nil {
		log.Printf("Adding genesis block: %+v\n", block)
		b.blocks = append(b.blocks, block)
		return nil
	}
	if !block.IsValid(previous) {
		return errors.New("New block is not valid")
	}
	b.blocks = append(b.blocks, block)
	return nil
}

// AddTransaction adds transaction to the pool of available transactions to include in next block
func (b *Blockchain) AddTransaction(transaction Transaction) error {
	if !transaction.ValidSignature() {
		return errors.New("Transaction has invalid signature")
	}
	b.pool[base64.StdEncoding.EncodeToString(transaction.Signature)] = transaction
	return nil
}

func (b *Blockchain) filterValidTransactions() []Transaction {
	validTransactions := make([]Transaction, 0)
	accounts := AccountsFromBlockchain(b.blocks)
	for _, transaction := range b.pool {
		if err := accounts.ApplyTransaction(transaction); err != nil {
			log.Println("Transaction is invalid", err)
			continue
		}
		validTransactions = append(validTransactions, transaction)
		if len(validTransactions) == maxTransactionsPerBlock-1 {
			break
		}
	}
	return validTransactions
}

func (b *Blockchain) clearSpentTransactions() {
	for _, transaction := range b.LastBlock().Transactions {
		delete(b.pool, base64.StdEncoding.EncodeToString(transaction.Signature))
	}
}

func (b *Blockchain) transactionsForNextBlock() []Transaction {
	return append([]Transaction{CoinbaseTransactionTo(b.keyPair.PublicKey)}, b.filterValidTransactions()...)
}

// ProofOfWorkRequest is a request to mine a new block
type ProofOfWorkRequest struct {
	blockNumber       int
	previousBlock     Block
	blockTransactions []Transaction
}

func blockWorker(nonces <-chan int, validBlock chan<- Block, request ProofOfWorkRequest) {
	for {
		nonce, more := <-nonces
		if !more {
			return
		}
		block := NewBlock(request.blockNumber, request.previousBlock.Hash, request.blockTransactions, nonce)
		if block.IsValid(&request.previousBlock) {
			validBlock <- *block
			return
		}
	}
}

// MineBlock mines a new valid block with transactions from the mempool
func (b *Blockchain) MineBlock() Block {
	nonces := make(chan int)
	validBlock := make(chan Block, runtime.NumCPU())

	defer func() {
		close(nonces)
		b.clearSpentTransactions()
	}()

	// Create a worker per core to mine for a valid block
	for worker := 0; worker < runtime.NumCPU(); worker++ {
		go blockWorker(nonces, validBlock, ProofOfWorkRequest{b.LastBlock().Number + 1, *b.LastBlock(), b.transactionsForNextBlock()})
	}

	// Send incremental nonces to workers until a valid block is found
	for nonce := 0; nonce < math.MaxInt64; nonce++ {
		select {
		// Another node found a valid block
		case block := <-b.externalBlocks:
			if block.IsValid(b.LastBlock()) {
				if err := b.addBlock(&block); err != nil {
					log.Fatalf("Failed to add external block to blockchain: %v\n", err)
				}
				log.Println("Remote node found valid block:", block)
				return *b.LastBlock()
			}
		// Found a valid block
		case block := <-validBlock:
			if err := b.addBlock(&block); err != nil {
				log.Fatalf("Failed to add internally generated block to blockchain: %v\n", err)
			}
			log.Printf("🎉 Found valid block: %+v\n", block)
			return *b.LastBlock()
		// No valid block found yet so keep sending nonces to workers
		default:
			nonces <- nonce
		}
	}

	panic("Exhausted possible nonce values")
}
