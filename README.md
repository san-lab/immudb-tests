# Reconciliation Bank Demo

## How it works - Comprehensive explanation

Running the program launches an instance of a "bank". It is a simple command line program, and the diagnostic output is simply printed out on the console. User can issue commands through the interactive prompt interface to create and manage accounts, send transactions between different customers, etc. Color coding is used to disitinguish between different pieces of information.

Each bank has a database holding its customer accounts. An account is defined by the following fields:

```
Iban      string
Holder    string
Balance   float32 
Currency  string
Bic       string
Suspended bool
```

There is a special type of account, Correspondent Account, which holds the prefunded balance of a bank at another bank. Each bank keeps a mirror in its records of its Correspondent Accounts at the counterpart banks. These accounts are visualized in red and green colors, respectively.

![picture of print accounts](./docs/print_accounts.png)

It is important that the state of these special accounts is kept in sync among different banks' private records. For this purpose, every certain amount of time, banks perform a consistency check by publishing a commitment to their state on blockchain, using a challenge - response scheme. It does not reveal information about their state to other participants.

More specifically, a bank hashes the state of the mirror of its Correspondent Account at other bank. This outputs a digest that does not reveal information about the state itself. Then it hashes the digest, and submits it as a challenge to the network. Its peer bank is expected to submit the preimage of this challenge, something that can only be accomplished by having performed the same calculations as the original bank, given the 1-way nature of hash functions.

