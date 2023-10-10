package main

import (
	"fmt"

	"github.com/manifoldco/promptui"
)

var vSet = "VerifiedSet"
var vGet = "VerifiedGet"
var health = "Health"
var interBankTx = "Transfer to other bank client"
var intraBankTx = "Transfer to another client of the same bank"
var currentStateRoot = "Current state root"
var printAllAccounts = "Print all key-values stored"
var txById = "Get transaction by ID"
var createAccount = "Create a new account"
var setBalance = "Set the balance of an account"

func TopUI() {
	for {
		items := []string{printAllAccounts, currentStateRoot, intraBankTx, interBankTx, createAccount, setBalance, txById, vSet, vGet, health}
		items = append(items, "EXIT")
		prompt := promptui.Select{
			Label: "ImmuDB",
			Items: items,
		}
		_, it, _ := prompt.Run()

		switch it {
		case vSet:
			pr := promptui.Prompt{Label: "Introduce the key", Default: "test_key"}
			key, _ := pr.Run()
			pr = promptui.Prompt{Label: "Introduce the value", Default: "test_value"}
			value, _ := pr.Run()
			VerifiedSet(key, value)

		case vGet:
			pr := promptui.Prompt{Label: "Introduce the key", Default: "test_key"}
			key, _ := pr.Run()
			VerifiedGet(key)

		case health:
			Health()

		case interBankTx:
			pr := promptui.Prompt{Label: "Introduce the sender of the transaction", Default: "test_userFrom"}
			userFrom, _ := pr.Run()
			pr = promptui.Prompt{Label: "Introduce the amount", Default: "33"}
			amount, _ := pr.Run()
			pr = promptui.Prompt{Label: "Introduce the recipient of the transaction", Default: "test_userTo"}
			userTo, _ := pr.Run()
			InterBankTx(userFrom, amount, userTo)

		case intraBankTx:
			pr := promptui.Prompt{Label: "Introduce the sender of the transaction", Default: "test_userFrom"}
			userFrom, _ := pr.Run()
			pr = promptui.Prompt{Label: "Introduce the amount", Default: "33"}
			amount, _ := pr.Run()
			pr = promptui.Prompt{Label: "Introduce the recipient of the transaction", Default: "test_userTo"}
			userTo, _ := pr.Run()
			IntraBankTx(userFrom, amount, userTo)

		case currentStateRoot:
			root, txId, _ := CurrentStateRoot()
			fmt.Printf("Current state root: 0x%x (last tx id: %d)", root, txId)

		case printAllAccounts:
			entries, _ := GetAllAccounts()
			for _, entry := range entries.Entries {
				fmt.Printf("(%s : %s)\n", entry.Key, entry.Value)
			}

		case txById:
			pr := promptui.Prompt{Label: "Introduce the ID of the transaction", Default: "0"}
			id, _ := pr.Run()
			tx, _ := TxById(id)
			entries := tx.GetEntries()
			for i, entry := range entries {
				fmt.Printf("Tx with id %s (%d): (%s : %s)\n", id, i, entry.Key, entry.Value)
			}

		case createAccount:
			pr := promptui.Prompt{Label: "Introduce the IBAN of the new account", Default: "test_IBAN"}
			userIban, _ := pr.Run()
			pr = promptui.Prompt{Label: "Introduce the owner name of the new account", Default: "test_ownerName"}
			userName, _ := pr.Run()
			CreateAccount(userIban, userName)

		case setBalance:
			pr := promptui.Prompt{Label: "Introduce the IBAN of the account", Default: "test_IBAN"}
			userIban, _ := pr.Run()
			pr = promptui.Prompt{Label: "Introduce the new balance of the account", Default: "1"}
			balance, _ := pr.Run()
			SetBalance(userIban, balance)

		case "EXIT":
			return

		default:
			fmt.Println("u shouldnt be here...")
		}
	}
}
