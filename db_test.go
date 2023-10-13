package main

import (
	"context"
	"fmt"
	"log"

	"testing"

	"github.com/codenotary/immudb/pkg/api/schema"
	"github.com/codenotary/immudb/pkg/client"
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
