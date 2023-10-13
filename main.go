package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"github.com/codenotary/immudb/pkg/api/schema"
	"github.com/codenotary/immudb/pkg/client"
)

var InstitutionName string
var CounterpartBanks []string

var LibP2PNode *Node
var Client client.ImmuClient

const StateDB = "defaultdb"
const MsgsDB = "msgdb"

func initDB(ip string, port int) {
	// even though the server address and port are defaults, setting them as a reference
	opts := client.DefaultOptions().WithAddress(ip).WithPort(port)
	c := client.NewClient().WithOptions(opts)

	// connect with immudb server (user, password, database)
	err := c.OpenSession(context.Background(), []byte("immudb"), []byte("immudb"), "defaultdb")
	if err != nil {
		log.Fatal(err)
	}

	// Tries to create both databases in case they dont exist (not thread safe without mutex or second client)
	c.CreateDatabaseV2(context.Background(), StateDB, &schema.DatabaseNullableSettings{})
	c.CreateDatabaseV2(context.Background(), MsgsDB, &schema.DatabaseNullableSettings{})

	_, err = c.UseDatabase(context.Background(), &schema.Database{DatabaseName: MsgsDB})
	fmt.Println("Connection check", err, c.GetOptions().CurrentDatabase)
	_, err = c.UseDatabase(context.Background(), &schema.Database{DatabaseName: StateDB})
	fmt.Println("Connection check", err, c.GetOptions().CurrentDatabase)

	Client = c
}

func main() {
	ipFlag := flag.String("ip", "127.0.0.1", "ip to connect to the ImmuDB instance")
	portFlag := flag.Int("port", 3322, "port to connect to the ImmuDB instance")
	topicFlag := flag.String("net", "ImmuDBTopic", "name of the topic for the network")
	institutionName := flag.String("name", "SampleBank", "name of the financial institution")
	flag.Parse()

	InstitutionName = *institutionName
	CounterpartBanks = []string{"GreenBank", "RedBank", "BlueBank"}

	NET = *topicFlag
	LibP2PNode, _ := GetNode()
	LibP2PNode.GetNodeID()

	initDB(*ipFlag, *portFlag)

	// ensure connection is closed
	defer Client.CloseSession(context.Background())

	// PromptUI to select action
	TopUI()
}

// TODO

// connect to blockchain
// store msg hash, banks and amount

//
// Important to add a nonce to each transaction at least for the MT messages, otherwise 2 equal messages wont show up on the database....
//
