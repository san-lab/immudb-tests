package main

import (
	"errors"
	"fmt"
	"math/big"

	"github.com/manifoldco/promptui"
	account "github.com/san-lab/immudb-tests/account"
	"github.com/san-lab/immudb-tests/bankinterop"
	"github.com/san-lab/immudb-tests/blockchainconnector"
	. "github.com/san-lab/immudb-tests/datastructs"
	sdk "github.com/san-lab/immudb-tests/immudbsdk"
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
const getAccountDigest = "Get the digest of the account"

const health = "Health"
const vSet = "VerifiedSet"
const vGet = "VerifiedGet"
const txById = "Get transaction by ID"

const seeMessagesDB = "See messages database...WIP"
const seeMessageByHash = "See a message by its hash...WIP"

const onChainOps = "Chain operations...for testing"
const getStateCheckByBlockNumber = "getStateCheckByBlockNumber"
const getStateCheckByIndex = "getStateCheckByIndex"
const getPendingSubmissions = "getPendingSubmissions"
const submitHash = "submitHash"
const submitPreimage = "submitPreimage"
const getBlockNumber = "getBlockNumber"
const getVersion = "getVersion"

func TopUI() {
	for {
		items := []string{bankInfo, findCounterpartBanks, seeCounterpartBanks, printAllAccounts, printAccount, currentStateRoot, intraBankTx, interBankTx,
			createAccount, manageAccount, txById, health, vSet, vGet, seeMessagesDB, seeMessageByHash, onChainOps}

		items = append(items, EXIT)
		prompt := promptui.Select{
			Label: THIS_BANK.Name + " - Actions",
			Items: items,
		}
		_, it, _ := prompt.Run()

		switch it {

		case bankInfo:
			bankinterop.PrintBankInfo()

		case findCounterpartBanks:
			err := bankinterop.FindCounterpartBanks()
			if err != nil {
				fmt.Println(err)
			}

		case seeCounterpartBanks:
			for k, v := range COUNTERPART_BANKS {
				if k != THIS_BANK.Name {
					fmt.Printf("| %s : %s\n", k, v)
				}
			}

		case interBankTx:
			pr := promptui.Prompt{Label: "Introduce the sender of the transaction", Default: "test_userFrom"}
			userFrom, _ := pr.Run()
			pr = promptui.Prompt{Label: "Introduce the amount", Default: "33"}
			amount, _ := pr.Run()
			pr = promptui.Prompt{Label: "Introduce the recipient of the transaction", Default: "test_userTo"}
			userTo, _ := pr.Run()

			// Grab the list of bank names, and have the user pick from it
			var CounterpartBankNames []string
			for k, _ := range COUNTERPART_BANKS {
				CounterpartBankNames = append(CounterpartBankNames, k)
			}
			prompt := promptui.Select{
				Label: "Select the bank of the recipient of the transaction",
				Items: append(CounterpartBankNames, UP),
			}
			_, bankTo, _ := prompt.Run()
			if bankTo != UP {
				err := bankinterop.InterBankTx(userFrom, amount, userTo, bankTo)
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

			err := bankinterop.IntraBankTx(userFrom, amount, userTo)
			if err != nil {
				fmt.Println(err)
				continue
			}

		case currentStateRoot:
			root, txId, err := sdk.CurrentStateRoot()
			if err != nil {
				fmt.Println(err)
				continue
			}
			fmt.Printf("Current state root: 0x%x (last tx id: %d)\n", root, txId)

		case printAllAccounts:
			accounts, _ := account.GetAllAccounts("")
			account.PrintAllAccounts(accounts)

		case printAccount:
			pr := promptui.Prompt{Label: "Introduce the IBAN of the new account", Default: "test_IBAN"}
			userIban, _ := pr.Run()

			acc, err := account.GetAccount(userIban)
			if err != nil {
				fmt.Println(err)
				continue
			}
			acc.PrintAccount(true)

		case txById:
			pr := promptui.Prompt{Label: "Introduce the ID of the transaction", Default: "0"}
			id, _ := pr.Run()

			tx, err := sdk.TxById(id)
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

			err := account.CreateAccount("", userIban, userName, "", "", 0, false, false)
			if err != nil {
				fmt.Println(err)
				continue
			}

		case manageAccount:
			pr := promptui.Prompt{Label: "Introduce the IBAN of the account to manage", Default: "test_IBAN"}
			userIban, _ := pr.Run()
			_, err := sdk.VerifiedGet(userIban)
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

			err := sdk.VerifiedSet(key, value)
			if err != nil {
				fmt.Println(err)
				continue
			}

		case vGet:
			pr := promptui.Prompt{Label: "Introduce the key", Default: "test_key"}
			key, _ := pr.Run()

			entry, err := sdk.VerifiedGet(key)
			if err != nil {
				fmt.Println(err)
				continue
			}
			fmt.Printf("Sucessfully got verified entry: ('%s', '%s') @ tx %d\n", entry.Key, entry.Value, entry.Tx)

		case health:
			health, err := sdk.Health()
			if err != nil {
				fmt.Println(err)
				continue
			}
			fmt.Printf("Health: Pending requests %d, Last request completed at %d\n", health.PendingRequests, health.LastRequestCompletedAt)

		case seeMessagesDB:
			msgentries, err := bankinterop.GetAllMessages()
			if err != nil {
				fmt.Println(err)
				continue
			}
			bankinterop.PrintAllMessages(msgentries)

		case seeMessageByHash:
			pr := promptui.Prompt{Label: "Introduce the hash of the message", Default: "test_hash"}
			hash, _ := pr.Run()
			message, err := bankinterop.GetMessage(hash)
			if err != nil {
				fmt.Println(err)
				continue
			}
			bankinterop.PrintMessage(message, true)

		case onChainOps:
			BlockchainOperationsUI()

		case EXIT:
			return

		default:
			fmt.Println("u shouldnt be here...")
		}
	}
}

func ManageAccountUI(userIban string) {
	for {
		items := []string{printAccount, setAccountBalance, depositToAccount, withdrawFromAccount, suspendAccount, unsuspendAccount, getAccountDigest}

		items = append(items, UP)
		prompt := promptui.Select{
			Label: "Manage Account",
			Items: items,
		}
		_, it, _ := prompt.Run()

		switch it {

		case printAccount:
			account, err := account.GetAccount(userIban)
			if err != nil {
				fmt.Println(err)
				continue
			}
			account.PrintAccount(true)

		case setAccountBalance:
			pr := promptui.Prompt{Label: "Introduce the new balance of the account", Default: "1"}
			balance, _ := pr.Run()

			err := account.SetAccountBalance(userIban, balance)
			if err != nil {
				fmt.Println(err)
				continue
			}

		case depositToAccount:
			pr := promptui.Prompt{Label: "Introduce the amount to deposit to the account", Default: "1"}
			amount, _ := pr.Run()

			err := account.DepositToAccount(userIban, amount)
			if err != nil {
				fmt.Println(err)
				continue
			}

		case withdrawFromAccount:
			pr := promptui.Prompt{Label: "Introduce the amount to withdraw from the account", Default: "1"}
			amount, _ := pr.Run()
			err := account.WithdrawFromAccount(userIban, amount)
			if err != nil {
				fmt.Println(err)
				continue
			}

		case suspendAccount:
			err := account.SuspendAccount(userIban)
			if err != nil {
				fmt.Println(err)
				continue
			}

		case unsuspendAccount:
			err := account.UnsuspendAccount(userIban)
			if err != nil {
				fmt.Println(err)
				continue
			}

		case getAccountDigest:
			accountDigest, err := account.GetAccountDigest(userIban)
			if err != nil {
				fmt.Println(err)
				continue
			}
			fmt.Println("Account digest: 0x", accountDigest)

		case UP:
			return

		default:
			fmt.Println("u shouldnt be here...")
		}
	}
}

func BlockchainOperationsUI() {
	for {
		items := []string{getStateCheckByBlockNumber, getStateCheckByIndex, getPendingSubmissions, submitHash, submitPreimage, getBlockNumber, getVersion}

		items = append(items, UP)
		prompt := promptui.Select{
			Label: "Blockchain Operations - Test",
			Items: items,
		}
		_, it, _ := prompt.Run()

		switch it {

		case getStateCheckByBlockNumber:
			originatorBank, _ := promptForBankName("Select the originator bank")
			originatorBank = COUNTERPART_BANKS[originatorBank]

			recipientBank, _ := promptForBankName("Select the recipient bank")
			recipientBank = COUNTERPART_BANKS[recipientBank]

			blockNumber, _ := promptForBigInt("Introduce the block number", "0")
			stateCheck, err := blockchainconnector.GetStateCheckByBlockNumber(originatorBank, recipientBank, blockNumber)
			if err != nil {
				fmt.Println(err)
				continue
			}
			blockchainconnector.PrintStateCheck(stateCheck)

		case getStateCheckByIndex:
			originatorBank, _ := promptForBankName("Select the originator bank")
			originatorBank = COUNTERPART_BANKS[originatorBank]

			recipientBank, _ := promptForBankName("Select the recipient bank")
			recipientBank = COUNTERPART_BANKS[recipientBank]

			index, _ := promptForBigInt("Introduce the index", "0")
			stateCheck, err := blockchainconnector.GetStateCheckByIndex(originatorBank, recipientBank, index)
			if err != nil {
				fmt.Println(err)
				continue
			}
			blockchainconnector.PrintStateCheck(stateCheck)

		case getPendingSubmissions:
			originatorBank, _ := promptForBankName("Select the originator bank")
			originatorBank = COUNTERPART_BANKS[originatorBank]
			pendingSubmissions, err := blockchainconnector.GetPendingSubmissions(originatorBank)
			if err != nil {
				fmt.Println(err)
				continue
			}
			fmt.Println("Pending Submissions:", pendingSubmissions)

		case submitHash:
			recipientBank, _ := promptForBankName("Select the recipient bank")
			recipientBank = COUNTERPART_BANKS[recipientBank]

			hash, _ := promptForString("Introduce the hash", "")

			err := blockchainconnector.SubmitHash(recipientBank, hash)
			if err != nil {
				fmt.Println(err)
				continue
			}

		case submitPreimage:
			originatorBank, _ := promptForBankName("Select the originator bank")
			originatorBank = COUNTERPART_BANKS[originatorBank]

			preimage, _ := promptForString("Introduce the preimage", "")

			blockNumber, _ := promptForBigInt("Introduce the block number", "0")
			err := blockchainconnector.SubmitPreimage(originatorBank, preimage, blockNumber)
			if err != nil {
				fmt.Println(err)
				continue
			}

		case getBlockNumber:
			blockNumber, err := blockchainconnector.GetBlockNumber()
			if err != nil {
				fmt.Println(err)
				continue
			}
			fmt.Println("Block Number:", blockNumber)

		case getVersion:
			version, err := blockchainconnector.Version()
			if err != nil {
				fmt.Println(err)
				continue
			}
			fmt.Println("Version:", version)

		case UP:
			return

		default:
			fmt.Println("u shouldnt be here...")
		}
	}
}

func promptForBankName(label string) (string, error) {
	var CounterpartBankNames []string
	for k, _ := range COUNTERPART_BANKS {
		CounterpartBankNames = append(CounterpartBankNames, k)
	}
	prompt := promptui.Select{
		Label: label,
		Items: append(CounterpartBankNames, UP),
	}
	_, bank, err := prompt.Run()
	return bank, err
}

func promptForBigInt(label, def string) (*big.Int, error) {
	pr := promptui.Prompt{Label: label, Default: def}
	numberString, _ := pr.Run()
	numberBigInt, ok := new(big.Int).SetString(numberString, 10)
	if !ok {
		return nil, errors.New("wrong block number format")
	}
	return numberBigInt, nil
}

func promptForString(label, def string) (string, error) {
	pr := promptui.Prompt{Label: label, Default: def}
	str, err := pr.Run()
	return str, err
}
