package main

import (
	"encoding/hex"
	"fmt"
	"time"

	account "github.com/san-lab/immudb-tests/account"
	"github.com/san-lab/immudb-tests/blockchainconnector"
	. "github.com/san-lab/immudb-tests/datastructs"
	"github.com/wealdtech/go-merkletree/keccak256"
)

// TODO: properly handle nonce ¿?
func periodicallySubmitHash(done chan bool, ticker *time.Ticker) {
	previousDigest := make(map[string]string)
	for {
		select {
		case <-done:
			return

		case t := <-ticker.C:
			fmt.Println("****** Tick at", t)
			mirrorAccounts, _ := account.GetAllAccounts("mirror")
			for _, mirrorAccount := range mirrorAccounts {
				recipientBank := COUNTERPART_BANKS[mirrorAccount.CABank]
				digest, _ := mirrorAccount.GetDigest()

				if ONLY_ON_CHANGES && digest != previousDigest[recipientBank] {
					digestBytes, err := hex.DecodeString(digest)
					if err != nil {
						fmt.Println(err)
						continue
					}
					hashBytes := keccak256.New().Hash(digestBytes)
					hash := fmt.Sprintf("%x", hashBytes)
					fmt.Println("* Submitting hash for", mirrorAccount.Iban)
					fmt.Println("* Digest", digest)
					fmt.Println("* Hash", hash)
					err = blockchainconnector.SubmitHash(recipientBank, hash)
					if err != nil {
						fmt.Println(err)
						continue
					}
					previousDigest[recipientBank] = digest
				}
				fmt.Println("* done!")
			}
		}
	}
}

// TODO: handle more than one pending submission properly
func periodicallyPollAndSubmitPreImage(done chan bool, ticker *time.Ticker) {
	for {
		select {
		case <-done:
			return

		case t := <-ticker.C:
			fmt.Println("------ Tick at", t)
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
					fmt.Println("- No pending submissions for this CA")
					continue
				}
				fmt.Println("- Pending submissions:", pending)
				digest, _ := CAAccount.GetDigest()
				fmt.Println("- Submitting preimage for", CAAccount.Iban)
				fmt.Println("- Digest", digest)
				err = blockchainconnector.SubmitPreimage(originatorBank, digest, pending[len(pending)-1])
				if err != nil {
					fmt.Println(err)
					continue
				}
				fmt.Println("- done!")
			}
		}
	}
}
