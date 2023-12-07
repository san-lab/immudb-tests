package main

import (
	"context"

	. "github.com/san-lab/immudb-tests/datastructs"
)

const STATE_DB = "defaultdb"
const MSGS_DB = "msgdb"

func main() {
	// Parse all config parameters
	initConfigParams()

	// Initialize DBs
	initDB()

	// Initialize digest history of onboarded CA banks
	initDigestHistory()

	// Get initial nonce
	initNonce()

	// Ensure connection is closed
	defer STATE_CLIENT.CloseSession(context.Background())
	defer MSGS_CLIENT.CloseSession(context.Background())

	// Go routines that interact with blockchain
	done := make(chan bool)
	go periodicallySubmitHash(done)

	done2 := make(chan bool)
	go periodicallyPollAndSubmitPreImage(done2)

	// HTTP server
	go startApiServer()

	// Automatically send broadcast message to find other libp2p banks
	done3 := make(chan bool)
	go periodicallyFindCounterpartBanks(done3)

	// PromptUI to select action
	TopUI()

	// time.Sleep(1600 * time.Millisecond)
	// ticker.Stop()
	// done <- true
}

// TODO use events and keep pendingSubmissions method as backup plan
// TODO refactor optimistic nonce update to have cleaner code
// TODO add previous submission blocknumber to keep submitting blocks (privacy purposes)
