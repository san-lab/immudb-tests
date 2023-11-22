package datastructs

import (
	"math/big"

	"github.com/codenotary/immudb/pkg/client"
)

var ThisBank NameAddress
var CounterpartBanks = make(map[string]string) // First entry is ThisBank

type NameAddress struct {
	Name    string
	Address string
}

var StateClient client.ImmuClient
var MsgsClient client.ImmuClient

const StateDB = "defaultdb"
const MsgsDB = "msgdb"

const Question = "question"
const Answer = "answer"

const MT103_string = "MT103"
const BankDiscoveryMessage_string = "BankDiscoveryMessage"

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

type BankDiscoveryMessage struct {
	Type              string // to prevent infinite loop
	SenderBankName    string
	SenderBankAddress string
}

type StateCheck struct {
	SubmittedHash     []byte
	SubmittedPreimage []byte
	Verified          bool
	BlockNumber       *big.Int
}
