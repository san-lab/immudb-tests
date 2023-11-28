package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/san-lab/immudb-tests/account"
	"github.com/san-lab/immudb-tests/bankinterop"
	. "github.com/san-lab/immudb-tests/datastructs"
)

var API_PORT int

type TransactionsStruct struct {
	UserFrom string `json:"userfrom"`
	Amount   string `json:"amount"`
	UserTo   string `json:"userto"`
	BankTo   string `json:"bankto"`
}

type RefillCAStruct struct {
	Amount string `json:"amount"`
	CABank string `json:"cabank"`
}

const HOME_ENDPOINT = "/"
const HEALTH_ENDPOINT = "/api/health"
const TRANSACTIONS_ENDPOINT = "/api/transactions"
const REFILL_CA_ENDPOINT = "/api/refill-ca"
const ACCOUNT_CREATION_ENDPOINT = "/api/account-creation"
const MIRROR_BALANCE_ENDPOINT = "/api/mirror-balance"

func startApiServer() {
	// HTTP routes
	http.HandleFunc(HOME_ENDPOINT, homeHandler)
	http.HandleFunc(HEALTH_ENDPOINT, healthHandler)
	http.HandleFunc(TRANSACTIONS_ENDPOINT, transactionsHandler)
	http.HandleFunc(REFILL_CA_ENDPOINT, refillCAHandler)
	http.HandleFunc(ACCOUNT_CREATION_ENDPOINT, accountCreationHandler)
	http.HandleFunc(MIRROR_BALANCE_ENDPOINT, mirrorBalanceHandler)

	fmt.Println("+ API Server running on port", API_PORT)
	err := http.ListenAndServe(fmt.Sprintf(":%d", API_PORT), nil)
	if err != nil {
		fmt.Println("Error starting the server:", err)
	}
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the home page of %s @ %s", THIS_BANK.Name, THIS_BANK.Address)
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Health check: Up!")
}

func transactionsHandler(w http.ResponseWriter, r *http.Request) {
	body, err := checkForPostRequest(w, r)
	if err != nil {
		return
	}

	// Parse array of transactions
	var transactionsBatch []TransactionsStruct
	json.Unmarshal(body, &transactionsBatch)

	fmt.Fprintf(w, "Received POST request. Processing...\n")

	for i, tx := range transactionsBatch {
		fmt.Fprintf(w, "tx %d: %s", i, tx)
		// Intrabank tx
		if tx.BankTo == THIS_BANK.Name {
			err = bankinterop.IntraBankTx(tx.UserFrom, tx.Amount, tx.UserTo)
			// Interbank tx
		} else {
			err = bankinterop.InterBankTx(tx.UserFrom, tx.Amount, tx.UserTo, tx.BankTo)
		}
		if err != nil {
			fmt.Fprintf(w, " - Error: %s\n", err.Error())
		} else {
			fmt.Fprintf(w, " - Successful!\n")
		}

	}
	fmt.Fprintf(w, "Done!")
}

func refillCAHandler(w http.ResponseWriter, r *http.Request) {
	body, err := checkForPostRequest(w, r)
	if err != nil {
		return
	}

	// Parse array of transactions
	var refillCAs []RefillCAStruct
	json.Unmarshal(body, &refillCAs)

	fmt.Fprintf(w, "Received POST request. Processing...\n")

	for i, tx := range refillCAs {
		fmt.Fprintf(w, "refill tx %d: %s", i, tx)

		err = bankinterop.RefillCA(tx.Amount, tx.CABank)
		if err != nil {
			fmt.Fprintf(w, " - Error: %s\n", err.Error())
		} else {
			fmt.Fprintf(w, " - Successful!\n")
		}

	}
	fmt.Fprintf(w, "Done!")
}

func accountCreationHandler(w http.ResponseWriter, r *http.Request) {
	body, err := checkForPostRequest(w, r)
	if err != nil {
		return
	}

	// Parse array of new accounts
	var accountCreationBatch []account.Account
	json.Unmarshal(body, &accountCreationBatch)

	fmt.Fprintf(w, "Received POST request. Processing...\n")

	for i, acc := range accountCreationBatch {
		fmt.Fprintf(w, "account %d: %s", i, acc.Holder)
		// Ignore CA fields...For now we can only create regular accounts
		err = account.CreateAccount(acc.Bic, acc.Iban, acc.Holder, acc.Currency, "", acc.Balance, acc.Suspended, false, false)

		if err != nil {
			fmt.Fprintf(w, " - Error: %s\n", err.Error())
		} else {
			fmt.Fprintf(w, " - Successful!\n")
		}
	}
	fmt.Fprintf(w, "Done!")
}

func mirrorBalanceHandler(w http.ResponseWriter, r *http.Request) {
	cABank := r.URL.Query()["cabank"]
	fmt.Println(cABank)
	//fmt.Fprintf(w, "Received POST request. Processing...\n")
	fmt.Println(cABank[0])
	fmt.Println(account.MirrorAccountIBAN(cABank[0]))
	mirrorAccount, err := account.GetAccount(account.MirrorAccountIBAN(cABank[0]))
	fmt.Println(mirrorAccount)
	fmt.Println(mirrorAccount.Balance)
	if err != nil {
		fmt.Fprintf(w, "Error %s\n", err.Error())
	} else {
		fmt.Fprintf(w, "%.2f", mirrorAccount.Balance)
	}

}

func checkForPostRequest(w http.ResponseWriter, r *http.Request) ([]byte, error) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return nil, errors.New("method not allowed")
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return nil, errors.New("error reading request body")
	}
	return body, nil
}
