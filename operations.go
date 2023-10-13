package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
)

const Question = "question"
const Answer = "answer"

const MT103_string = "MT103"
const BankDiscoveryMessage_string = "BankDiscoveryMessage"

type BankDiscoveryMessage struct {
	Type       string // to prevent infinite loop
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

func GetAccount(userIban string) (*Account, error) {
	accountState, err := GetAndDeserializeAccount(userIban)
	if err != nil {
		return nil, err
	}
	return accountState, nil
}

func GetAllAccounts() ([]*Account, error) {
	entries, err := GetAllEntries()
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

func GetMessage(key string) (*MT103Message, error) {
	// Pick message
	messageRaw, err := VerifiedGetMsg(key)
	if err != nil {
		return nil, err
	}
	message := new(MT103Message)
	err = json.Unmarshal(messageRaw.Value, message)
	return message, err
}

func GetAllMessages() ([]*MT103Message, error) {
	entries, err := GetAllMsgsEntries()
	if err != nil {
		return nil, err
	}
	messages := []*MT103Message{}
	for _, entry := range entries.Entries {
		message := new(MT103Message)
		err := json.Unmarshal(entry.Value, message)
		if err != nil {
			return messages, err
		}
		messages = append(messages, message)
	}
	return messages, nil
}

func GetAndDeserializeAccount(key string) (*Account, error) {
	// Pick state
	accountStateRaw, err := VerifiedGet(key)
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
	err = VerifiedSet(key, string(finalAccountState))
	return err
}

func FindCounterpartBanks() error {
	discoveryMsg := &BankDiscoveryMessage{Type: Question, MyBankName: InstitutionName}
	bytes, err := json.Marshal(discoveryMsg)
	if err != nil {
		return err
	}
	node.SendMessage(BankDiscoveryMessage_string, bytes)
	return nil
}

func ProcessBankDiscovery(discoveryMsg *BankDiscoveryMessage) error {
	// Pick the other bank name
	if !contains(CounterpartBanks, discoveryMsg.MyBankName) /* && discoveryMsg.MyBankName != InstitutionName */ {
		CounterpartBanks = append(CounterpartBanks, discoveryMsg.MyBankName)
	}

	// Answer if needed
	if discoveryMsg.Type == Question {
		discoveryAnswer := &BankDiscoveryMessage{Type: Answer, MyBankName: InstitutionName}
		bytes, err := json.Marshal(discoveryAnswer)
		if err != nil {
			return err
		}
		node.SendMessage(BankDiscoveryMessage_string, bytes)
	}
	return nil
}

func contains(list []string, elem string) bool {
	for _, a := range list {
		if a == elem {
			return true
		}
	}
	return false
}

func PrintBankInfo() {
	fmt.Println("| Bank Name:", InstitutionName)
	fmt.Println("| ImmuDB instance running on IP:", Client.GetOptions().Address)
	fmt.Println("| ImmuDB instance running on port:", Client.GetOptions().Port)
	fmt.Println("| ...")
}
