package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
)

const MT103_string = "MT103"
const BankDiscoveryMessage_string = "BankDiscoveryMessage"
const BankDiscoveryAnswer_string = "BankDiscoveryAnswer"

type BankDiscoveryMessage struct {
	Hi string // ¿?
}

type BankDiscoveryAnswer struct {
	MyBankName string
}

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
	Suspend(accountState)
	err = SerializeAndSetAccount(userIban, accountState)
	return err
}

func UnsuspendAccount(userIban string) error {
	accountState, err := GetAndDeserializeAccount(userIban)
	if err != nil {
		return err
	}
	Unsuspend(accountState)
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
	err = SetBalance(accountState, float32(balance))
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
	err = Deposit(accountState, float32(amount))
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
	err = Withdraw(accountState, float32(amount))
	if err != nil {
		return err
	}
	err = SerializeAndSetAccount(userIban, accountState)
	return err
}

func PrintAccount(userIban string) error {
	accountState, err := GetAndDeserializeAccount(userIban)
	if err != nil {
		return err
	}
	PrintDetails(accountState)
	return nil
}

func GetAndDeserializeAccount(key string) (*Account, error) {
	// Pick state
	accountStateRaw, err := VerifiedGet(key)
	if err != nil {
		return nil, err
	}
	accountState := new(Account)
	json.Unmarshal(accountStateRaw.Value, accountState)
	return accountState, nil
}

func SerializeAndSetAccount(key string, accountState *Account) error {
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

func FindCounterpartBanks() error {
	discoveryMsg := &BankDiscoveryMessage{Hi: "hi!"}
	bytes, err := json.Marshal(discoveryMsg)
	if err != nil {
		return err
	}
	node.SendMessage(BankDiscoveryMessage_string, bytes)
	return nil
}

func AnswerBankDiscovery(discoveryMsg *BankDiscoveryMessage) error {
	discoveryAnswer := &BankDiscoveryAnswer{MyBankName: InstitutionName}
	bytes, err := json.Marshal(discoveryAnswer)
	if err != nil {
		return err
	}
	node.SendMessage(BankDiscoveryAnswer_string, bytes)
	return nil
}

func ProcessBankDiscoveryAnswer(discoveryAnswer *BankDiscoveryAnswer) error {
	CounterpartBanks = append(CounterpartBanks, discoveryAnswer.MyBankName)
	return nil
}

// TODO
func PrintBankInfo() {
	fmt.Println("Bank Name:", InstitutionName)
}
