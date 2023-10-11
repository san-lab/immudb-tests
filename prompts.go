package main

import (
	"fmt"

	"github.com/manifoldco/promptui"
)

const UP = "UP"
const EXIT = "EXIT"

var bankInfo = "Show bank information"
var findCounterpartBanks = "Broadcast call to find other banks"
var seeCounterpartBanks = "See current list of counterpart banks"

var interBankTx = "Transfer to other bank client"
var intraBankTx = "Transfer to another client of the same bank"

var currentStateRoot = "Current state root"
var printAllAccounts = "Print all key-values stored"
var printAccount = "Print all values of an account"

var manageAccount = "Manage an account"
var createAccount = "Create a new account"
var setAccountBalance = "Set the balance of an account"
var depositToAccount = "Deposit the specified amount to an account"
var withdrawFromAccount = "Withdraw the specified amount from an account"
var suspendAccount = "Suspend an account"
var unsuspendAccount = "Unsuspend an account"

var health = "Health"
var vSet = "VerifiedSet"
var vGet = "VerifiedGet"
var txById = "Get transaction by ID"

func TopUI() {
	for {
		items := []string{bankInfo, findCounterpartBanks, seeCounterpartBanks, printAllAccounts, printAccount, currentStateRoot, intraBankTx, interBankTx,
			manageAccount, txById, health, vSet, vGet}

		items = append(items, EXIT)
		prompt := promptui.Select{
			Label: InstitutionName + " - Actions",
			Items: items,
		}
		_, it, _ := prompt.Run()

		switch it {

		case bankInfo:
			PrintBankInfo()

		case findCounterpartBanks:
			err := FindCounterpartBanks()
			if err != nil {
				fmt.Println(err)
			}

		case seeCounterpartBanks:
			fmt.Println("Current list of banks: ", CounterpartBanks)

		case interBankTx:
			pr := promptui.Prompt{Label: "Introduce the sender of the transaction", Default: "test_userFrom"}
			userFrom, _ := pr.Run()
			pr = promptui.Prompt{Label: "Introduce the amount", Default: "33"}
			amount, _ := pr.Run()
			pr = promptui.Prompt{Label: "Introduce the recipient of the transaction", Default: "test_userTo"}
			userTo, _ := pr.Run()

			// TODO add a prompt to select from a list of banks
			// CounterpartBanks
			pr = promptui.Prompt{Label: "Introduce the institution recipient of the transaction", Default: "test_bankTo"}
			bankTo, _ := pr.Run()

			err := InterBankTx(userFrom, amount, userTo, bankTo)
			if err != nil {
				fmt.Println(err)
			}

		case intraBankTx:
			pr := promptui.Prompt{Label: "Introduce the sender of the transaction", Default: "test_userFrom"}
			userFrom, _ := pr.Run()
			pr = promptui.Prompt{Label: "Introduce the amount", Default: "33"}
			amount, _ := pr.Run()
			pr = promptui.Prompt{Label: "Introduce the recipient of the transaction", Default: "test_userTo"}
			userTo, _ := pr.Run()

			err := IntraBankTx(userFrom, amount, userTo)
			if err != nil {
				fmt.Println(err)
			}

		case currentStateRoot:
			root, txId, _ := CurrentStateRoot()
			fmt.Printf("Current state root: 0x%x (last tx id: %d)", root, txId)

		case printAllAccounts:
			accounts, _ := GetAllAccounts()
			PrintAllAccounts(accounts)

		case printAccount:
			pr := promptui.Prompt{Label: "Introduce the IBAN of the new account", Default: "test_IBAN"}
			userIban, _ := pr.Run()

			account, err := GetAccount(userIban)
			if err != nil {
				fmt.Println(err)
			}
			PrintAccount(account, true)

		case txById:
			pr := promptui.Prompt{Label: "Introduce the ID of the transaction", Default: "0"}
			id, _ := pr.Run()
			tx, _ := TxById(id)
			entries := tx.GetEntries()
			for i, entry := range entries {
				fmt.Printf("Tx with id %s (%d): (%s : %s)\n", id, i, entry.Key, entry.Value)
			}

		case manageAccount:
			pr := promptui.Prompt{Label: "Introduce the IBAN of the account to manage", Default: "test_IBAN"}
			userIban, _ := pr.Run()
			ManageAccountUI(userIban)

		case vSet:
			pr := promptui.Prompt{Label: "Introduce the key", Default: "test_key"}
			key, _ := pr.Run()
			pr = promptui.Prompt{Label: "Introduce the value", Default: "test_value"}
			value, _ := pr.Run()

			err := VerifiedSet(key, value)
			if err != nil {
				fmt.Println(err)
			}

		case vGet:
			pr := promptui.Prompt{Label: "Introduce the key", Default: "test_key"}
			key, _ := pr.Run()

			entry, err := VerifiedGet(key)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Printf("Sucessfully got verified entry: ('%s', '%s') @ tx %d\n", entry.Key, entry.Value, entry.Tx)

		case health:
			health, err := Health()
			if err != nil {
				fmt.Println(err)
			}
			fmt.Printf("Health: Pending requests %d, Last request completed at %d\n", health.PendingRequests, health.LastRequestCompletedAt)

		case EXIT:
			return

		default:
			fmt.Println("u shouldnt be here...")
		}
	}
}

func ManageAccountUI(userIban string) {
	items := []string{createAccount, printAccount, setAccountBalance, depositToAccount, withdrawFromAccount, suspendAccount, unsuspendAccount}

	items = append(items, UP)
	prompt := promptui.Select{
		Label: "Manage Account",
		Items: items,
	}
	_, it, _ := prompt.Run()

	switch it {

	case createAccount:
		pr := promptui.Prompt{Label: "Introduce the owner name of the new account", Default: "test_ownerName"}
		userName, _ := pr.Run()

		err := CreateAccount(userIban, userName)
		if err != nil {
			fmt.Println(err)
		}

	case printAccount:
		account, err := GetAccount(userIban)
		if err != nil {
			fmt.Println(err)
		}
		PrintAccount(account, true)

	case setAccountBalance:
		pr := promptui.Prompt{Label: "Introduce the new balance of the account", Default: "1"}
		balance, _ := pr.Run()

		err := SetAccountBalance(userIban, balance)
		if err != nil {
			fmt.Println(err)
		}

	case depositToAccount:
		pr := promptui.Prompt{Label: "Introduce the amount to deposit to the account", Default: "1"}
		amount, _ := pr.Run()

		err := DepositToAccount(userIban, amount)
		if err != nil {
			fmt.Println(err)
		}

	case withdrawFromAccount:
		pr := promptui.Prompt{Label: "Introduce the amount to withdraw from the account", Default: "1"}
		amount, _ := pr.Run()
		WithdrawFromAccount(userIban, amount)

	case suspendAccount:
		err := SuspendAccount(userIban)
		if err != nil {
			fmt.Println(err)
		}

	case unsuspendAccount:
		err := UnsuspendAccount(userIban)
		if err != nil {
			fmt.Println(err)
		}

	case UP:
		return

	default:
		fmt.Println("u shouldnt be here...")
	}
}
