package blockchain

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/coocos/cryptocurrency/internal/keys"
)

func TestBlockChain(t *testing.T) {
	t.Run("Test that blockchain always includes genesis block", func(t *testing.T) {
		chain := NewBlockchain(nil)

		if chain.LastBlock() == nil {
			t.Errorf("Blockchain has no genesis block\n")
		}
	})
	t.Run("Test that mined block includes coinbase transaction to miner", func(t *testing.T) {
		miner := keys.NewKeyPair()

		chain := NewBlockchain(miner)
		chain.MineBlock()

		expectedTransactions := 1
		if len(chain.LastBlock().Transactions) != expectedTransactions {
			t.Errorf("Expected %d transactions but there are %d\n", expectedTransactions, len(chain.LastBlock().Transactions))
		}
		coinbaseTransaction := chain.LastBlock().Transactions[0]
		if coinbaseTransaction.Sender != nil {
			t.Error("Coinbase transaction should have no sender")
		}
		if !bytes.Equal(coinbaseTransaction.Receiver, miner.EncodedPublicKey) {
			t.Error("Coinbase transaction not sent to miner")
		}
		if !coinbaseTransaction.ValidSignature() {
			t.Error("Coinbase transaction does not have a valid signature")
		}
		if coinbaseTransaction.Amount != 10 {
			t.Error("Coinbase transaction should be 10 coins")
		}
	})
	t.Run("Test that next mined block includes transaction", func(t *testing.T) {
		miner := keys.NewKeyPair()
		receiver := keys.NewKeyPair()

		// Mine one block so that miner has some coins
		chain := NewBlockchain(miner)
		chain.MineBlock()

		// Mine next block and send coins from miner to receiver
		transaction := NewTransaction(miner.EncodedPublicKey, receiver.EncodedPublicKey, 5)
		transaction.Sign(miner.PrivateKey)
		if err := chain.AddTransaction(*transaction); err != nil {
			t.Errorf("Failed to add transaction to blockchain: %v", err)
		}
		chain.MineBlock()

		if len(chain.LastBlock().Transactions) != 2 {
			t.Error("Transaction was not included in block")
		}
		if !reflect.DeepEqual(chain.LastBlock().Transactions[1], *transaction) {
			t.Error("Included transaction does not match submitted transaction")
		}
	})
}
