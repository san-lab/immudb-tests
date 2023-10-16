package main

import (
	"fmt"

	"github.com/manifoldco/promptui"
)

const UP = "UP"
const EXIT = "EXIT"

const bankInfo = "Show bank information"
const findCounterpartBanks = "Broadcast call to find other banks"
const seeCounterpartBanks = "See current list of counterpart banks"

const interBankTx = "Transfer to other bank client"
const intraBankTx = "Transfer to another client of the same bank"

const currentStateRoot = "Current state root"
const printAllAccounts = "Print all key-values stored"
const printAccount = "Print all values of an account"

const createAccount = "Create a new account"
const manageAccount = "Manage an account"
const setAccountBalance = "Set the balance of the account"
const depositToAccount = "Deposit the specified amount to the account"
const withdrawFromAccount = "Withdraw the specified amount from the account"
const suspendAccount = "Suspend the account"
const unsuspendAccount = "Unsuspend the account"

const health = "Health"
const vSet = "VerifiedSet"
const vGet = "VerifiedGet"
const txById = "Get transaction by ID"

const seeMessagesDB = "See messages database...WIP"
const seeMessageByHash = "See a message by its hash...WIP"

func TopUI() {
	for {
		items := []string{bankInfo, findCounterpartBanks, seeCounterpartBanks, printAllAccounts, printAccount, currentStateRoot, intraBankTx, interBankTx,
			createAccount, manageAccount, txById, health, vSet, vGet, seeMessagesDB, seeMessageByHash}

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

			// Substitute by discoveredCounterpartBanks
			prompt := promptui.Select{
				Label: "Select the bank of the recipient of the transaction",
				Items: append(CounterpartBanks, UP),
			}
			_, bankTo, _ := prompt.Run()
			if bankTo != UP {
				err := InterBankTx(userFrom, amount, userTo, bankTo)
				if err != nil {
					fmt.Println(err)
					continue
				}
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
				continue
			}

		case currentStateRoot:
			root, txId, err := CurrentStateRoot()
			if err != nil {
				fmt.Println(err)
				continue
			}
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
				continue
			}
			PrintAccount(account, true)

		case txById:
			pr := promptui.Prompt{Label: "Introduce the ID of the transaction", Default: "0"}
			id, _ := pr.Run()

			tx, err := TxById(id)
			if err != nil {
				fmt.Println(err)
				continue
			}
			entries := tx.GetEntries()
			for i, entry := range entries {
				fmt.Printf("Tx with id %s (%d): (%s : %s)\n", id, i, entry.Key, entry.Value)
			}

		case createAccount:
			pr := promptui.Prompt{Label: "Introduce the IBAN of the new account", Default: "test_IBAN"}
			userIban, _ := pr.Run()
			pr = promptui.Prompt{Label: "Introduce the owner name of the new account", Default: "test_ownerName"}
			userName, _ := pr.Run()

			err := CreateAccount(userIban, userName)
			if err != nil {
				fmt.Println(err)
				continue
			}

		case manageAccount:
			pr := promptui.Prompt{Label: "Introduce the IBAN of the account to manage", Default: "test_IBAN"}
			userIban, _ := pr.Run()
			_, err := VerifiedGet(userIban)
			if err != nil {
				fmt.Println(err)
				continue
			}
			ManageAccountUI(userIban)

		case vSet:
			pr := promptui.Prompt{Label: "Introduce the key", Default: "test_key"}
			key, _ := pr.Run()
			pr = promptui.Prompt{Label: "Introduce the value", Default: "test_value"}
			value, _ := pr.Run()

			err := VerifiedSet(key, value)
			if err != nil {
				fmt.Println(err)
				continue
			}

		case vGet:
			pr := promptui.Prompt{Label: "Introduce the key", Default: "test_key"}
			key, _ := pr.Run()

			entry, err := VerifiedGet(key)
			if err != nil {
				fmt.Println(err)
				continue
			}
			fmt.Printf("Sucessfully got verified entry: ('%s', '%s') @ tx %d\n", entry.Key, entry.Value, entry.Tx)

		case health:
			health, err := Health()
			if err != nil {
				fmt.Println(err)
				continue
			}
			fmt.Printf("Health: Pending requests %d, Last request completed at %d\n", health.PendingRequests, health.LastRequestCompletedAt)

		case seeMessagesDB:
			msgentries, err := GetAllMessages()
			if err != nil {
				fmt.Println(err)
				continue
			}
			PrintAllMessages(msgentries)

		case seeMessageByHash:
			pr := promptui.Prompt{Label: "Introduce the hash of the message", Default: "test_hash"}
			hash, _ := pr.Run()
			message, err := GetMessage(hash)
			if err != nil {
				fmt.Println(err)
				continue
			}
			PrintMessage(message, true)

		case EXIT:
			return

		default:
			fmt.Println("u shouldnt be here...")
		}
	}
}

func ManageAccountUI(userIban string) {
	for {
		items := []string{printAccount, setAccountBalance, depositToAccount, withdrawFromAccount, suspendAccount, unsuspendAccount}

		items = append(items, UP)
		prompt := promptui.Select{
			Label: "Manage Account",
			Items: items,
		}
		_, it, _ := prompt.Run()

		switch it {

		case printAccount:
			account, err := GetAccount(userIban)
			if err != nil {
				fmt.Println(err)
				continue
			}
			PrintAccount(account, true)

		case setAccountBalance:
			pr := promptui.Prompt{Label: "Introduce the new balance of the account", Default: "1"}
			balance, _ := pr.Run()

			err := SetAccountBalance(userIban, balance)
			if err != nil {
				fmt.Println(err)
				continue
			}

		case depositToAccount:
			pr := promptui.Prompt{Label: "Introduce the amount to deposit to the account", Default: "1"}
			amount, _ := pr.Run()

			err := DepositToAccount(userIban, amount)
			if err != nil {
				fmt.Println(err)
				continue
			}

		case withdrawFromAccount:
			pr := promptui.Prompt{Label: "Introduce the amount to withdraw from the account", Default: "1"}
			amount, _ := pr.Run()
			err := WithdrawFromAccount(userIban, amount)
			if err != nil {
				fmt.Println(err)
				continue
			}

		case suspendAccount:
			err := SuspendAccount(userIban)
			if err != nil {
				fmt.Println(err)
				continue
			}

		case unsuspendAccount:
			err := UnsuspendAccount(userIban)
			if err != nil {
				fmt.Println(err)
				continue
			}

		case UP:
			return

		default:
			fmt.Println("u shouldnt be here...")
		}
	}
}
