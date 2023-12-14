package main

import (
	"encoding/hex"
	"fmt"
	"log"
	"time"

	"github.com/san-lab/immudb-tests/account"
	"github.com/san-lab/immudb-tests/bankinterop"
	"github.com/san-lab/immudb-tests/blockchainconnector"
	"github.com/san-lab/immudb-tests/color"
	. "github.com/san-lab/immudb-tests/datastructs"
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
			log.Printf("\n****** Tick at %s", t)
			mirrorAccounts, _ := account.GetAllAccounts(account.MIRROR)
			for _, mirrorAccount := range mirrorAccounts {
				log.Printf("*** Checking Mirror Account %s", mirrorAccount.Holder)
				recipientBankAddress := COUNTERPART_BANKS[mirrorAccount.CABank]
				digest, _ := mirrorAccount.GetDigest()
				// TODO: initialize previousDigest quering blockchain StateCheck
				log.Printf("* Previous digest compare %s =? %s", color.Shorten(digest, 10), color.Shorten(previousDigest[mirrorAccount.CABank], 10))
				if ONLY_ON_CHANGES && digest != previousDigest[mirrorAccount.CABank] {
					digestBytes, err := hex.DecodeString(digest)
					if err != nil {
						fmt.Println(err)
						continue
					}
					hashBytes := keccak256.New().Hash(digestBytes)
					hash := fmt.Sprintf("%x", hashBytes)
					log.Printf("*** Submitting hash for %s", mirrorAccount.Holder)
					log.Printf("* Digest %s", color.Shorten(digest, 10))
					log.Printf("* Hash %s", color.Shorten(hash, 10))
					log.Printf("* RecipientBankAddress %s", recipientBankAddress)
					err = blockchainconnector.SubmitHash(recipientBankAddress, hash)
					if err != nil {
						fmt.Println(err)
						continue
					}
					previousDigest[mirrorAccount.CABank] = digest
				}
				log.Printf("** Done!")
			}
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
			log.Printf("\n------ Tick at %s", t)
			CAAccounts, _ := account.GetAllAccounts(account.CA)
			for _, CAAccount := range CAAccounts {
				log.Printf("--- Checking CA %s", CAAccount.Holder)
				originatorBankAddress := COUNTERPART_BANKS[CAAccount.CABank]
				pending, err := blockchainconnector.GetPendingSubmissions(originatorBankAddress)
				if err != nil {
					fmt.Println(err)
					continue
				}
				log.Printf("- Pending submissions: %v", pending)
				// Either no pending, or pending older that what we can provide
				if len(pending) == 0 || int(pending[len(pending)-1].Int64()) < FIRST_BLOCK_NUMBER {
					log.Printf("- No pending submissions for this CA")
					continue
				}

				//log.Printf("- Full digest history %v", bankinterop.DigestHistory[CAAccount.CABank])
				log.Printf("-- Full digest history:")
				for k, v := range bankinterop.DigestHistory[CAAccount.CABank] {
					log.Printf("- %d -> %s", k, color.Shorten(v, 10))
				}

				digest, err := bankinterop.PickLatestDigestPriorToResquestedBlockNumber(CAAccount.CABank, pending[len(pending)-1])
				if err != nil {
					fmt.Println(err)
					continue
				}
				log.Printf("- Submitting preimage for %s", CAAccount.Holder)
				log.Printf("- Digest %s", color.Shorten(digest, 10))
				log.Printf("- OriginatorBankAddress %s", originatorBankAddress)
				err = blockchainconnector.SubmitPreimage(originatorBankAddress, digest, pending[len(pending)-1])
				if err != nil {
					fmt.Println(err)
					continue
				}
				log.Printf("-- Done!")
			}
		}
	}
}
