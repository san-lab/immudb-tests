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

const printAllAccounts = "Print all accounts stored"
const printAccount = "Print all values of an account"

const interBankTx = "Transfer to a client of a different bank"
const intraBankTx = "Transfer to a client of the same bank"
const refillCA = "Refill our correspondent account at other bank"

const currentStateRoot = "Current state root"

const createAccount = "Create a new account"
const manageAccount = "Manage an account"
const setAccountBalance = "Set the balance of the account"
const depositToAccount = "Deposit the specified amount to the account"
const withdrawFromAccount = "Withdraw the specified amount from the account"
const suspendAccount = "Suspend the account"
const unsuspendAccount = "Unsuspend the account"
const getAccountDigest = "Get the digest of the account"

const DBOps = "Low level database operations"
const health = "Health"
const vSet = "VerifiedSet"
const vGet = "VerifiedGet"
const txById = "Get transaction by ID"

const seeMessagesDB = "See messages database...WIP"
const seeMessageByHash = "See a message by its hash...WIP"

const onChainOps = "On-chain operations"
const getStateCheckByBlockNumber = "getStateCheckByBlockNumber"
const getStateCheckByIndex = "getStateCheckByIndex"
const getPendingSubmissions = "getPendingSubmissions"
const submitHash = "submitHash"
const submitPreimage = "submitPreimage"
const getBlockNumber = "getBlockNumber"
const getVersion = "getVersion"

