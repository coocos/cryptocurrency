package blockchain

import (
	"testing"

	"github.com/coocos/cryptocurrency/internal/keys"
)

func TestAccount(t *testing.T) {

	t.Run("Test applying coinbase transaction", func(t *testing.T) {
		keys := keys.NewKeyPair()
		accounts := NewAccounts()

		if err := accounts.ApplyTransaction(CoinbaseTransactionTo(keys.PublicKey)); err != nil {
			t.Error("Failed to apply coinbase transaction:", err)
		}

		account, err := accounts.Read(keys.PublicKey)
		if err != nil {
			t.Error("Failed to read account balance:", err)
		}
		if account.Balance != 10 {
			t.Errorf("Expected account balance %v, real account balance %v\n", 10, account.Balance)
		}
	})
	t.Run("Test applying regular transaction", func(t *testing.T) {
		sender := keys.NewKeyPair()
		receiver := keys.NewKeyPair()
		accounts := NewAccounts()

		accounts.ApplyTransaction(CoinbaseTransactionTo(sender.PublicKey))
		transaction := NewTransaction(sender.PublicKey, receiver.PublicKey, 5, 1)
		transaction.Sign(sender.PrivateKey)
		if err := accounts.ApplyTransaction(*transaction); err != nil {
			t.Error("Failed to apply transaction", err)
		}

		senderAccount, _ := accounts.Read(sender.PublicKey)
		if senderAccount.Balance != 5 {
			t.Errorf("Expected account balance %v, real account balance %v\n", 5, senderAccount.Balance)
		}
		if senderAccount.Nonce != 1 {
			t.Error("Sender account nonce not incremented")
		}
		receiverAccount, _ := accounts.Read(receiver.PublicKey)
		if receiverAccount.Nonce != 0 {
			t.Error("Receiver account nonce should be zero")
		}
	})
	t.Run("Test rejecting replayed transactions", func(t *testing.T) {
		sender := keys.NewKeyPair()
		receiver := keys.NewKeyPair()
		accounts := NewAccounts()

		accounts.ApplyTransaction(CoinbaseTransactionTo(sender.PublicKey))
		transaction := NewTransaction(sender.PublicKey, receiver.PublicKey, 5, 1)
		transaction.Sign(sender.PrivateKey)
		accounts.ApplyTransaction(*transaction)

		if err := accounts.ApplyTransaction(*transaction); err == nil {
			t.Error("Applied same transaction a second time", err)
		}
	})
	t.Run("Test rejecting an overspent transaction", func(t *testing.T) {
		sender := keys.NewKeyPair()
		receiver := keys.NewKeyPair()
		accounts := NewAccounts()

		transaction := NewTransaction(sender.PublicKey, receiver.PublicKey, 10, 1)
		transaction.Sign(sender.PrivateKey)
		if err := accounts.ApplyTransaction(*transaction); err == nil {
			t.Error("Applied invalid transaction")
		}
	})
	t.Run("Test reading balances from blockchain", func(t *testing.T) {
		miner := keys.NewKeyPair()
		chain := NewBlockchain(miner)

		chain.MineBlock()
		accounts := ReadAccounts(chain.blocks)
		account, err := accounts.Read(miner.PublicKey)
		if err != nil {
			t.Error("Failed to read accounts from blockchain:", err)
		}
		if account.Balance != 10 {
			t.Errorf("Account balance based on blockchain should be %v but is %v\n", 10, account.Balance)
		}
	})
}
