# How it works - Comprehensive explanation

Running the program launches an instance of a "bank". It is a simple command line program, and the diagnostic output is simply printed out on the console. Color coding is used to disitinguish between different pieces of information.

A bank has a database holding its clients accounts. An account is defined by the following fields:

```
Iban      string
Holder    string
Balance   float32 
Currency  string
Bic       string
Suspended bool
```

There is a special type of account, Correspondent Account, which holds the prefunded balance of a bank at another bank. Each bank keeps a mirror in its records of its Correspondent Accounts at the counterpart banks. These accounts are visualized in red and green colors, respectively.

![picture of print accounts](print_accounts.png)

It is important that the state of these special accounts is kept in sync among different banks' private records. For this purpose, every certain amount of time, banks perform a consistency check by publishing a commitment to their state on blockchain, using a challenge - response scheme. It does not reveal information about their state to other participants.

More specifically, a bank hashes the state of the mirror of its Correspondent Account at other bank. This outputs a digest that does not reveal information about the state itself. Then it hashes the digest, and submits it as a challenge to the network. Its peer bank is expected to submit the preimage of this challenge, something that can only be accomplished by having performed the same calculations as the original bank, given the 1-way nature of hash functions.

TODO link to exercise on hashing: play with hashing...

![picture of challenge - response](challenge_response.png)

### Intra-bank transaction
Intra-bank transactions move assets from one client account to another within the same bank, thus not having to interact with different banks.

### Inter-bank transaction
Inter-bank transactions move assets between client accounts at different banks. This requires the originator bank to send an off-chain message to the recipient bank. The recipient bank transfers the requested amount from sending bank Correspondent Account to the final destination customer (this is an internal transfer). The recipient bank sends a confirmation of this operation. Upon receiving the confirmation the originator bank updates its mirror account accordignly.

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

