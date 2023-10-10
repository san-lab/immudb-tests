package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
)

func CreateAccount(userIban, userName string) error {
	// Check if IBAN already in database
	_, err := VerifiedGet(userIban)
	if err == nil {
		fmt.Println("User with that IBAN is already in the database")
		return errors.New("user with that IBAN is already in the database")
	}

	accountState := AccountState{OwnerName: userName, Balance: 0, Suspended: false} // TODO add attributes
	accountStateRaw, err := json.Marshal(accountState)
	if err != nil {
		fmt.Println(err)
		return err
	}

	err = VerifiedSet(userIban, string(accountStateRaw))
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func SuspendAccount(userIban string) error {
	err := UpdateStateAttributes(userIban, []string{"Suspended"}, &AccountState{Suspended: false})
	return err
}

func SetBalance(userIban, balanceString string) error {
	balance, err := strconv.Atoi(balanceString)
	if err != nil {
		fmt.Println(err)
		return err
	}
	err = UpdateStateAttributes(userIban, []string{"Balance"}, &AccountState{Balance: balance})
	return err
}

func AddSubstractBalanceIfPossible(user, amount string, add bool) error {
	// Pick state
	accountStateRaw, err := VerifiedGet(user)
	if err != nil {
		return err
	}
	accountState := new(AccountState)
	json.Unmarshal(accountStateRaw, accountState)

	// Update state balance if possible
	amountInt, err := strconv.Atoi(string(amount))
	if err != nil {
		return err
	}
	if add {
		accountState.Balance += amountInt
	} else {
		accountState.Balance -= amountInt
		if accountState.Balance <= 0 {
			return errors.New("not enough balance to perform the transaction")
		}
	}

	// Marshal and set new state into the DB
	finalAccountState, err := json.Marshal(accountState)
	if err != nil {
		return err
	}
	err = VerifiedSet(user, string(finalAccountState))
	if err != nil {
		return err
	}
	return nil
}

func UpdateStateAttributes(key string, attributes []string, newValues *AccountState) error {
	// Pick state
	accountStateRaw, err := VerifiedGet(key)
	if err != nil {
		return err
	}
	accountState := new(AccountState)
	json.Unmarshal(accountStateRaw, accountState)

	// Update state
	if contains(attributes, "OwnerName") {
		accountState.OwnerName = newValues.OwnerName
	}
	if contains(attributes, "Balance") {
		accountState.Balance = newValues.Balance
	}
	if contains(attributes, "Suspended") {
		accountState.Suspended = newValues.Suspended
	}
	// TODO add attributes... identify a better way

	// Marshal and set new state into the DB
	finalAccountState, err := json.Marshal(accountState)
	if err != nil {
		return err
	}
	err = VerifiedSet(key, string(finalAccountState))
	if err != nil {
		return err
	}
	return nil
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
