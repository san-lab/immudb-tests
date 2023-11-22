package main

import (
	"context"
	"encoding/hex"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/codenotary/immudb/pkg/api/schema"
	"github.com/codenotary/immudb/pkg/client"
	"github.com/spf13/viper"

	account "github.com/san-lab/immudb-tests/account"
	"github.com/san-lab/immudb-tests/bankinterop"
	"github.com/san-lab/immudb-tests/blockchainconnector"
	. "github.com/san-lab/immudb-tests/datastructs"
	"github.com/wealdtech/go-merkletree/keccak256"
)

const STATE_DB = "defaultdb"
const MSGS_DB = "msgdb"

const UPDATE_FREQUENCY = 60
const POLL_FREQUENCY = 30

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
	fmt.Println("Connection check State client", err, c.GetOptions().CurrentDatabase)

	_, err = c2.UseDatabase(context.Background(), &schema.Database{DatabaseName: MSGS_DB})
	fmt.Println("Connection check Msgs client", err, c2.GetOptions().CurrentDatabase)

	STATE_CLIENT = c
	MSGS_CLIENT = c2
}

func main() {
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

	// Initialize DBs
	initDB(viper.GetString("DB_IP"), viper.GetInt("DB_PORT"))

	// Ensure connection is closed
	defer STATE_CLIENT.CloseSession(context.Background())
	defer MSGS_CLIENT.CloseSession(context.Background())

	// Go routines that interact with blockchain
	ticker := time.NewTicker(UPDATE_FREQUENCY * time.Second)
	done := make(chan bool)
	go periodicallySubmitHash(done, ticker)

	ticker2 := time.NewTicker(POLL_FREQUENCY * time.Second)
	done2 := make(chan bool)
	go preiodicallyPollAndSubmitPreImage(done2, ticker2)

	// PromptUI to select action
	TopUI()

	//time.Sleep(1600 * time.Millisecond)
	//ticker.Stop()
	//done <- true
}

// TODO: handle more than one pending submission properly
// TODO: check for changes in the state before updating
func periodicallySubmitHash(done chan bool, ticker *time.Ticker) {
	for {
		select {
		case <-done:
			return

		case t := <-ticker.C:
			fmt.Println("--- Tick at", t)
			mirrorAccounts, _ := account.GetAllAccounts("mirror")
			for _, mirrorAccount := range mirrorAccounts {
				recipientBank := COUNTERPART_BANKS[mirrorAccount.CABank]
				digest, _ := mirrorAccount.GetDigest()
				digestBytes, err := hex.DecodeString(digest)
				if err != nil {
					fmt.Println(err)
					continue
				}
				hashBytes := keccak256.New().Hash(digestBytes)
				hash := fmt.Sprintf("%x", hashBytes)
				fmt.Println("Submitting hash for", mirrorAccount.Iban)
				fmt.Println("Digest", digest)
				fmt.Println("Hash", hash)
				err = blockchainconnector.SubmitHash(recipientBank, hash)
				if err != nil {
					fmt.Println(err)
					continue
				}
				fmt.Println("done!")
			}
		}
	}
}

func preiodicallyPollAndSubmitPreImage(done chan bool, ticker *time.Ticker) {
	for {
		select {
		case <-done:
			return

		case t := <-ticker.C:
			fmt.Println("--- Tick at", t)
			CAAccounts, _ := account.GetAllAccounts("ca")
			for _, CAAccount := range CAAccounts {
				fmt.Println("- Checking CA ", CAAccount.Iban)
				originatorBank := COUNTERPART_BANKS[CAAccount.CABank]
				pending, err := blockchainconnector.GetPendingSubmissions(originatorBank)
				if err != nil {
					fmt.Println(err)
					continue
				}
				if len(pending) == 0 {
					fmt.Println("No pending submissions for this CA")
					continue
				}
				fmt.Println("Pending submissions:", pending)
				digest, _ := CAAccount.GetDigest()
				fmt.Println("Submitting preimage for", CAAccount.Iban)
				fmt.Println("Digest", digest)
				err = blockchainconnector.SubmitPreimage(originatorBank, digest, pending[len(pending)-1])
				if err != nil {
					fmt.Println(err)
					continue
				}
				fmt.Println("done!")
			}
		}
	}
}
