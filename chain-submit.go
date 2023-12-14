package main

import (
	"encoding/hex"
	"fmt"
	"time"

	"github.com/san-lab/immudb-tests/account"
	"github.com/san-lab/immudb-tests/bankinterop"
	"github.com/san-lab/immudb-tests/color"
	. "github.com/san-lab/immudb-tests/datastructs"
	"github.com/san-lab/reconciliation-bank/blockchainconnector"
	"github.com/wealdtech/go-merkletree/keccak256"
)

const ONLY_ON_CHANGES = true

var UPDATE_FREQUENCY int
var POLL_FREQUENCY int

// TODO: properly handle nonce ¿?
func periodicallySubmitHash(done chan bool) {
	ticker := time.NewTicker(time.Duration(UPDATE_FREQUENCY) * time.Second)
	previousDigest := make(map[string]string)

	for {
		select {
		case <-done:
			return

		case t := <-ticker.C:
			color.CPrintln(color.CYAN, "\n****** Tick at %s", t)
			mirrorAccounts, _ := account.GetAllAccounts(account.MIRROR)
			for _, mirrorAccount := range mirrorAccounts {
				color.CPrintln(color.CYAN, "*** Checking Mirror Account %s", mirrorAccount.Holder)
				recipientBankAddress := COUNTERPART_BANKS[mirrorAccount.CABank]
				digest, _ := mirrorAccount.GetDigest()
				// TODO: initialize previousDigest quering blockchain StateCheck
				color.CPrintln(color.CYAN, "* Previous digest compare %s =? %s", color.Shorten(digest, 10), color.Shorten(previousDigest[mirrorAccount.CABank], 10))
				if ONLY_ON_CHANGES && digest != previousDigest[mirrorAccount.CABank] {
					digestBytes, err := hex.DecodeString(digest)
					if err != nil {
						fmt.Println(err)
						continue
					}
					hashBytes := keccak256.New().Hash(digestBytes)
					hash := fmt.Sprintf("%x", hashBytes)
					color.CPrintln(color.CYAN, "*** Submitting hash for %s", mirrorAccount.Holder)
					color.CPrintln(color.CYAN, "* Digest %s", color.Shorten(digest, 10))
					color.CPrintln(color.CYAN, "* Hash %s", color.Shorten(hash, 10))
					color.CPrintln(color.CYAN, "* RecipientBankAddress %s", recipientBankAddress)
					err = blockchainconnector.SubmitHash(recipientBankAddress, hash)
					if err != nil {
						fmt.Println(err)
						continue
					}
					previousDigest[mirrorAccount.CABank] = digest
				}
				color.CPrintln(color.CYAN, "** Done!")
			}
		}
	}
}

func listenForEventsAndSubmitPreImage(done chan bool) {
	sink := make(chan *blockchainconnector.OnChainVerifierHashSubmitted)
	err := blockchainconnector.WatchHashSubmittedEvent(sink)
	if err != nil {
		fmt.Println(err)
	}
	for {
		select {
		case <-done:
			return

		case <-sink:
			fmt.Println("got event, do something")
		}
	}
}

// TODO: handle more than one pending submission properly
func periodicallyPollAndSubmitPreImage(done chan bool) {
	ticker := time.NewTicker(time.Duration(POLL_FREQUENCY) * time.Second)

	for {
		select {
		case <-done:
			return

		case t := <-ticker.C:
			color.CPrintln(color.WHITE, "\n------ Tick at %s", t)
			CAAccounts, _ := account.GetAllAccounts(account.CA)
			for _, CAAccount := range CAAccounts {
				color.CPrintln(color.WHITE, "--- Checking CA %s", CAAccount.Holder)
				originatorBankAddress := COUNTERPART_BANKS[CAAccount.CABank]
				pending, err := blockchainconnector.GetPendingSubmissions(originatorBankAddress)
				if err != nil {
					fmt.Println(err)
					continue
				}
				color.CPrintln(color.WHITE, "- Pending submissions: %v", pending)
				// Either no pending, or pending older that what we can provide
				if len(pending) == 0 || int(pending[len(pending)-1].Int64()) < FIRST_BLOCK_NUMBER {
					color.CPrintln(color.WHITE, "- No pending submissions for this CA")
					continue
				}

				//color.CPrintln(color.WHITE, "- Full digest history %v", bankinterop.DigestHistory[CAAccount.CABank])
				color.CPrintln(color.WHITE, "-- Full digest history:")
				for k, v := range bankinterop.DigestHistory[CAAccount.CABank] {
					color.CPrintln(color.WHITE, "- %d -> %s", k, color.Shorten(v, 10))
				}

				digest, err := bankinterop.PickLatestDigestPriorToResquestedBlockNumber(CAAccount.CABank, pending[len(pending)-1])
				if err != nil {
					fmt.Println(err)
					continue
				}
				color.CPrintln(color.WHITE, "- Submitting preimage for %s", CAAccount.Holder)
				color.CPrintln(color.WHITE, "- Digest %s", color.Shorten(digest, 10))
				color.CPrintln(color.WHITE, "- OriginatorBankAddress %s", originatorBankAddress)
				err = blockchainconnector.SubmitPreimage(originatorBankAddress, digest, pending[len(pending)-1])
				if err != nil {
					fmt.Println(err)
					continue
				}
				color.CPrintln(color.WHITE, "-- Done!")
			}
		}
	}
}
