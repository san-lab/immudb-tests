package main

import (
	"fmt"
)

type Account struct {
	suspended bool
	bic       string
	iban      string
	balance   float32
	holder    string
	currency  string
}

// Opens a new bank account with the specified parameters
func SetAccount(bic string, iban string, balance float32, holder string, currency string) *Account {
	return &Account{
		suspended: false,
		bic:       bic,
		iban:      iban,
		balance:   balance,
		holder:    holder,
		currency:  currency}
}

// Suspends an account
func Suspend(account *Account) {
	account.suspended = true
}

// Deposits an amount into an account
func Deposit(account *Account, amount float32) {
	if !account.suspended {
		account.balance += amount
	} else {
		fmt.Println("Account suspended: operation not performed")
	}
}

// Withdraws an amount from an account
func Withdraw(account *Account, amount float32) {
	if !account.suspended {
		account.balance -= amount
	} else {
		fmt.Println("Account suspended: operation not performed")
	}
}

// Returns current balance of an account
func GetBalance(account *Account) float32 {
	return account.balance
}

//Print account details
func AccountDetails (account *Account){
	fmt.Println("BIC:",account.bic)
	fmt.Println("IBAN:",account.iban)
	fmt.Println("Currency:",account.currency)
	fmt.Println("Holder:",account.holder)
	fmt.Println("Balance:",account.balance)
}
