package bankinterop

import (
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"time"

	account "github.com/san-lab/immudb-tests/account"
	"github.com/san-lab/immudb-tests/blockchainconnector"
	. "github.com/san-lab/immudb-tests/datastructs"
	sdk "github.com/san-lab/immudb-tests/immudbsdk"
)

const MT103_MESSAGE = "MT103"
const BANK_DISCOVERY_MESSAGE = "BankDiscoveryMessage"

const QUESTION = "question"
const ANSWER = "answer"

const INITIAL_AMOUNT = float32(100.0)

type BankDiscoveryMessage struct {
	Type              string // to prevent infinite loop
	SenderBankName    string
	SenderBankAddress string
}

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

// CABank -> BlockNumber -> digest
var DigestHistory = make(map[string]map[int]string)

func IntraBankTx(userFrom, amount, userTo string) error {
	err := account.WithdrawFromAccount(userFrom, amount)
	if err != nil {
		return err
	}

	err = account.DepositToAccount(userTo, amount)
	return err
}

func InterBankTx(userFrom, amount, userTo, bankTo string) error {
	_, set := COUNTERPART_BANKS[bankTo]
	if !set {
		return errors.New("cannot perform the inter bank transaction. could not find recipient bank")
	}

	err := account.WithdrawFromAccount(userFrom, amount)
	if err != nil {
		return err
	}

	// Send event to the topic and store it in MsgsDB
	txmsg := &MT103Message{TimeIndication: time.Now().String(), OrderingInstitution: THIS_BANK.Name, OrderingCustomer: userFrom, BeneficiaryInstitution: bankTo, BeneficiaryCustomer: userTo, Amount: amount}
	bytes, err := json.Marshal(txmsg)
	if err != nil {
		return err
	}
	hash, err := StoreInMsgsDB(txmsg)
	if err != nil {
		return err
	}

	// Replicate what the other bank should do with our correspondent account
	err = account.WithdrawFromAccount(account.MirrorAccountIBAN(bankTo), amount)
	if err != nil {
		return err
	}

	fmt.Println("Hash of the message sent:", hash)
	LIBP2P_NODE.SendMessage(MT103_MESSAGE, bytes)
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
	err = account.DepositToAccount(txmsg.BeneficiaryCustomer, txmsg.Amount)
	if err != nil {
		return err
	}

	// Move funds from ordering bank correspondent account
	err = account.WithdrawFromAccount(account.CAAccountIBAN(txmsg.OrderingInstitution), txmsg.Amount)
	if err != nil {
		return err
	}

	// Find out what blockNumber this new state belongs to
	blockNumber, err := blockchainconnector.GetBlockNumber()
	if err != nil {
		return err
	}
	digest, err := account.GetAccountDigest(account.CAAccountIBAN(txmsg.OrderingInstitution))
	if err != nil {
		return err
	}
	DigestHistory[txmsg.OrderingInstitution][blockNumber] = digest
	fmt.Println("debug map", DigestHistory[txmsg.OrderingInstitution])
	return err
}

func validAndAddressedToUs(txmsg *MT103Message) bool {
	if txmsg.BeneficiaryInstitution != THIS_BANK.Name {
		return false
	}

	_, err := sdk.VerifiedGet(txmsg.BeneficiaryCustomer)
	if err != nil {
		fmt.Println("Beneficiary customer is not in the database")
		return false
	}
	return true
	// TODO: check more stuff..
}

func FindCounterpartBanks() error {
	discoveryMsg := &BankDiscoveryMessage{Type: QUESTION, SenderBankName: THIS_BANK.Name, SenderBankAddress: THIS_BANK.Address}
	bytes, err := json.Marshal(discoveryMsg)
	if err != nil {
		return err
	}
	LIBP2P_NODE.SendMessage(BANK_DISCOVERY_MESSAGE, bytes)
	return nil
}

func ProcessBankDiscovery(discoveryMsg *BankDiscoveryMessage) error {
	// Pick the other bank name
	_, set := COUNTERPART_BANKS[discoveryMsg.SenderBankName]
	if !set {
		// Register he discovered bank
		COUNTERPART_BANKS[discoveryMsg.SenderBankName] = discoveryMsg.SenderBankAddress

		// Onboard the discovered bank
		err := account.CreateCAAccount("", "", discoveryMsg.SenderBankName, INITIAL_AMOUNT)
		if err != nil {
			// It means the bank has been onboarded in the DB already
			// fmt.Println(err)
			return err
		}

		// Assume the other bank has done the same, and create a mirror of our account

		err = account.CreateMirrorAccount("", "", discoveryMsg.SenderBankName, INITIAL_AMOUNT)
		if err != nil {
			// It means the bank has been onboarded in the DB already
			// fmt.Println(err)
			return err
		}

		// Initialize digest history
		DigestHistory[discoveryMsg.SenderBankName] = make(map[int]string)
		blockNumber, err := blockchainconnector.GetBlockNumber()
		if err != nil {
			return err
		}
		digest, err := account.GetAccountDigest(account.CAAccountIBAN(discoveryMsg.SenderBankName))
		if err != nil {
			return err
		}
		DigestHistory[discoveryMsg.SenderBankName][blockNumber] = digest
		fmt.Println("debug map", DigestHistory[discoveryMsg.SenderBankName])

	} else {
		// TODO handle properly
	}

	// Answer if needed
	if discoveryMsg.Type == QUESTION {
		discoveryAnswer := &BankDiscoveryMessage{Type: ANSWER, SenderBankName: THIS_BANK.Name, SenderBankAddress: THIS_BANK.Address}
		bytes, err := json.Marshal(discoveryAnswer)
		if err != nil {
			return err
		}
		LIBP2P_NODE.SendMessage(BANK_DISCOVERY_MESSAGE, bytes)
	}
	return nil
}

func HandleMessage(msgtype string, data []byte) {
	switch msgtype {

	case MT103_MESSAGE:
		txMsg := new(MT103Message)
		err := json.Unmarshal(data, txMsg)
		if err != nil {
			fmt.Println("bad frame:", err)
			return
		}
		ProcessInterBankTx(txMsg)

	case BANK_DISCOVERY_MESSAGE:
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

func PrintBankInfo() {
	fmt.Println("| Bank Name:", THIS_BANK.Name)
	fmt.Println("| Bank Address:", THIS_BANK.Address)
	fmt.Println("| ImmuDB instance running on IP:", STATE_CLIENT.GetOptions().Address)
	fmt.Println("| ImmuDB instance running on port:", STATE_CLIENT.GetOptions().Port)
	fmt.Println("| ...")
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

func PickLatestDigestPriorToResquestedBlockNumber(cABank string, blockNumber *big.Int) (string, error) {
	digest := ""
	number := int(blockNumber.Int64())
	for digest == "" && number > 0 {
		number = number - 1
		digest = DigestHistory[cABank][number]
	}
	fmt.Println("- debug pick digest:", cABank, number, digest)
	if digest == "" {
		return "", errors.New("couldnt find a digest for the block number requested")
	}
	return digest, nil
}
