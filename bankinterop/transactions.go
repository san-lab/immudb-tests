package bankinterop

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"strconv"
	"time"

	"github.com/san-lab/immudb-tests/account"
	"github.com/san-lab/immudb-tests/blockchainconnector"
	"github.com/san-lab/immudb-tests/color"
	. "github.com/san-lab/immudb-tests/datastructs"
)

const MT103_MESSAGE = "MT103"
const REFILL_CA_MESSAGE = "RefillCAMessage"
const BANK_DISCOVERY_MESSAGE = "BankDiscoveryMessage"

const QUESTION = "question"
const ANSWER = "answer"

const INITIAL_AMOUNT = float32(100.0)
const DEBT_LIMIT = -100.0

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

// Refill CA account
type RefillCAMessage struct {
	TxReferenceNumber string
	TimeIndication    string
	ValueDate         string // always today?
	Currency          string
	ExchangeRate      string

	OrderingInstitution    string
	BeneficiaryInstitution string
	Amount                 string
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

	// Check balance at mirror, and take action
	mirrorAccount, _ := account.GetAccount(account.MirrorAccountIBAN(bankTo))
	amountFloat, _ := strconv.ParseFloat(amount, 32)
	mirrorBalance := mirrorAccount.Balance - float32(amountFloat)
	if mirrorBalance < DEBT_LIMIT {
		return errors.New("cannot perform the inter bank transaction. our correspondent account at counterpart bank would be over the limit")
	} else if mirrorBalance < 0 {
		color.CPrintln(color.RED, "Warning: balance of our mirror account @ %s is %.2f", bankTo, mirrorBalance)
	}

	err := account.WithdrawFromAccount(userFrom, amount)
	if err != nil {
		return err
	}

	// Send event to the topic and store it in MsgsDB
	txmsg := &MT103Message{
		TimeIndication:         time.Now().String(),
		OrderingInstitution:    THIS_BANK.Name,
		OrderingCustomer:       userFrom,
		BeneficiaryInstitution: bankTo,
		BeneficiaryCustomer:    userTo,
		Amount:                 amount}
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
	if !validAndAddressedToUsMT103(txmsg) {
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

	err = updateCADigestHistory(txmsg.OrderingInstitution)
	return err
}

func RefillCA(amount, bankTo string) error {
	// TODO: dont update until confirmation
	err := account.DepositToAccount(account.MirrorAccountIBAN(bankTo), amount)
	if err != nil {
		return err
	}
	txmsg := &RefillCAMessage{
		TimeIndication:         time.Now().String(),
		OrderingInstitution:    THIS_BANK.Name,
		BeneficiaryInstitution: bankTo,
		Amount:                 amount}

	bytes, err := json.Marshal(txmsg)
	if err != nil {
		return err
	}
	LIBP2P_NODE.SendMessage(REFILL_CA_MESSAGE, bytes)
	return nil
}

func ProcessRefillCA(refillMsg *RefillCAMessage) error {
	if !validAndAddressedToUsRefillCA(refillMsg) {
		return errors.New("received refill message is invalid")
	}

	err := account.DepositToAccount(account.CAAccountIBAN(refillMsg.OrderingInstitution), refillMsg.Amount)
	if err != nil {
		return err
	}

	err = updateCADigestHistory(refillMsg.OrderingInstitution)
	return err
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

	case REFILL_CA_MESSAGE:
		txMsg := new(RefillCAMessage)
		err := json.Unmarshal(data, txMsg)
		if err != nil {
			fmt.Println("bad frame:", err)
			return
		}
		ProcessRefillCA(txMsg)

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

func PickLatestDigestPriorToResquestedBlockNumber(cABank string, blockNumber *big.Int) (string, error) {
	digest := ""
	number := int(blockNumber.Int64())
	for digest == "" && number > 0 {
		number = number - 1
		digest = DigestHistory[cABank][number]
	}
	color.CPrintln(color.WHITE, "- Picked digest: %s, %d, %s", cABank, number, color.Shorten(digest, 10))
	if digest == "" {
		return "", errors.New("couldnt find a digest for the block number requested")
	}
	return digest, nil
}
