package main

import (
	"context"
	"fmt"
	"strconv"

	"github.com/codenotary/immudb/pkg/api/schema"
)

func VerifiedSet(key, value string) error {
	// write an entry
	// upon submission, the SDK validates proofs and updates the local state under the hood
	hdr, err := Client.VerifiedSet(context.Background(), []byte(key), []byte(value))
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Printf("Sucessfully set a verified entry: ('%s', '%s') @ tx %d\n", []byte(key), []byte(value), hdr.Id)
	fmt.Printf("Current state root is 0x%x\n", hdr.GetBlRoot())
	return nil
}

func VerifiedGet(key string) ([]byte, error) {
	// read an entry
	// upon submission, the SDK validates proofs and updates the local state under the hood
	entry, err := Client.VerifiedGet(context.Background(), []byte(key))
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	fmt.Printf("Sucessfully got verified entry: ('%s', '%s') @ tx %d\n", entry.Key, entry.Value, entry.Tx)
	return entry.Value, nil
}

func Health() error {
	health, err := Client.Health(context.Background())
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Printf("Health: Pending requests %d, Last request completed at %d\n", health.PendingRequests, health.LastRequestCompletedAt)
	return nil
}

func CurrentStateRoot() ([]byte, uint64, error) {
	state, err := Client.CurrentState(context.Background())
	if err != nil {
		fmt.Println(err)
		return nil, 0, err
	}
	tx, err := Client.TxByID(context.Background(), state.GetTxId())
	if err != nil {
		fmt.Println(err)
		return nil, 0, err
	}
	root := tx.GetHeader().GetBlRoot()
	return root, state.GetTxId(), nil
}

func GetAllAccounts() (*schema.Entries, error) {
	req := &schema.ScanRequest{Limit: 10} // only 10 users
	entries, err := Client.Scan(context.Background(), req)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return entries, nil
}

func TxById(idString string) (*schema.Tx, error) {
	id, err := strconv.Atoi(idString)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	tx, err := Client.TxByID(context.Background(), uint64(id))
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return tx, nil
}
