package transactions

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	account "github.com/san-lab/immudb-tests/account"
	. "github.com/san-lab/immudb-tests/datastructs"
	sdk "github.com/san-lab/immudb-tests/immudbsdk"
)

func InterBankTx(userFrom, amount, userTo, bankTo string) error { // TODO maybe another struct for the parameters
	err := account.WithdrawFromAccount(userFrom, amount)
	if err != nil {
		return err
	}

	// Send event to the topic and store it in MsgsDB
	txmsg := &MT103Message{TimeIndication: time.Now().String(), OrderingInstitution: ThisBank.Name, OrderingCustomer: userFrom, BeneficiaryInstitution: bankTo, BeneficiaryCustomer: userTo, Amount: amount}
	bytes, err := json.Marshal(txmsg)
	if err != nil {
		return err
	}
	hash, err := sdk.StoreInMsgsDB(txmsg)
	if err != nil {
		return err
	}

	// Replicate what the other bank should do with our correspondent account
	mirrorAccount := ThisBank.Name + "@" + bankTo + " - Mirror"
	err = account.WithdrawFromAccount(mirrorAccount, amount)
	if err != nil {
		return err
	}

	fmt.Println("Hash of the message sent:", hash)
	LibP2PNode.SendMessage(MT103_string, bytes)
	return nil
}

// When receiveing a transaction
func ProcessInterBankTx(txmsg *MT103Message) error {
	if !validAndAddressedToUs(txmsg) {
		return errors.New("received transaction message is invalid")
	}
	hash, err := sdk.StoreInMsgsDB(txmsg)
	if err != nil {
		return err
	}
	fmt.Println("Hash of the message received:", hash)
	err = account.DepositToAccount(txmsg.BeneficiaryCustomer, txmsg.Amount)
	if err != nil {
		return err
	}

	// Move funds from ordering bank correspondent account
	correspondentAccount := txmsg.OrderingInstitution + " - CA"
	err = account.WithdrawFromAccount(correspondentAccount, txmsg.Amount)
	if err != nil {
		return err
	}

	return err
}

func IntraBankTx(userFrom, amount, userTo string) error {
	err := account.WithdrawFromAccount(userFrom, amount)
	if err != nil {
		return err
	}

	err = account.DepositToAccount(userTo, amount)
	return err
}

func validAndAddressedToUs(txmsg *MT103Message) bool {
	if txmsg.BeneficiaryInstitution != ThisBank.Name {
		return false
	}

	_, err := sdk.VerifiedGet(txmsg.BeneficiaryCustomer)
	if err != nil {
		fmt.Println("Beneficiary customer is not in the database")
		return false
	}
	return true
	// TODO check more stuff..
}

func FindCounterpartBanks() error {
	discoveryMsg := &BankDiscoveryMessage{Type: Question, SenderBankName: ThisBank.Name, SenderBankAddress: ThisBank.Address}
	bytes, err := json.Marshal(discoveryMsg)
	if err != nil {
		return err
	}
	LibP2PNode.SendMessage(BankDiscoveryMessage_string, bytes)
	return nil
}

func ProcessBankDiscovery(discoveryMsg *BankDiscoveryMessage) error {
	// Pick the other bank name
	_, set := CounterpartBanks[discoveryMsg.SenderBankName]
	if !set {
		initialAmount := "100"

		// Register he discovered bank
		CounterpartBanks[discoveryMsg.SenderBankName] = discoveryMsg.SenderBankAddress

		// Onboard the discovered bank
		accName := discoveryMsg.SenderBankName + " - CA"
		err := account.CreateAccount(accName, accName)
		if err != nil {
			// It means the bank has been onboarded in the DB already
			// fmt.Println(err)
			return err
		}
		err = account.SetAccountBalance(accName, initialAmount)
		if err != nil {
			fmt.Println(err, "cannot set CA balance")
		}

		// Assume the other bank has done the same, and create a mirror of our account
		accName = ThisBank.Name + "@" + discoveryMsg.SenderBankName + " - Mirror"
		err = account.CreateAccount(accName, accName)
		if err != nil {
			// It means the bank has been onboarded in the DB already
			// fmt.Println(err)
			return err
		}
		err = account.SetAccountBalance(accName, initialAmount)
		if err != nil {
			fmt.Println(err, "cannot set mirror account balance")
		}

	} else {
		// TODO handle properly
	}

	// Answer if needed
	if discoveryMsg.Type == Question {
		discoveryAnswer := &BankDiscoveryMessage{Type: Answer, SenderBankName: ThisBank.Name, SenderBankAddress: ThisBank.Address}
		bytes, err := json.Marshal(discoveryAnswer)
		if err != nil {
			return err
		}
		LibP2PNode.SendMessage(BankDiscoveryMessage_string, bytes)
	}
	return nil
}

func HandleMessage(msgtype string, data []byte) {
	switch msgtype {

	case MT103_string:
		txMsg := new(MT103Message)
		err := json.Unmarshal(data, txMsg)
		if err != nil {
			fmt.Println("bad frame:", err)
			return
		}
		ProcessInterBankTx(txMsg)

	case BankDiscoveryMessage_string:
		discoveryMsg := new(BankDiscoveryMessage)
		err := json.Unmarshal(data, discoveryMsg)
		if err != nil {
			fmt.Println("bad frame:", err)
			return
		}
		ProcessBankDiscovery(discoveryMsg)

	default:
		fmt.Println("u shouldnt be here...")
	}
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

func PrintBankInfo() {
	fmt.Println("| Bank Name:", ThisBank.Name)
	fmt.Println("| Bank Address:", ThisBank.Address)
	fmt.Println("| ImmuDB instance running on IP:", StateClient.GetOptions().Address)
	fmt.Println("| ImmuDB instance running on port:", StateClient.GetOptions().Port)
	fmt.Println("| ...")
}

func PrintMessage(mtmsg *MT103Message, spacing bool) {
	incOutGoing := "Outgoing message"
	if mtmsg.BeneficiaryInstitution == ThisBank.Name {
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
