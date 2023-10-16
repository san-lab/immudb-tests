package main

import (
	"errors"
	"fmt"
)

type Account struct {
	Suspended bool
	Bic       string
	Iban      string
	Balance   float32
	Holder    string
	Currency  string
}

// Opens a new bank account with the specified parameters
func SetAccount(bic string, iban string, balance float32, holder string, currency string) *Account {
	return &Account{
		Suspended: false,
		Bic:       bic,
		Iban:      iban,
		Balance:   balance,
		Holder:    holder,
		Currency:  currency}
}

// Suspends an account
func Suspend(account *Account) {
	account.Suspended = true
}

func Unsuspend(account *Account) {
	account.Suspended = false
}

// Deposits an amount into an account
func Deposit(account *Account, amount float32) error {
	if !account.Suspended {
		account.Balance += amount
	} else {
		return errors.New("account suspended: operation not performed")
	}
	return nil
}

// Withdraws an amount from an account
func Withdraw(account *Account, amount float32) error {
	if !account.Suspended {
		if amount <= account.Balance {
			account.Balance -= amount
		} else {
			return errors.New("balance cannot be negative")
		}

	} else {
		return errors.New("account suspended: operation not performed")
	}
	return nil
}

// Returns current balance of an account
func GetBalance(account *Account) float32 {
	return account.Balance
}

func SetBalance(account *Account, newBalance float32) error {
	if newBalance < 0 {
		return errors.New("balance cannot be negative")
	}
	account.Balance = newBalance
	return nil
}

// Print account details
func PrintAccount(account *Account, spacing bool) {
	if spacing {
		fmt.Println(" -----------------")
		fmt.Printf("| IBAN: %s\n| Holder: %s\n| Balance: %.2f\n| Currency: %s\n| BIC: %s\n| Suspended: %t\n",
			account.Iban, account.Holder, account.Balance, account.Currency, account.Bic, account.Suspended)
		fmt.Println(" -----------------")
	} else {
		fmt.Printf("| IBAN: %s | Holder: %s | Balance: %.2f | Currency: %s | BIC: %s | Suspended: %t\n",
			account.Iban, account.Holder, account.Balance, account.Currency, account.Bic, account.Suspended)
	}
}

func PrintAllAccounts(accounts []*Account) {
	fmt.Println(" -----------------")
	for _, account := range accounts {
		PrintAccount(account, false)
	}
	fmt.Println(" -----------------")
}
