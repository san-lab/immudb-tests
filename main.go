package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"github.com/codenotary/immudb/pkg/api/schema"
	"github.com/codenotary/immudb/pkg/client"
	"github.com/spf13/viper"

	"github.com/san-lab/immudb-tests/blockchainconnector"
	. "github.com/san-lab/immudb-tests/datastructs"
	"github.com/san-lab/immudb-tests/transactions"
)

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
	c.CreateDatabaseV2(context.Background(), StateDB, &schema.DatabaseNullableSettings{})
	c2.CreateDatabaseV2(context.Background(), MsgsDB, &schema.DatabaseNullableSettings{})

	_, err = c.UseDatabase(context.Background(), &schema.Database{DatabaseName: StateDB})
	fmt.Println("Connection check State client", err, c.GetOptions().CurrentDatabase)

	_, err = c2.UseDatabase(context.Background(), &schema.Database{DatabaseName: MsgsDB})
	fmt.Println("Connection check Msgs client", err, c2.GetOptions().CurrentDatabase)

	StateClient = c
	MsgsClient = c2
}

func main() {
	/*
		ipFlag := flag.String("ip", "127.0.0.1", "ip to connect to the ImmuDB instance")
		portFlag := flag.Int("port", 3322, "port to connect to the ImmuDB instance")
		topicFlag := flag.String("net", "ImmuDBTopic", "name of the topic for the network")
		institutionName := flag.String("name", "SampleBank", "name of the financial institution")
	*/
	configFile := flag.String("config", "config.env", "path to config file")
	flag.Parse()
	viper.SetConfigFile(*configFile)
	viper.ReadInConfig()

	ThisBank.Name = viper.GetString("BANK_NAME")
	ThisBank.Address = viper.GetString("BANK_ADDRESS")
	CounterpartBanks["SampleBank"] = "0x1234"

	blockchainconnector.NETWORK = viper.GetString("NETWORK")
	blockchainconnector.CHAIN_ID = viper.GetString("CHAIN_ID")
	blockchainconnector.VERIFIER_ADDRESS = viper.GetString("VERIFIER_ADDRESS")
	blockchainconnector.PRIV_KEY_FILE = viper.GetString("PRIV_KEY_FILE")

	transactions.NET = viper.GetString("LIBP2P_TOPIC")
	transactions.LibP2PNode, _ = transactions.GetNode()

	initDB(viper.GetString("DB_IP"), viper.GetInt("DB_PORT"))

	// ensure connection is closed
	defer StateClient.CloseSession(context.Background())
	defer MsgsClient.CloseSession(context.Background())

	// PromptUI to select action
	TopUI()
}
