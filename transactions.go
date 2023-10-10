package main

import (
	"errors"
	"fmt"
)

// MT103-like message
type MT103Message struct {
	TxReferenceNumber      string
	TimeIndication         string
	BankOperationCode      string
	ValueDate              string // always today?
	Currency               string
	ExchangeRate           string
	OrderingInstitution    string
	BeneficiaryInstitution string

	OrderingCustomer    string // Sender IBAN
	BeneficiaryCustomer string // Recipient IBAN
	Amount              string
}

// Deprecated
type TxMessage struct {
	UserFrom string
	UserTo   string
	Amount   string
}

type AccountState struct {
	Balance   int
	OwnerName string
	Suspended bool

	Other      string
	Meaningful string
	Attributes string
}

func InterBankTx(userFrom, amount, userTo string) error { // TODO maybe another struct for the parameters
	AddSubstractBalanceIfPossible(userFrom, amount, false)

	// Sent event to the topic
	txmsg := &MT103Message{OrderingCustomer: userFrom, BeneficiaryCustomer: userTo, Amount: amount}
	node.SendMsg(txmsg)
	return nil
}

// When receiveing a transaction
func ProcessInterBankTx(txmsg *MT103Message) error {
	if !isIncomingTxValid(txmsg) {
		return errors.New("received transaction message is invalid")
	}

	err := AddSubstractBalanceIfPossible(txmsg.BeneficiaryCustomer, txmsg.Amount, true)
	return err
}

func IntraBankTx(userFrom, amount, userTo string) error {
	err := AddSubstractBalanceIfPossible(userFrom, amount, false)
	if err != nil {
		return err
	}

	err = AddSubstractBalanceIfPossible(userTo, amount, true)
	return err
}

func isIncomingTxValid(txmsg *MT103Message) bool {
	_, err := VerifiedGet(txmsg.BeneficiaryCustomer)
	if err != nil {
		fmt.Println("BeneficiaryCustomer is not in the database")
		return false
	}
	return true
	// TODO check more stuff..
}
