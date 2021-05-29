package blockchain

import (
	"testing"

	"github.com/coocos/cryptocurrency/internal/keys"
)

func TestAccount(t *testing.T) {
	t.Run("Test adding to account", func(t *testing.T) {
		account := keys.NewKeyPair().PublicKey
		accounts := NewAccounts()
		accounts.Add(account, 10)

		balance, err := accounts.BalanceFor(account)
		if err != nil {
			t.Errorf("Failed to read account balance: %v\n", err)
		}
		if balance != 10 {
			t.Errorf("Expected account balance %v, real account balance %v\n", 10, balance)
		}
	})
	t.Run("Test subtracting from account", func(t *testing.T) {
		account := keys.NewKeyPair().PublicKey
		accounts := NewAccounts()
		accounts.Add(account, 10)
		err := accounts.Subtract(account, 5)
		if err != nil {
			t.Errorf("Failed to subtract from account: %v\n", err)
		}
		balance, err := accounts.BalanceFor(account)
		if err != nil {
			t.Errorf("Failed to read account balance: %v\n", err)
		}
		if balance != 5 {
			t.Errorf("Expected account balance %v, real account balance %v\n", 5, balance)
		}
	})
	t.Run("Test subtracting more than available", func(t *testing.T) {
		account := keys.NewKeyPair().PublicKey
		accounts := NewAccounts()
		accounts.Add(account, 10)
		err := accounts.Subtract(account, 15)
		if err == nil {
			t.Error("Subtracted more from account than available")
		}
	})
	t.Run("Test reading balances from blockchain", func(t *testing.T) {
		minerAccount := keys.NewKeyPair()
		coinbaseTransaction := NewTransaction(nil, minerAccount.PublicKey, 10)
		coinbaseTransaction.Sign(minerAccount.PrivateKey)
		chain := NewBlockchain(nil)
		nonce := 0
		block := NewBlock(chain.LastBlock().Number+1, chain.LastBlock().Hash, []Transaction{*coinbaseTransaction}, nonce)
		// FIXME: Instead of mining a valid block, use deterministic key generation and hardcode the block hash
		for !block.IsValid(chain.LastBlock()) {
			nonce += 1
			block = NewBlock(chain.LastBlock().Number+1, chain.LastBlock().Hash, []Transaction{*coinbaseTransaction}, nonce)
		}
		err := chain.addBlock(block)
		if err != nil {
			t.Errorf("Failed to add block to blockchain: %v\n", err)
		}

		accounts := ReadAccounts(chain)
		balance, err := accounts.BalanceFor(minerAccount.PublicKey)
		if err != nil {
			t.Errorf("Failed to read accounts from blockchain: %v\n", err)
		}
		if balance != 10 {
			t.Errorf("Account balance based on blockchain should be %v but is %v\n", 10, balance)
		}
	})
}
