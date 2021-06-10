package blockchain

import (
	"crypto/ed25519"
	"encoding/base64"
	"errors"
	"fmt"
)

// Account holds coins and an incrementing nonce
type Account struct {
	Address ed25519.PublicKey `json:"address"`
	Nonce   uint              `json:"nonce"`
	Balance uint              `json:"balance"`
}

// Accounts represents all the accounts within the blockchain
type Accounts struct {
	accounts map[string]*Account
}

// NewAccounts returns an empty Accounts struct
func NewAccounts() *Accounts {
	return &Accounts{make(map[string]*Account)}
}

// Read returns the account matching the address or an error if the account is unknown
func (a *Accounts) Read(address ed25519.PublicKey) (*Account, error) {
	accountId := base64.StdEncoding.EncodeToString(address)
	account, exists := a.accounts[accountId]
	if !exists {
		return nil, fmt.Errorf("Account %s does not exist", accountId)
	}
	return account, nil
}

// ListAccounts returns all known accounts
func (a *Accounts) ListAccounts() []Account {
	accounts := make([]Account, len(a.accounts))
	for _, account := range a.accounts {
		accounts = append(accounts, *account)
	}
	return accounts
}

// ApplyTransaction applies the transaction if it's valid
func (a *Accounts) ApplyTransaction(transaction Transaction) error {
	if !transaction.ValidSignature() {
		return errors.New("Invalid transaction signature")
	}
	if !transaction.IsCoinbase() {
		if err := a.subtract(transaction.Sender, transaction.Amount, transaction.Nonce); err != nil {
			return err
		}
	}
	a.add(transaction.Receiver, transaction.Amount)
	return nil
}

// AccountsFromBlockchain generates the current account states from the blockchain
func AccountsFromBlockchain(blocks []*Block) *Accounts {
	accounts := NewAccounts()
	for _, block := range blocks {
		for _, transaction := range block.Transactions {
			accounts.ApplyTransaction(transaction)
		}
	}
	return accounts
}

func (a *Accounts) add(address ed25519.PublicKey, amount uint) {
	accountId := base64.StdEncoding.EncodeToString(address)
	account, exists := a.accounts[accountId]
	if !exists {
		a.accounts[accountId] = &Account{
			Address: address,
			Balance: amount,
			Nonce:   0,
		}
		return
	}
	account.Balance += amount
}

func (a *Accounts) subtract(address ed25519.PublicKey, amount uint, nonce uint) error {
	accountId := base64.StdEncoding.EncodeToString(address)
	account, exists := a.accounts[accountId]
	if !exists {
		return fmt.Errorf("Account %v not found", accountId)
	}
	if nonce != account.Nonce+1 {
		return errors.New("Transaction has invalid nonce")
	}
	if amount > account.Balance {
		return fmt.Errorf("Account %v has insufficient balance", accountId)
	}
	account.Balance -= amount
	account.Nonce += 1
	return nil
}
