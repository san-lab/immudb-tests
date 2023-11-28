package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"github.com/codenotary/immudb/pkg/api/schema"
	"github.com/codenotary/immudb/pkg/client"
	account "github.com/san-lab/immudb-tests/account"
	"github.com/san-lab/immudb-tests/bankinterop"
	"github.com/san-lab/immudb-tests/blockchainconnector"
	. "github.com/san-lab/immudb-tests/datastructs"
	"github.com/spf13/viper"
)

var FIRST_BLOCK_NUMBER int

func initDB(ip string, port int) {
	// even though the server address and port are defaults, setting them as a reference
	opts := client.DefaultOptions().WithAddress(ip).WithPort(port)
	opts2 := client.DefaultOptions().WithAddress(ip).WithPort(port)

	c := client.NewClient().WithOptions(opts)
	c2 := client.NewClient().WithOptions(opts2)

	// connect with immudb server (user, password, database)
	err := c.OpenSession(context.Background(), []byte("immudb"), []byte("immudb"), "defaultdb")
	if err != nil {
		log.Fatal(err)
	}
	err = c2.OpenSession(context.Background(), []byte("immudb"), []byte("immudb"), "defaultdb")
	if err != nil {
		log.Fatal(err)
	}

	// Tries to create both databases in case they dont exist (not thread safe without mutex or second client)
	c.CreateDatabaseV2(context.Background(), STATE_DB, &schema.DatabaseNullableSettings{})
	c2.CreateDatabaseV2(context.Background(), MSGS_DB, &schema.DatabaseNullableSettings{})

	_, err = c.UseDatabase(context.Background(), &schema.Database{DatabaseName: STATE_DB})
	fmt.Println("+ Connection check State client:", err, c.GetOptions().CurrentDatabase)

	_, err = c2.UseDatabase(context.Background(), &schema.Database{DatabaseName: MSGS_DB})
	fmt.Println("+ Connection check Msgs client:", err, c2.GetOptions().CurrentDatabase)

	STATE_CLIENT = c
	MSGS_CLIENT = c2
}

func initConfigParams() {
	// Parsing all parameters
	configFile := flag.String("config", "config/config.env", "path to config file")
	flag.Parse()

	viper.SetConfigFile(*configFile)
	viper.ReadInConfig()

	THIS_BANK.Name = viper.GetString("BANK_NAME")
	THIS_BANK.Address = viper.GetString("BANK_ADDRESS")
	COUNTERPART_BANKS[THIS_BANK.Name] = THIS_BANK.Address
	COUNTERPART_BANKS["SampleBank"] = "0x1234"

	blockchainconnector.NETWORK = viper.GetString("NETWORK")
	blockchainconnector.CHAIN_ID = viper.GetString("CHAIN_ID")
	blockchainconnector.VERIFIER_ADDRESS = viper.GetString("VERIFIER_ADDRESS")
	blockchainconnector.PRIV_KEY_FILE = viper.GetString("PRIV_KEY_FILE")

	bankinterop.NET = viper.GetString("LIBP2P_TOPIC")
	bankinterop.LIBP2P_NODE, _ = bankinterop.GetNode()

	API_PORT = viper.GetInt("API_PORT")

	DB_IP = viper.GetString("DB_IP")
	DB_PORT = viper.GetInt("DB_PORT")
}

// Initialize digest history
func initDigestHistory() {
	CAAccounts, _ := account.GetAllAccounts("ca")
	for _, CAAccount := range CAAccounts {
		bankinterop.DigestHistory[CAAccount.CABank] = make(map[int]string)
		blockNumber, err := blockchainconnector.GetBlockNumber()
		if err != nil {
			fmt.Println(err)
		}
		digest, err := account.GetAccountDigest(CAAccount.Iban)
		if err != nil {
			fmt.Println(err)
		}
		bankinterop.DigestHistory[CAAccount.CABank][blockNumber] = digest
		fmt.Println("+ Initial digest history:", CAAccount.CABank, bankinterop.DigestHistory[CAAccount.CABank])
		FIRST_BLOCK_NUMBER = blockNumber
	}
}

func PrintBankInfo() {
	fmt.Println("| Bank Name:", THIS_BANK.Name)
	fmt.Println("| Bank Address:", THIS_BANK.Address)
	fmt.Println("| ImmuDB instance running on IP:", STATE_CLIENT.GetOptions().Address)
	fmt.Println("| ImmuDB instance running on port:", STATE_CLIENT.GetOptions().Port)
	fmt.Println("| ...")
}
