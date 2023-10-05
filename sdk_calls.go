package main

import (
	"context"
	"fmt"
	"strconv"

	"github.com/codenotary/immudb/pkg/api/schema"
)

func VerifiedSet(key, value string) {
	// write an entry
	// upon submission, the SDK validates proofs and updates the local state under the hood
	hdr, err := Client.VerifiedSet(context.Background(), []byte(key), []byte(value))
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("Sucessfully set a verified entry: ('%s', '%s') @ tx %d\n", []byte(key), []byte(value), hdr.Id)
	fmt.Printf("Current state root is 0x%x\n", hdr.GetBlRoot())
}

func VerifiedGet(key string) []byte {
	// read an entry
	// upon submission, the SDK validates proofs and updates the local state under the hood
	entry, err := Client.VerifiedGet(context.Background(), []byte(key))
	if err != nil {
		fmt.Println(err)
		return nil
	}
	fmt.Printf("Sucessfully got verified entry: ('%s', '%s') @ tx %d\n", entry.Key, entry.Value, entry.Tx)
	return entry.Value
}

func Health() {
	health, err := Client.Health(context.Background())
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("Health: Pending requests %d, Last request completed at %d\n", health.PendingRequests, health.LastRequestCompletedAt)
}

func CurrentStateRoot() {
	state, err := Client.CurrentState(context.Background())
	if err != nil {
		fmt.Println(err)
		return
	}
	tx, err := Client.TxByID(context.Background(), state.GetTxId())
	if err != nil {
		fmt.Println(err)
		return
	}
	root := tx.GetHeader().GetBlRoot()
	fmt.Printf("Current state root: 0x%x (last tx id: %d)", root, state.GetTxId())
}

func InterBankTx(userFrom, amount, userTo string) {
	value := VerifiedGet(userFrom)

	valueInt, _ := strconv.Atoi(string(value))
	amountInt, _ := strconv.Atoi(string(amount))
	newBalanceInt := valueInt - amountInt
	if newBalanceInt <= 0 {
		fmt.Println("Not enough balance to perform the transaction")
		return
	}

	newBalance := strconv.Itoa(newBalanceInt)
	VerifiedSet(userFrom, newBalance)

	// Sent event to the topic
	txmsg := &TxMessage{userFrom, amount, userTo}
	node.SendMsg(txmsg)
}

func ProcessInterBankTx(amount, userTo string) {
	value := VerifiedGet(userTo)

	valueInt, _ := strconv.Atoi(string(value))
	amountInt, _ := strconv.Atoi(string(amount))
	newBalanceInt := valueInt + amountInt

	newBalance := strconv.Itoa(newBalanceInt)
	VerifiedSet(userTo, newBalance)
}

func IntraBankTx(userFrom, amount, userTo string) {
	valueFrom := VerifiedGet(userFrom)

	valueFromInt, _ := strconv.Atoi(string(valueFrom))
	amountInt, _ := strconv.Atoi(string(amount))
	newBalanceFromInt := valueFromInt - amountInt
	if newBalanceFromInt <= 0 {
		fmt.Println("Not enough balance to perform the transaction")
		return
	}
	newBalanceFrom := strconv.Itoa(newBalanceFromInt)
	VerifiedSet(userFrom, newBalanceFrom)

	valueTo := VerifiedGet(userTo)
	valueToInt, _ := strconv.Atoi(string(valueTo))
	newBalanceToInt := valueToInt + amountInt
	newBalanceTo := strconv.Itoa(newBalanceToInt)
	VerifiedSet(userTo, newBalanceTo)

}

func PrintAllBalances() {
	req := &schema.ScanRequest{Limit: 10}
	entries, err := Client.Scan(context.Background(), req)
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, entry := range entries.Entries {
		fmt.Printf("(%s : %s)\n", entry.Key, entry.Value)
	}
}
