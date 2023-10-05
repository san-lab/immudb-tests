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
var printAllBalances = "Print all balances"

func TopUI() {
	for {
		items := []string{printAllBalances, currentStateRoot, intraBankTx, interBankTx, vSet, vGet, health}
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
			CurrentStateRoot()

		case printAllBalances:
			PrintAllBalances()

		case "EXIT":
			return

		default:
			fmt.Println("u shouldnt be here...")
		}
	}
}
