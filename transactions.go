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

	// Send event to the topic and store it in MsgsDB
	txmsg := &MT103Message{OrderingInstitution: InstitutionName, OrderingCustomer: userFrom, BeneficiaryInstitution: bankTo, BeneficiaryCustomer: userTo, Amount: amount}
	bytes, err := json.Marshal(txmsg)
	if err != nil {
		return err
	}
	hash, err := StoreInMsgsDB(txmsg)
	if err != nil {
		return err
	}
	fmt.Println("Hash of the message sent:", hash)
	node.SendMessage(MT103_string, bytes)
	return nil
}

// When receiveing a transaction
func ProcessInterBankTx(txmsg *MT103Message) error {
	if !validAndAddressedToUs(txmsg) {
		return errors.New("received transaction message is invalid")
	}
	hash, err := StoreInMsgsDB(txmsg)
	if err != nil {
		return err
	}
	fmt.Println("Hash of the message received:", hash)
	err = DepositToAccount(txmsg.BeneficiaryCustomer, txmsg.Amount)
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

func PrintMessage(mtmsg *MT103Message, spacing bool) {
	incOutGoing := "Outgoing message"
	if mtmsg.BeneficiaryInstitution == InstitutionName {
		incOutGoing = "Incoming message"
	}
	if spacing {
		fmt.Println(" -----------------")
		fmt.Printf("| %s\n| TxReferenceNumber: %s\n| TimeIndication: %s\n| BankOperationCode: %s\n| ValueDate: %s\n| Currency: %s\n| ExchangeRate: %s\n| OrderingInstitution: %s\n| BeneficiaryInstitution: %s\n| OrderingCustomer: %s\n| BeneficiaryCustomer: %s\n| Amount: %s\n",
			incOutGoing,
			mtmsg.TxReferenceNumber,
			mtmsg.TimeIndication,
			mtmsg.BankOperationCode,
			mtmsg.ValueDate,
			mtmsg.Currency,
			mtmsg.ExchangeRate,
			mtmsg.OrderingInstitution,
			mtmsg.BeneficiaryInstitution,
			mtmsg.OrderingCustomer,
			mtmsg.BeneficiaryCustomer,
			mtmsg.Amount)
		fmt.Println(" -----------------")
	} else {
		fmt.Printf("| %s | TxReferenceNumber: %s | TimeIndication: %s | BankOperationCode: %s | ValueDate: %s | Currency: %s | ExchangeRate: %s | OrderingInstitution: %s | BeneficiaryInstitution: %s | OrderingCustomer: %s | BeneficiaryCustomer: %s | Amount: %s\n",
			incOutGoing,
			mtmsg.TxReferenceNumber,
			mtmsg.TimeIndication,
			mtmsg.BankOperationCode,
			mtmsg.ValueDate,
			mtmsg.Currency,
			mtmsg.ExchangeRate,
			mtmsg.OrderingInstitution,
			mtmsg.BeneficiaryInstitution,
			mtmsg.OrderingCustomer,
			mtmsg.BeneficiaryCustomer,
			mtmsg.Amount)
	}
}

func PrintAllMessages(mtmsgs []*MT103Message) {
	fmt.Println(" -----------------")
	for _, msg := range mtmsgs {
		PrintMessage(msg, false)
	}
	fmt.Println(" -----------------")
}
