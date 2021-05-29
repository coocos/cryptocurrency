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
		chain := NewBlockchain(minerAccount)
		chain.MineBlock()
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
