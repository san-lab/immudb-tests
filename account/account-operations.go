package account

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	sdk "github.com/san-lab/immudb-tests/immudbsdk"
)

func CreateAccount(userIban, userName string) error {
	// Check if IBAN already in database
	_, err := GetAndDeserializeAccount(userIban)
	if err == nil {
		fmt.Println("User with that IBAN is already in the database")
		return errors.New("user with that IBAN is already in the database")
	}

	accountState := SetAccount("", userIban, 0, userName, "")

	err = SerializeAndSetAccount(userIban, accountState)
	return err
}

func SuspendAccount(userIban string) error {
	accountState, err := GetAndDeserializeAccount(userIban)
	if err != nil {
		return err
	}
	accountState.Suspend()
	err = SerializeAndSetAccount(userIban, accountState)
	return err
}

func UnsuspendAccount(userIban string) error {
	accountState, err := GetAndDeserializeAccount(userIban)
	if err != nil {
		return err
	}
	accountState.Unsuspend()
	err = SerializeAndSetAccount(userIban, accountState)
	return err
}

func SetAccountBalance(userIban, balanceString string) error {
	balance, err := strconv.ParseFloat(balanceString, 32)
	if err != nil {
		fmt.Println(err)
		return err
	}

	accountState, err := GetAndDeserializeAccount(userIban)
	if err != nil {
		return err
	}
	err = accountState.SetBalance(float32(balance))
	if err != nil {
		return err
	}
	err = SerializeAndSetAccount(userIban, accountState)
	return err
}

func DepositToAccount(userIban, amountString string) error {
	amount, err := strconv.ParseFloat(amountString, 32)
	if err != nil {
		fmt.Println(err)
		return err
	}

	accountState, err := GetAndDeserializeAccount(userIban)
	if err != nil {
		return err
	}
	err = accountState.Deposit(float32(amount))
	if err != nil {
		return err
	}
	err = SerializeAndSetAccount(userIban, accountState)
	return err
}

func WithdrawFromAccount(userIban, amountString string) error {
	amount, err := strconv.ParseFloat(amountString, 32)
	if err != nil {
		fmt.Println(err)
		return err
	}

	accountState, err := GetAndDeserializeAccount(userIban)
	if err != nil {
		return err
	}
	err = accountState.Withdraw(float32(amount))
	if err != nil {
		return err
	}
	err = SerializeAndSetAccount(userIban, accountState)
	return err
}

func GetAccount(userIban string) (*Account, error) {
	accountState, err := GetAndDeserializeAccount(userIban)
	if err != nil {
		return nil, err
	}
	return accountState, nil
}

func GetAccountDigest(userIban string) ([]byte, error) {
	accountState, err := GetAndDeserializeAccount(userIban)
	if err != nil {
		return nil, err
	}
	accountDigest, err := accountState.GetDigest()
	return accountDigest, nil
}

func GetAndDeserializeAccount(key string) (*Account, error) {
	// Pick state
	accountStateRaw, err := sdk.VerifiedGet(key)
	if err != nil {
		return nil, err
	}
	accountState := new(Account)
	err = json.Unmarshal(accountStateRaw.Value, accountState)
	return accountState, err
}

func SerializeAndSetAccount(key string, accountState *Account) error {
	// Marshal and set new state into the DB
	finalAccountState, err := json.Marshal(accountState)
	if err != nil {
		return err
	}
	err = sdk.VerifiedSet(key, string(finalAccountState))
	return err
}

func GetAllAccounts() ([]*Account, error) {
	entries, err := sdk.GetAllEntries()
	if err != nil {
		return nil, err
	}
	accounts := []*Account{}
	for _, entry := range entries.Entries {
		accountState := new(Account)
		err := json.Unmarshal(entry.Value, accountState)
		if err != nil {
			return accounts, err
		}
		accounts = append(accounts, accountState)
	}
	return accounts, nil
}
