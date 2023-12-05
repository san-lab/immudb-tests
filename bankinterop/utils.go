package bankinterop

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/san-lab/immudb-tests/account"
	"github.com/san-lab/immudb-tests/color"
	. "github.com/san-lab/immudb-tests/datastructs"
	sdk "github.com/san-lab/immudb-tests/immudbsdk"
)

func validInterBankTx(txmsg *MT103Message) bool {
	if txmsg.BeneficiaryInstitution != THIS_BANK.Name {
		return false
	}

	_, err := sdk.VerifiedGet(txmsg.BeneficiaryCustomer)
	if err != nil {
		color.CPrintln(color.RED, "%s is not in the database", txmsg.BeneficiaryCustomer)
		return false
	}

	_, err = sdk.VerifiedGet(account.CAAccountIBAN(txmsg.OrderingInstitution))
	if err != nil {
		color.CPrintln(color.RED, "%s ordering bank correspondent account is not in the database", account.CAAccountIBAN(txmsg.OrderingInstitution))
		return false
	}
	return true
}

func validInterBankTxConfirmation(txmsg *MT103Message) bool {
	if txmsg.OrderingInstitution != THIS_BANK.Name {
		return false
	}

	_, err := sdk.VerifiedGet(txmsg.OrderingCustomer)
	if err != nil {
		color.CPrintln(color.RED, "%s is not in the database", txmsg.OrderingCustomer)
		return false
	}

	_, err = sdk.VerifiedGet(account.MirrorAccountIBAN(txmsg.BeneficiaryInstitution))
	if err != nil {
		color.CPrintln(color.RED, "%s mirror account of recipient bank CA is not in the database", account.MirrorAccountIBAN(txmsg.BeneficiaryInstitution))
		return false
	}

	return true
}

func isMirrorBalanceEnough(bankTo, amount string) bool {
	mirrorAccount, _ := account.GetAccount(account.MirrorAccountIBAN(bankTo))
	amountFloat, _ := strconv.ParseFloat(amount, 32)
	finalMirrorBalance := mirrorAccount.Balance - float32(amountFloat)
	if finalMirrorBalance < DEBT_LIMIT {
		color.CPrintln(color.RED, "Error: cannot perform the inter bank transaction. our correspondent account at counterpart bank would be over the limit")
		return false
	} else if finalMirrorBalance < 0 {
		color.CPrintln(color.RED, "Warning: balance of our mirror account @ %s is %.2f", bankTo, finalMirrorBalance)
	}
	return true
}

func validRefillCA(refillMsg *RefillCAMessage) bool {
	if refillMsg.BeneficiaryInstitution != THIS_BANK.Name {
		return false
	}

	_, err := sdk.VerifiedGet(account.CAAccountIBAN(refillMsg.OrderingInstitution))
	if err != nil {
		color.CPrintln(color.RED, "%s correspondent account is not in the database", refillMsg.OrderingInstitution)
		return false
	}
	return true
}

func validRefillCAConfirmation(refillMsg *RefillCAMessage) bool {
	if refillMsg.OrderingInstitution != THIS_BANK.Name {
		return false
	}
	_, err := sdk.VerifiedGet(account.MirrorAccountIBAN(refillMsg.BeneficiaryInstitution))
	if err != nil {
		color.CPrintln(color.RED, "%s mirror account is not in the database", refillMsg.BeneficiaryInstitution)
		return false
	}
	return true
}

func updateCADigestHistory(cABank string, refBlockNumber int) error {
	// Find out what blockNumber this new state belongs to
	/*
		blockNumber, err := blockchainconnector.GetBlockNumber()
		if err != nil {
			return err
		}
	*/
	digest, err := account.GetAccountDigest(account.CAAccountIBAN(cABank))
	if err != nil {
		return err
	}
	DigestHistory[cABank][refBlockNumber] = digest
	return nil
}

func GetMessage(key string) (*MT103Message, error) {
	// Pick message
	messageRaw, err := sdk.VerifiedGetMsg(key)
	if err != nil {
		return nil, err
	}
	message := new(MT103Message)
	err = json.Unmarshal(messageRaw.Value, message)
	return message, err
}

func GetAllMessages() ([]*MT103Message, error) {
	entries, err := sdk.GetAllMsgsEntries()
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

func StoreInMsgsDB(txmsg *MT103Message) (string, error) {
	value, err := json.Marshal(txmsg)
	if err != nil {
		return "", err
	}
	hash := sha256.Sum256(value)
	key := fmt.Sprintf("0x%x", hash[:])
	err = sdk.VerifiedSetMsg(key, string(value))
	return key, err
}

func PrintMessage(mtmsg *MT103Message, spacing bool) {
	incOutGoing := "Outgoing message"
	if mtmsg.BeneficiaryInstitution == THIS_BANK.Name {
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
