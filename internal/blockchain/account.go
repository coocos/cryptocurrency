package blockchain

import (
	"encoding/base64"
	"errors"
	"fmt"
)

// Accounts represents all the accounts within the blockchain
type Accounts struct {
	accounts map[string]uint
}

// NewAccounts returns an empty Accounts struct
func NewAccounts() *Accounts {
	return &Accounts{make(map[string]uint)}
}

// BalanceFor returns the amount of coins held in account
func (a *Accounts) BalanceFor(account []byte) (uint, error) {
	accountId := base64.StdEncoding.EncodeToString(account)
	balance, exists := a.accounts[accountId]
	if !exists {
		return 0, errors.New(fmt.Sprintf("Account %v not found", accountId))
	}
	return balance, nil
}

// Add adds a number of coins to the given account
func (a *Accounts) Add(account []byte, amount uint) {
	accountId := base64.StdEncoding.EncodeToString(account)
	_, exists := a.accounts[accountId]
	if !exists {
		a.accounts[accountId] = 0
	}
	a.accounts[accountId] += amount
}

// Subtract removes a number of coins from the given account
func (a *Accounts) Subtract(account []byte, amount uint) error {
	accountId := base64.StdEncoding.EncodeToString(account)
	balance, exists := a.accounts[accountId]
	if !exists {
		return errors.New(fmt.Sprintf("Account %v not found", accountId))
	}
	if amount > balance {
		return errors.New(fmt.Sprintf("Account %v has insufficient balance", accountId))
	}
	a.accounts[accountId] -= amount
	return nil
}

// ReadAccounts goes through the blockchain and returns all the accounts
func ReadAccounts(blockchain *Blockchain) *Accounts {
	accounts := NewAccounts()
	for _, block := range blockchain.blocks {
		for _, transaction := range block.Transactions {
			accounts.Add(transaction.Receiver, transaction.Amount)
			// Coinbase transactions do not have a sender
			if transaction.Sender != nil {
				accounts.Subtract(transaction.Sender, transaction.Amount)
			}
		}
	}
	return accounts
}
