package main

import (
	"encoding/json"
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

func InterBankTx(userFrom, amount, userTo, bankTo string) error { // TODO maybe another struct for the parameters
	err := WithdrawFromAccount(userFrom, amount)
	if err != nil {
		return err
	}

	// Sent event to the topic
	txMsg := &MT103Message{OrderingInstitution: InstitutionName, OrderingCustomer: userFrom, BeneficiaryInstitution: bankTo, BeneficiaryCustomer: userTo, Amount: amount}
	bytes, err := json.Marshal(txMsg)
	if err != nil {
		return err
	}
	node.SendMessage(MT103_string, bytes)
	return nil
}

// When receiveing a transaction
func ProcessInterBankTx(txmsg *MT103Message) error {
	if !validAndAddressedToUs(txmsg) {
		return errors.New("received transaction message is invalid")
	}

	err := DepositToAccount(txmsg.BeneficiaryCustomer, txmsg.Amount)
	return err
}

func IntraBankTx(userFrom, amount, userTo string) error {
	err := WithdrawFromAccount(userFrom, amount)
	if err != nil {
		return err
	}

	err = DepositToAccount(userTo, amount)
	return err
}

func validAndAddressedToUs(txmsg *MT103Message) bool {
	if txmsg.BeneficiaryInstitution != InstitutionName {
		return false
	}

	_, err := VerifiedGet(txmsg.BeneficiaryCustomer)
	if err != nil {
		fmt.Println("Beneficiary customer is not in the database")
		return false
	}
	return true
	// TODO check more stuff..
}