> **Note:** the reader not familiarized with hash functions might want to play with some online tool like [Keccak256 Hash Online](https://emn178.github.io/online-tools/keccak_256.html). For every input, a different, constant-size, randomly-looking string is outputted (often called *digest*). It is unfeasible to come up with an input that hashes to a provided output.

![picture of challenge - response](./docs/challenge_response.png)

This periodic polling and submitting information to blockchain is shown in white and light blue color logs in the program.

### Intra-bank transaction
Intra-bank transactions move assets from one client account to another within the same bank, thus not having to interact with different banks.

### Inter-bank transaction
Inter-bank transactions move assets between client accounts at different banks. This requires the originator bank to send an off-chain message to the recipient bank. The recipient bank transfers the requested amount from sending bank Correspondent Account to the final destination customer (this is an internal transfer). The recipient bank sends a confirmation of this operation. Upon receiving the confirmation the originator bank updates its mirror account accordingly.

The format of the Inter-bank transaction message is similar to Swift message MT103:
```
TxReferenceNumber      string
TimeIndication         string
BankOperationCode      string
ValueDate              string
Currency               string
ExchangeRate           string
OrderingInstitution    string
BeneficiaryInstitution string

OrderingCustomer    string
BeneficiaryCustomer string
Amount              string

ReferenceBlockNumber int
```

#### Complete proccess diagram
![picture of the complete process](./docs/demo_bank_reconciliation.png)

---

### Current Deployment

In the current deployment, we are using the new Ethereum test network, [Holesky](https://holesky.etherscan.io/). Our verifier [contract](https://github.com/san-lab/immudb-tests/blockchainconnector/onchainverifier.sol) is deployed in the blockchain address `0x5a0F1c0A4482a6CE88C190dE396d154A2149544a`. The transaction related to the contract can be inspected with the public blockchain [explorer](https://holesky.etherscan.io/) ([Link to contract](https://holesky.etherscan.io/address/0x5a0F1c0A4482a6CE88C190dE396d154A2149544a)).

The sample bank blockchain addresses are the following:

Bank  | Address
------------- | -------------
Santa Bank | [0x7eC027cF7f470983030167d2FACE94745E1AFfE3](https://holesky.etherscan.io/address/0x7eC027cF7f470983030167d2FACE94745E1AFfE3)
Blue Bank | [0x6e7786c888Fe08E9360E830bC5806eca6186fB89](https://holesky.etherscan.io/address/0x6e7786c888Fe08E9360E830bC5806eca6186fB89)
Green Bank | [0x3e2a6b7E74447bC16c10E1a5E5da7D1af5e5c2e3](https://holesky.etherscan.io/address/0x3e2a6b7E74447bC16c10E1a5E5da7D1af5e5c2e3)

These account have been prefunded with some Ethereum test tokens for the network gas fees. The corresponding private keys are configured in the deployed environment.

---

# Setup

0. Setup ImmuDB https://docs.immudb.io/0.9.2/quickstart.html

1. Delete `data` directories to reset the databases

2. Start multiple ImmuDB instances (with different data directories, and running on different ports). Parameters: https://docs.immudb.io/1.2.1/reference/configuration.html
```
./immudb
./immudb --dir=./data2 --port=3323 --web-server-port=8081 --pgsql-server-port=5433
```
For some reason, the web server still points to the first instance...

3. Delete `.identity-XXXXX` and `.state-XXXXX` files (generated by the clients) if reseted the databases

4. Start multiple ImmuDB clients pointing to the respectives databases
```
./binary_name
./binary_name --config config/config2.env
```

Config files must have the following data:
```
BANK_NAME        string
BANK_ADDRESS     string

DB_IP            string
DB_PORT          int

LIBP2P_TOPIC     string
API_PORT         int

NETWORK          string
CHAIN_ID         string
VERIFIER_ADDRESS string
PRIV_KEY_FILE    string

UPDATE_FREQUENCY int
POLL_FREQUENCY   int

FIND_FREQUENCY   int

LOG_FILE         string
```

Sample:
```
BANK_NAME="Santa Bank"
BANK_ADDRESS="0x7eC027cF7f470983030167d2FACE94745E1AFfE3"

DB_IP="127.0.0.1"
DB_PORT=3322

LIBP2P_TOPIC="ImmuDBTopic"
API_PORT=3301

NETWORK="https://ethereum-holesky.publicnode.com/"
CHAIN_ID="17000"
VERIFIER_ADDRESS="0xff330b4c20d04602f2e522354a1e1d2c91835a7f"
PRIV_KEY_FILE="config/priv_key.txt"

UPDATE_FREQUENCY=60
POLL_FREQUENCY=25

FIND_FREQUENCY=20

LOG_FILE="chain-logs.txt"
```
Private key must be hex encoded, **without** *"0x"* character **nor** quotation marks.

5. Deploy ```blockchainconnector/onchainverfier.go``` on your blockchain network of choice (compiler version 0.5.0)


# APIs

### GET /
Will answer with the bank name and address

### GET /api/health
Will answer ```Up!``` if the server is running


### POST /api/transactions
Provide an array of structs with the following data. *BankTo* field determines whether is an intrabank or interbank transaction:
```
type TransactionsStruct []struct 
{
	UserFrom string `json:"userfrom"`
	Amount   string `json:"amount"`
	UserTo   string `json:"userto"`
	BankTo   string `json:"bankto"`
}
```

Sample: 
```
[
{
    "userfrom": "UserFromIBAN",
    "amount":   "10",
    "userto":   "UserToIBAN",
    "bankto":   "MyOwnBank"
},
{
    "userfrom": "UserFromIBAN",
    "amount":   "2",
    "userto":   "UserToIBAN",
    "bankto":   "CounterpartBank"
},
...
]
```

### POST /api/account-creation
Provide an array of structs with the following data. Correspondent account fields (*isCA*, *isMirror*, *CABank*) will be ignored for now.
```
type Account struct 
{
	Suspended bool    `json:"suspended"`
	Bic       string  `json:"bic"`
	Iban      string  `json:"iban"`
	Balance   float32 `json:"balance"`
	Holder    string  `json:"holder"`
	Currency  string  `json:"currency"`
	IsCA      bool    `json:"isca"`
	IsMirror  bool    `json:"ismirror"`
	CABank    string  `json:"cabank"`
}
```

Sample:
```
[
{
	"suspended": false,
	"bic":       "",
	"iban":      "NewUserIBAN",
	"balance":   15.5,
	"holder":    "NewUserName",
	"currency":  "",
	"isca":      "false",
	"ismirror":  "false",
	"cabank":    ""
},
...
]
```

### POST /api/refill-ca
Refill correspondent account at another bank. Provide an array of structs with the following data.
```
type RefillCAStruct struct {
	Amount string `json:"amount"`
	CABank string `json:"cabank"`
}
```

Sample:
```
[
{
    "amount": "10",
    "cabank": "Test Bank"
},
{
    "amount": "1",
    "cabank": "Test Bank 2"
}
]
```