func TopUI() {
	for {
		items := []string{bankInfo, seeCounterpartBanks, printAllAccounts, printAccount,
			intraBankTx, interBankTx, refillCA, createAccount, manageAccount, seeMessagesDB,
			seeMessageByHash, DBOps, currentStateRoot, onChainOps, findCounterpartBanks}

		items = append(items, EXIT)
		prompt := promptui.Select{
			Label: THIS_BANK.Name + " - Actions",
			Items: items,
		}
		_, it, _ := prompt.Run()

		switch it {

		case bankInfo:
			PrintBankInfo()

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
			userFrom, _ := promptForString("Introduce the sender of the transaction", "test_userFrom")
			amount, _ := promptForString("Introduce the amount", "33")
			userTo, _ := promptForString("Introduce the recipient of the transaction", "test_userTo")
			bankTo, _ := promptForBankName(false, "Select the bank of the recipient of the transaction")

			if bankTo != UP {
				err := bankinterop.RequestInterBankTx(userFrom, amount, userTo, bankTo)
				if err != nil {
					fmt.Println(err)
					continue
				}
			}

		case intraBankTx:
			userFrom, _ := promptForString("Introduce the sender of the transaction", "test_userFrom")
			amount, _ := promptForString("Introduce the amount", "33")
			userTo, _ := promptForString("Introduce the recipient of the transaction", "test_userTo")

			err := bankinterop.IntraBankTx(userFrom, amount, userTo)
			if err != nil {
				fmt.Println(err)
				continue
			}

		case refillCA:
			amount, _ := promptForString("Introduce the amount to deposit", "100")
			bankTo, _ := promptForBankName(false, "Select the bank recipient of the transaction")

			if bankTo != UP {
				err := bankinterop.RequestRefillCA(amount, bankTo)
				if err != nil {
					fmt.Println(err)
					continue
				}
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
			userIban, _ := promptForString("Introduce the IBAN of the new account", "test_IBAN")

			acc, err := account.GetAccount(userIban)
			if err != nil {
				fmt.Println(err)
				continue
			}
			acc.PrintAccount(true)

		case createAccount:
			userIban, _ := promptForString("Introduce the IBAN of the new account", "test_IBAN")
			userName, _ := promptForString("Introduce the owner name of the new account", "test_ownerName")

			err := account.CreateAccount("", userIban, userName, "", "", 0, false, false, false)
			if err != nil {
				fmt.Println(err)
				continue
			}

		case manageAccount:
			userIban, _ := promptForString("Introduce the IBAN of the account to manage", "test_IBAN")

			_, err := sdk.VerifiedGet(userIban)
			if err != nil {
				fmt.Println(err)
				continue
			}
			ManageAccountUI(userIban)

		case seeMessagesDB:
			msgentries, err := bankinterop.GetAllMessages()
			if err != nil {
				fmt.Println(err)
				continue
			}
			bankinterop.PrintAllMessages(msgentries)

		case seeMessageByHash:
			hash, _ := promptForString("Introduce the hash of the message", "test_hash")

			message, err := bankinterop.GetMessage(hash)
			if err != nil {
				fmt.Println(err)
				continue
			}
			bankinterop.PrintMessage(message, true)

		case DBOps:
			DatabaseOperationsUI()

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
			balance, _ := promptForString("Introduce the new balance of the account", "1")

			err := account.SetAccountBalance(userIban, balance)
			if err != nil {
				fmt.Println(err)
				continue
			}

		case depositToAccount:
			amount, _ := promptForString("Introduce the amount to deposit to the account", "1")

			err := account.DepositToAccount(userIban, amount)
			if err != nil {
				fmt.Println(err)
				continue
			}

		case withdrawFromAccount:
			amount, _ := promptForString("Introduce the amount to withdraw from the account", "1")

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
			fmt.Printf("Account digest: 0x%s\n", accountDigest)

		case UP:
			return

		default:
			fmt.Println("u shouldnt be here...")
		}
	}
}

func BlockchainOperationsUI() {
	for {
		items := []string{getStateCheckByBlockNumber, getStateCheckByIndex, getPendingSubmissions, submitHash,
			submitPreimage, getBlockNumber, getVersion}

		items = append(items, UP)
		prompt := promptui.Select{
			Label: "Blockchain Operations",
			Items: items,
		}
		_, it, _ := prompt.Run()

		switch it {

		case getStateCheckByBlockNumber:
			originatorBank, _ := promptForBankName(true, "Select the originator bank")
			originatorBank = COUNTERPART_BANKS[originatorBank]

			recipientBank, _ := promptForBankName(true, "Select the recipient bank")
			recipientBank = COUNTERPART_BANKS[recipientBank]

			blockNumber, _ := promptForBigInt("Introduce the block number", "0")
			stateCheck, err := blockchainconnector.GetStateCheckByBlockNumber(originatorBank, recipientBank, blockNumber)
			if err != nil {
				fmt.Println(err)
				continue
			}
			blockchainconnector.PrintStateCheck(stateCheck)

		case getStateCheckByIndex:
			originatorBank, _ := promptForBankName(true, "Select the originator bank")
			originatorBank = COUNTERPART_BANKS[originatorBank]

			recipientBank, _ := promptForBankName(true, "Select the recipient bank")
			recipientBank = COUNTERPART_BANKS[recipientBank]

			index, _ := promptForBigInt("Introduce the index", "0")
			stateCheck, err := blockchainconnector.GetStateCheckByIndex(originatorBank, recipientBank, index)
			if err != nil {
				fmt.Println(err)
				continue
			}
			blockchainconnector.PrintStateCheck(stateCheck)

		case getPendingSubmissions:
			originatorBank, _ := promptForBankName(true, "Select the originator bank")
			originatorBank = COUNTERPART_BANKS[originatorBank]
			pendingSubmissions, err := blockchainconnector.GetPendingSubmissions(originatorBank)
			if err != nil {
				fmt.Println(err)
				continue
			}
			fmt.Println("Pending Submissions:", pendingSubmissions)

		case submitHash:
			recipientBank, _ := promptForBankName(true, "Select the recipient bank")
			recipientBank = COUNTERPART_BANKS[recipientBank]

			hash, _ := promptForString("Introduce the hash", "")

			err := blockchainconnector.SubmitHash(recipientBank, hash)
			if err != nil {
				fmt.Println(err)
				continue
			}

		case submitPreimage:
			originatorBank, _ := promptForBankName(true, "Select the originator bank")
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

func DatabaseOperationsUI() {
	for {
		items := []string{health, vGet, vSet, txById}

		items = append(items, UP)
		prompt := promptui.Select{
			Label: "Database Operations",
			Items: items,
		}
		_, it, _ := prompt.Run()

		switch it {

		case vSet:
			key, _ := promptForString("Introduce the key", "test_key")
			value, _ := promptForString("Introduce the value", "test_value")

			err := sdk.VerifiedSet(key, value)
			if err != nil {
				fmt.Println(err)
				continue
			}

		case vGet:
			key, _ := promptForString("Introduce the key", "test_key")

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

		case txById:
			id, _ := promptForString("Introduce the ID of the transaction", "0")

			tx, err := sdk.TxById(id)
			if err != nil {
				fmt.Println(err)
				continue
			}
			entries := tx.GetEntries()
			for i, entry := range entries {
				fmt.Printf("Tx with id %s (%d): (%s : %s)\n", id, i, entry.Key, entry.Value)
			}

		case UP:
			return

		default:
			fmt.Println("u shouldnt be here...")
		}
	}
}

func promptForBankName(includeOurselves bool, label string) (string, error) {
	var CounterpartBankNames []string
	for k, _ := range COUNTERPART_BANKS {
		if includeOurselves || k != THIS_BANK.Name {
			CounterpartBankNames = append(CounterpartBankNames, k)
		}
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
