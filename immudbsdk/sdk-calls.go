package immudbsdk

import (
	"context"
	"fmt"
	"strconv"

	"github.com/codenotary/immudb/pkg/api/schema"
	. "github.com/san-lab/immudb-tests/datastructs"
)

func VerifiedSet(key, value string) error {
	// write an entry
	// upon submission, the SDK validates proofs and updates the local state under the hood
	_, err := STATE_CLIENT.VerifiedSet(context.Background(), []byte(key), []byte(value))
	if err != nil {
		fmt.Println(err)
	}
	return err
}

func VerifiedGet(key string) (*schema.Entry, error) {
	// read an entry
	// upon submission, the SDK validates proofs and updates the local state under the hood
	entry, err := STATE_CLIENT.VerifiedGet(context.Background(), []byte(key))
	return entry, err
}

func Health() (*schema.DatabaseHealthResponse, error) {
	health, err := STATE_CLIENT.Health(context.Background())
	return health, err
}

func CurrentStateRoot() ([]byte, uint64, error) {
	state, err := STATE_CLIENT.CurrentState(context.Background())
	if err != nil {
		fmt.Println(err)
		return nil, 0, err
	}
	tx, err := STATE_CLIENT.TxByID(context.Background(), state.GetTxId())
	if err != nil {
		fmt.Println(err)
		return nil, 0, err
	}
	root := tx.GetHeader().GetBlRoot()
	return root, state.GetTxId(), nil
}

func GetAllEntries() (*schema.Entries, error) {
	req := &schema.ScanRequest{Limit: 100} // 100 users...
	entries, err := STATE_CLIENT.Scan(context.Background(), req)
	return entries, err
}

func TxById(idString string) (*schema.Tx, error) {
	id, err := strconv.Atoi(idString)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	tx, err := STATE_CLIENT.TxByID(context.Background(), uint64(id))
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return tx, nil
}

// ----- Msgs database methods -----

func VerifiedSetMsg(key, value string) error {
	_, err := MSGS_CLIENT.VerifiedSet(context.Background(), []byte(key), []byte(value))
	if err != nil {
		fmt.Println(err)
	}
	return err
}

func VerifiedGetMsg(key string) (*schema.Entry, error) {
	entry, err := MSGS_CLIENT.VerifiedGet(context.Background(), []byte(key))
	return entry, err
}

func GetAllMsgsEntries() (*schema.Entries, error) {
	req := &schema.ScanRequest{Limit: 100} // 100 users...
	entries, err := MSGS_CLIENT.Scan(context.Background(), req)
	return entries, err
}
