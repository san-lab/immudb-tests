package main

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/codenotary/immudb/pkg/api/schema"
)

func VerifiedSet(key, value string) error {
	// write an entry
	// upon submission, the SDK validates proofs and updates the local state under the hood
	hdr, err := StateClient.VerifiedSet(context.Background(), []byte(key), []byte(value))
	if err != nil {
		fmt.Println(err)
		return err
	}
	_ = hdr
	//fmt.Printf("Sucessfully set a verified entry: ('%s', '%s') @ tx %d\n", []byte(key), []byte(value), hdr.Id)
	//fmt.Printf("Current state root is 0x%x\n", hdr.GetBlRoot())
	return nil
}

func VerifiedGet(key string) (*schema.Entry, error) {
	// read an entry
	// upon submission, the SDK validates proofs and updates the local state under the hood
	entry, err := StateClient.VerifiedGet(context.Background(), []byte(key))
	if err != nil {
		return nil, err
	}
	return entry, nil
}

func Health() (*schema.DatabaseHealthResponse, error) {
	health, err := StateClient.Health(context.Background())
	if err != nil {
		return nil, err
	}
	return health, nil
}

func CurrentStateRoot() ([]byte, uint64, error) {
	state, err := StateClient.CurrentState(context.Background())
	if err != nil {
		fmt.Println(err)
		return nil, 0, err
	}
	tx, err := StateClient.TxByID(context.Background(), state.GetTxId())
	if err != nil {
		fmt.Println(err)
		return nil, 0, err
	}
	root := tx.GetHeader().GetBlRoot()
	return root, state.GetTxId(), nil
}

func GetAllEntries() (*schema.Entries, error) {
	req := &schema.ScanRequest{Limit: 100} // 100 users...
	entries, err := StateClient.Scan(context.Background(), req)
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
	tx, err := StateClient.TxByID(context.Background(), uint64(id))
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return tx, nil
}

// ----- Msgs database methods -----

func StoreInMsgsDB(txmsg *MT103Message) (string, error) {
	value, err := json.Marshal(txmsg)
	if err != nil {
		return "", err
	}
	hash := sha256.Sum256(value)
	key := fmt.Sprintf("0x%x", hash[:])
	_, err = MsgsClient.VerifiedSet(context.Background(), []byte(key), value)
	return key, err
}

func VerifiedGetMsg(key string) (*schema.Entry, error) {
	entry, err := MsgsClient.VerifiedGet(context.Background(), []byte(key))
	return entry, err
}

func GetAllMsgsEntries() (*schema.Entries, error) {
	req := &schema.ScanRequest{Limit: 100} // 100 users...
	entries, err := MsgsClient.Scan(context.Background(), req)
	return entries, err
}
