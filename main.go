package main

import (
	"context"
	"time"

	. "github.com/san-lab/immudb-tests/datastructs"
)

const STATE_DB = "defaultdb"
const MSGS_DB = "msgdb"

var DB_IP string
var DB_PORT int

func main() {
	// Parse all config parameters
	initConfigParams()

	// Initialize DBs
	initDB(DB_IP, DB_PORT)

	// Initialize digest history of onboarded CA banks
	initDigestHistory()

	// Ensure connection is closed
	defer STATE_CLIENT.CloseSession(context.Background())
	defer MSGS_CLIENT.CloseSession(context.Background())

	// Go routines that interact with blockchain

	ticker := time.NewTicker(UPDATE_FREQUENCY * time.Second)
	done := make(chan bool)
	go periodicallySubmitHash(done, ticker)

	ticker2 := time.NewTicker(POLL_FREQUENCY * time.Second)
	done2 := make(chan bool)
	go periodicallyPollAndSubmitPreImage(done2, ticker2)

	// HTTP server
	go startApiServer()

	// PromptUI to select action
	TopUI()

	// time.Sleep(1600 * time.Millisecond)
	// ticker.Stop()
	// done <- true
}

// TODO include blocknumber in message (maybe within a mutex between hash submisison and hash retrieval)
// TODO refill wait confirmation form other bank (same with interbank tx)
