package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/codenotary/immudb/pkg/api/schema"
	"github.com/codenotary/immudb/pkg/client"
	account "github.com/san-lab/immudb-tests/account"
	"github.com/san-lab/immudb-tests/bankinterop"
	"github.com/san-lab/immudb-tests/blockchainconnector"
	. "github.com/san-lab/immudb-tests/datastructs"
	"github.com/spf13/viper"
)

var FIRST_BLOCK_NUMBER int

var DB_IP string
var DB_PORT int

var FIND_FREQUENCY int

func initDB() {
	// even though the server address and port are defaults, setting them as a reference
	opts := client.DefaultOptions().WithAddress(DB_IP).WithPort(DB_PORT)
	opts2 := client.DefaultOptions().WithAddress(DB_IP).WithPort(DB_PORT)

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

	blockchainconnector.NETWORK = viper.GetString("NETWORK")
	blockchainconnector.CHAIN_ID = viper.GetString("CHAIN_ID") // may not be needed anymore...
	blockchainconnector.VERIFIER_ADDRESS = viper.GetString("VERIFIER_ADDRESS")
	blockchainconnector.PRIV_KEY_FILE = viper.GetString("PRIV_KEY_FILE")

	bankinterop.NET = viper.GetString("LIBP2P_TOPIC")
	bankinterop.LIBP2P_NODE, _ = bankinterop.GetNode()

	API_PORT = viper.GetInt("API_PORT")

	DB_IP = viper.GetString("DB_IP")
	DB_PORT = viper.GetInt("DB_PORT")

	UPDATE_FREQUENCY = viper.GetInt("UPDATE_FREQUENCY")
	POLL_FREQUENCY = viper.GetInt("POLL_FREQUENCY")

	FIND_FREQUENCY = viper.GetInt("FIND_FREQUENCY")
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

func initNonce() {
	nonce, err := blockchainconnector.GetBlockchainNonce()
	if err != nil {
		fmt.Println(err)
	}
	blockchainconnector.NONCE = nonce
	fmt.Println("+ Initial nonce:", blockchainconnector.NONCE)
}

func periodicallyFindCounterpartBanks(done chan (bool)) {
	// Try at least once on initialization
	time.Sleep(2 * time.Second)
	err := bankinterop.FindCounterpartBanks()
	if err != nil {
		fmt.Println(err)
	}

	if FIND_FREQUENCY == 0 {
		return
	}

	ticker := time.NewTicker(time.Duration(FIND_FREQUENCY) * time.Second)
	for {
		select {
		case <-done:
			return

		case t := <-ticker.C:
			// fmt.Println("debug find counterparties", t)
			_ = t
			err := bankinterop.FindCounterpartBanks()
			if err != nil {
				fmt.Println(err)
			}
		}
	}
}

func PrintBankInfo() {
	fmt.Println("| Bank Name:", THIS_BANK.Name)
	fmt.Println("| Bank Address:", THIS_BANK.Address)
	fmt.Printf("| ImmuDB instance running on %s:%d\n", STATE_CLIENT.GetOptions().Address, STATE_CLIENT.GetOptions().Port)
	fmt.Println("| API server running on port:", API_PORT)
	fmt.Printf("| Connected to blockchain network %s, with chain ID: %s\n", blockchainconnector.NETWORK, blockchainconnector.CHAIN_ID)
	fmt.Println("| Verifier contract address:", blockchainconnector.VERIFIER_ADDRESS)
}
