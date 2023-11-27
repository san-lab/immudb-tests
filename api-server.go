package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/san-lab/immudb-tests/account"
	"github.com/san-lab/immudb-tests/bankinterop"
	. "github.com/san-lab/immudb-tests/datastructs"
)

var API_PORT int

type TransactionsStruct []struct {
	UserFrom string `json:"userfrom"`
	Amount   string `json:"amount"`
	UserTo   string `json:"userto"`
	BankTo   string `json:"bankto"`
}

const HOME_ENDPOINT = "/"
const HEALTH_ENDPOINT = "/api/health"
const TRANSACTIONS_ENDPOINT = "/api/transactions"
const ACCOUNT_CREATION_ENDPOINT = "/api/account-creation"

func startApiServer() {
	// Define your HTTP routes
	http.HandleFunc(HOME_ENDPOINT, homeHandler)
	http.HandleFunc(HEALTH_ENDPOINT, healthHandler)
	http.HandleFunc(TRANSACTIONS_ENDPOINT, transactionsHandler)
	http.HandleFunc(ACCOUNT_CREATION_ENDPOINT, accountCreationHandler)

	fmt.Printf("Server is running on :%d...\n", API_PORT)
	err := http.ListenAndServe(fmt.Sprintf(":%d", API_PORT), nil)
	if err != nil {
		fmt.Println("Error starting the server:", err)
	}
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the home page of %s @ %s", THIS_BANK.Name, THIS_BANK.Address)
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Up!")
}

func transactionsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}

	// Parse array of transactions
	var transactionsBatch TransactionsStruct
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

func accountCreationHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}

	// Parse array of new accounts
	var accountCreationBatch []account.Account
	json.Unmarshal(body, &accountCreationBatch)
	fmt.Println(accountCreationBatch)
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
