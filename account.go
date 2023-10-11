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
func PrintDetails(account *Account) {
	fmt.Println("BIC:", account.Bic)
	fmt.Println("IBAN:", account.Iban)
	fmt.Println("Currency:", account.Currency)
	fmt.Println("Holder:", account.Holder)
	fmt.Println("Balance:", account.Balance)
	fmt.Println("Suspended:", account.Suspended)
}
