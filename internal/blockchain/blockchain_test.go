package blockchain

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/coocos/cryptocurrency/internal/keys"
)

func TestBlockChain(t *testing.T) {

	miner := keys.NewKeyPair()
	receiver := keys.NewKeyPair()

	t.Run("Test that blockchain always includes genesis block", func(t *testing.T) {
		chain := NewBlockchain(nil)

		if chain.LastBlock() == nil {
			t.Error("Blockchain has no genesis block")
		}
	})
	t.Run("Test that mined block includes coinbase transaction to miner", func(t *testing.T) {
		chain := NewBlockchain(miner)
		block := chain.MineBlock()

		expectedTransactions := 1
		if len(block.Transactions) != expectedTransactions {
			t.Errorf("Expected %d transactions but there are %d\n", expectedTransactions, len(chain.LastBlock().Transactions))
		}
		coinbaseTransaction := block.Transactions[0]
		if coinbaseTransaction.Sender != nil {
			t.Error("Coinbase transaction should have no sender")
		}
		if !bytes.Equal(coinbaseTransaction.Receiver, miner.PublicKey) {
			t.Error("Coinbase transaction not sent to miner")
		}
		if !coinbaseTransaction.ValidSignature() {
			t.Error("Coinbase transaction does not have a valid signature")
		}
		if coinbaseTransaction.Amount != 10 {
			t.Error("Coinbase transaction should be 10 coins")
		}
	})
	t.Run("Test that mined block includes transaction", func(t *testing.T) {
		// Mine one block so that miner has some coins
		chain := NewBlockchain(miner)
		chain.MineBlock()

		// Mine next block to send coins from miner to receiver
		transaction := NewTransaction(miner.PublicKey, receiver.PublicKey, 5)
		transaction.Sign(miner.PrivateKey)
		if err := chain.AddTransaction(*transaction); err != nil {
			t.Errorf("Failed to add transaction to blockchain: %v", err)
		}
		block := chain.MineBlock()
		if len(block.Transactions) != 2 {
			t.Error("Transaction was not included in block")
		}
		if !reflect.DeepEqual(block.Transactions[1], *transaction) {
			t.Error("Included transaction does not match submitted transaction")
		}
	})
	t.Run("Test that mined block does not include overspent transaction", func(t *testing.T) {
		// Mine one block so that miner has some coins
		chain := NewBlockchain(miner)
		chain.MineBlock()

		// Mine next block to send coins from miner to receiver
		transaction := NewTransaction(miner.PublicKey, receiver.PublicKey, 15)
		transaction.Sign(miner.PrivateKey)
		if err := chain.AddTransaction(*transaction); err != nil {
			t.Errorf("Failed to add transaction to blockchain: %v", err)
		}
		block := chain.MineBlock()
		if len(block.Transactions) > 1 {
			t.Fatal("Invalid transaction was included in block")
		}
	})
	t.Run("Test that spent transaction is not included in the next block", func(t *testing.T) {
		// Mine one block so that miner has some coins
		chain := NewBlockchain(miner)
		chain.MineBlock()

		// Mine next block to send coins from miner to receiver
		transaction := NewTransaction(miner.PublicKey, receiver.PublicKey, 5)
		transaction.Sign(miner.PrivateKey)
		if err := chain.AddTransaction(*transaction); err != nil {
			t.Errorf("Failed to add transaction to blockchain: %v", err)
		}
		chain.MineBlock()

		// Mine another block to see that it does include the spent transaction
		block := chain.MineBlock()
		if len(block.Transactions) > 1 {
			t.Error("Block included an already spent transaction")
		}
	})
	t.Run("Test that blockchain accepts valid blocks from other chains", func(t *testing.T) {
		firstChain := NewBlockchain(miner)
		secondChain := NewBlockchain(miner)

		firstBlock := firstChain.MineBlock()
		secondChain.SubmitExternalBlock(firstChain.LastBlock())
		secondBlock := secondChain.MineBlock()

		if !reflect.DeepEqual(firstBlock, secondBlock) {
			t.Error("Blockchain did not accept block from other chain")
		}
	})
}
