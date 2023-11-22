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
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/san-lab/immudb-tests/blockchainconnector"
	"github.com/spf13/viper"
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

func TestGetPendingSubmissions(t *testing.T) {
	var NETWORK string
	//var CHAIN_ID string
	var VERIFIER_ADDRESS string
	//var PRIV_KEY_FILE string

	viper.SetConfigFile("config/config.env")
	viper.ReadInConfig()

	//ThisBankName := viper.GetString("BANK_NAME")
	ThisBankAddress := viper.GetString("BANK_ADDRESS")

	NETWORK = viper.GetString("NETWORK")
	//CHAIN_ID = viper.GetString("CHAIN_ID")
	VERIFIER_ADDRESS = viper.GetString("VERIFIER_ADDRESS")
	//PRIV_KEY_FILE = viper.GetString("PRIV_KEY_FILE")

	originatorBank := "0x6e7786c888Fe08E9360E830bC5806eca6186fB89"

	client, err := ethclient.Dial(NETWORK)
	if err != nil {
		fmt.Println("dial", err)
	}

	address := common.HexToAddress(VERIFIER_ADDRESS)
	instance, err := blockchainconnector.NewOnchainVerifier(address, client)
	if err != nil {
		fmt.Println("instance", err)
	}

	// Recipient must be ThisBank
	originatorBankAddress := common.HexToAddress(originatorBank)
	recipientBankAddress := common.HexToAddress(ThisBankAddress)

	pendingSubmissions, err := instance.GetPendingSubmissions(&bind.CallOpts{From: recipientBankAddress}, originatorBankAddress, recipientBankAddress)
	fmt.Println(err, pendingSubmissions)
	//return pendingSubmissions, err

}

func TestVersion(t *testing.T) {
	var NETWORK string
	//var CHAIN_ID string
	var VERIFIER_ADDRESS string
	//var PRIV_KEY_FILE string

	viper.SetConfigFile("config/config.env")
	viper.ReadInConfig()

	NETWORK = viper.GetString("NETWORK")
	//CHAIN_ID = viper.GetString("CHAIN_ID")
	VERIFIER_ADDRESS = viper.GetString("VERIFIER_ADDRESS")
	//PRIV_KEY_FILE = viper.GetString("PRIV_KEY_FILE")

	client, err := ethclient.Dial(NETWORK)
	if err != nil {
		fmt.Println("dial", err)
	}

	address := common.HexToAddress(VERIFIER_ADDRESS)
	instance, err := blockchainconnector.NewOnchainVerifier(address, client)
	if err != nil {
		fmt.Println("instance", err)
	}

	version, err := instance.Version(&bind.CallOpts{})
	fmt.Println(err, version)
	//return pendingSubmissions, err

}
