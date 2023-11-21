package main

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"os"
	"time"

	"testing"

	"github.com/codenotary/immudb/pkg/api/schema"
	"github.com/codenotary/immudb/pkg/client"
	"github.com/ethereum/go-ethereum/crypto"
)

var ip = "127.0.0.1"
var port = 3322

func TestDBConnection(t *testing.T) {
	// even though the server address and port are defaults, setting them as a reference
	opts := client.DefaultOptions().WithAddress(ip).WithPort(port)

	c := client.NewClient().WithOptions(opts)

	// connect with immudb server (user, password, database)
	err := c.OpenSession(context.Background(), []byte("immudb"), []byte("immudb"), "defaultdb")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("db1:", c.GetOptions().CurrentDatabase)

	fmt.Println("--------")

	opts2 := client.DefaultOptions().WithAddress(ip).WithPort(port)

	c2 := client.NewClient().WithOptions(opts2)

	newDB, err := c.CreateDatabaseV2(context.Background(), "msgsdb", &schema.DatabaseNullableSettings{})
	fmt.Println("err:", err)
	fmt.Println(newDB.GetName(), newDB.GetAlreadyExisted())

	// connect with immudb server (user, password, database)
	err = c2.OpenSession(context.Background(), []byte("immudb"), []byte("immudb"), "msgsssdb")
	if err != nil {
		fmt.Println(err)
	}

	_, err = c.UseDatabase(context.Background(), &schema.Database{DatabaseName: "msgssssdb"})
	fmt.Println("err:", err)
	fmt.Println("db2:", c.GetOptions().CurrentDatabase)

}

func TestTimestamp(t *testing.T) {
	timestamp := time.Now()
	fmt.Println(timestamp.String())
}

func TestReadKey(t *testing.T) {
	privKeyBytes, err := os.ReadFile("config/priv_key.txt")
	fmt.Println(err, string(privKeyBytes))
	/*
		privKey := new(datastructs.PrivateKey)
		err = json.Unmarshal(privKeyBytes, privKey)
		if err != nil {
			return nil, err
		}
	*/

	privateKey, err := crypto.HexToECDSA(string(privKeyBytes))
	fmt.Println(err, privateKey)
	publicKey := privateKey.Public()
	publicKeyECDSA, _ := publicKey.(*ecdsa.PublicKey)
	fmt.Println(err, publicKeyECDSA)
}
