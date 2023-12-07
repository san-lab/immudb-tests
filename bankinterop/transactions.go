package bankinterop

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/san-lab/immudb-tests/account"
	"github.com/san-lab/immudb-tests/blockchainconnector"
	"github.com/san-lab/immudb-tests/color"
	. "github.com/san-lab/immudb-tests/datastructs"
)

const MT103_MESSAGE = "MT103"
const MT103_CONFIRMATION = "MT103Confirmation"
const REFILL_CA_MESSAGE = "RefillCAMessage"
const REFILL_CA_CONFIRMATION = "RefillCAConfirmation"
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

	ReferenceBlockNumber int
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

	ReferenceBlockNumber int
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

// TODO: keep track of pending Interbank and Refill pending responses with a reference number in the message. Dont act on anything that is not pending
func RequestInterBankTx(userFrom, amount, userTo, bankTo string) error {
	_, set := COUNTERPART_BANKS[bankTo]
	if !set {
		return errors.New("cannot perform the inter bank transaction. could not find recipient bank")
	}

	// Check balance at mirror, and take action
	if !isMirrorBalanceEnough(bankTo, amount) {
		return errors.New("not enough balance at correspondent account")
	}

	refBlockNumber, err := blockchainconnector.GetBlockNumber()
	// fmt.Println("debug interbank tx refBlkNumber", refBlockNumber)
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
		Amount:                 amount,
		ReferenceBlockNumber:   refBlockNumber,
	}
	bytes, err := json.Marshal(txmsg)
	if err != nil {
		return err
	}
	LIBP2P_NODE.SendMessage(MT103_MESSAGE, bankTo, bytes)
	// fmt.Println("debug tx sent", txmsg)
	/*
		hash, err := StoreInMsgsDB(txmsg)
		if err != nil {
			return err
		}
		fmt.Println("Hash of the message sent:", hash)

		err = account.WithdrawFromAccount(userFrom, amount)
		if err != nil {
			return err
		}

		// Replicate what the other bank should do with our correspondent account
		err = account.WithdrawFromAccount(account.MirrorAccountIBAN(bankTo), amount)
		if err != nil {
			return err
		}
	*/
	return nil
}

// When receiveing a transaction
func ProcessInterBankTx(txmsg *MT103Message) error {
	if !validInterBankTx(txmsg) {
		return errors.New("received transaction message is invalid or not addressed to us")
	}

	hash, err := StoreInMsgsDB(txmsg)
	if err != nil {
		return err
	}
	fmt.Println("Hash of the transaction received:", hash)

	err = account.DepositToAccount(txmsg.BeneficiaryCustomer, txmsg.Amount)
	if err != nil {
		return err
	}

	// Move funds from ordering bank correspondent account
	err = account.WithdrawFromAccount(account.CAAccountIBAN(txmsg.OrderingInstitution), txmsg.Amount)
	if err != nil {
		return err
	}

	bytes, err := json.Marshal(txmsg)
	if err != nil {
		return err
	}
	LIBP2P_NODE.SendMessage(MT103_CONFIRMATION, txmsg.OrderingInstitution, bytes)
	// fmt.Println("debug tx received", txmsg)

	err = updateCADigestHistory(txmsg.OrderingInstitution, txmsg.ReferenceBlockNumber)
	return err
}

func ConfirmedInterBankTx(txmsg *MT103Message) error {
	if !validInterBankTxConfirmation(txmsg) {
		return errors.New("received transaction message is invalid or not addressed to us")
	}
	// fmt.Println("debug tx confirmation", txmsg)

	hash, err := StoreInMsgsDB(txmsg)
	if err != nil {
		return err
	}
	fmt.Println("Hash of the transaction sent (confirmed):", hash)

	err = account.WithdrawFromAccount(txmsg.OrderingCustomer, txmsg.Amount)
	if err != nil {
		return err
	}

	// Replicate what the other bank should do with our correspondent account
	err = account.WithdrawFromAccount(account.MirrorAccountIBAN(txmsg.BeneficiaryInstitution), txmsg.Amount)
	if err != nil {
		return err
	}
	return nil
}

/*
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
*/

func RequestRefillCA(amount, bankTo string) error {
	refBlockNumber, err := blockchainconnector.GetBlockNumber()
	fmt.Println("debug refill ca refBlkNumber", refBlockNumber)
	if err != nil {
		return err
	}
	txmsg := &RefillCAMessage{
		TimeIndication:         time.Now().String(),
		OrderingInstitution:    THIS_BANK.Name,
		BeneficiaryInstitution: bankTo,
		Amount:                 amount,
		ReferenceBlockNumber:   refBlockNumber,
	}
	bytes, err := json.Marshal(txmsg)
	if err != nil {
		return err
	}
	LIBP2P_NODE.SendMessage(REFILL_CA_MESSAGE, bankTo, bytes)
	fmt.Printf("Sent refill CA request to %s (%s)\n", bankTo, amount)
	// fmt.Println("debug refill sent", txmsg)
	return nil
}

func ProcessRefillCA(refillMsg *RefillCAMessage) error {
	if !validRefillCA(refillMsg) {
		return errors.New("received refill message is invalid or not addressed to us")
	}

	err := account.DepositToAccount(account.CAAccountIBAN(refillMsg.OrderingInstitution), refillMsg.Amount)
	if err != nil {
		return err
	}

	fmt.Printf("Received and processed refill CA request from %s (%s)\n", refillMsg.OrderingInstitution, refillMsg.Amount)

	bytes, err := json.Marshal(refillMsg)
	if err != nil {
		return err
	}
	LIBP2P_NODE.SendMessage(REFILL_CA_CONFIRMATION, refillMsg.OrderingInstitution, bytes)
	// fmt.Println("debug refill received", refillMsg)

	err = updateCADigestHistory(refillMsg.OrderingInstitution, refillMsg.ReferenceBlockNumber)
	return err
}

func ConfirmedRefillCA(refillMsg *RefillCAMessage) error {
	if !validRefillCAConfirmation(refillMsg) {
		return errors.New("received refill confirmation message is invalid or not addressed to us")
	}
	// fmt.Println("debug refill confirmation", refillMsg)
	err := account.DepositToAccount(account.MirrorAccountIBAN(refillMsg.BeneficiaryInstitution), refillMsg.Amount)
	return err
}

func FindCounterpartBanks() error {
	discoveryMsg := &BankDiscoveryMessage{Type: QUESTION, SenderBankName: THIS_BANK.Name, SenderBankAddress: THIS_BANK.Address}
	bytes, err := json.Marshal(discoveryMsg)
	if err != nil {
		return err
	}
	LIBP2P_NODE.SendMessage(BANK_DISCOVERY_MESSAGE, "", bytes)
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
		LIBP2P_NODE.SendMessage(BANK_DISCOVERY_MESSAGE, "", bytes)
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

	case MT103_CONFIRMATION:
		txMsg := new(MT103Message)
		err := json.Unmarshal(data, txMsg)
		if err != nil {
			fmt.Println("bad frame:", err)
			return
		}
		ConfirmedInterBankTx(txMsg)

	case REFILL_CA_MESSAGE:
		txMsg := new(RefillCAMessage)
		err := json.Unmarshal(data, txMsg)
		if err != nil {
			fmt.Println("bad frame:", err)
			return
		}
		ProcessRefillCA(txMsg)

	case REFILL_CA_CONFIRMATION:
		txMsg := new(RefillCAMessage)
		err := json.Unmarshal(data, txMsg)
		if err != nil {
			fmt.Println("bad frame:", err)
			return
		}
		ConfirmedRefillCA(txMsg)

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
