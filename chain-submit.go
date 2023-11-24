package main

import (
	"encoding/hex"
	"fmt"
	"time"

	"github.com/fatih/color"
	account "github.com/san-lab/immudb-tests/account"
	"github.com/san-lab/immudb-tests/bankinterop"
	"github.com/san-lab/immudb-tests/blockchainconnector"
	. "github.com/san-lab/immudb-tests/datastructs"
	"github.com/wealdtech/go-merkletree/keccak256"
)

const UPDATE_FREQUENCY = 60
const POLL_FREQUENCY = 25
const ONLY_ON_CHANGES = true

// TODO: properly handle nonce ¿?
func periodicallySubmitHash(done chan bool, ticker *time.Ticker) {
	previousDigest := make(map[string]string)
	for {
		select {
		case <-done:
			return

		case t := <-ticker.C:
			color.Cyan("\n****** Tick at %s", t)
			//fmt.Println(promptui.Styler(promptui.FGCyan)(fmt.Sprintf("****** Tick at %s", t)))
			mirrorAccounts, _ := account.GetAllAccounts(account.MIRROR)
			for _, mirrorAccount := range mirrorAccounts {
				recipientBankAddress := COUNTERPART_BANKS[mirrorAccount.CABank]
				digest, _ := mirrorAccount.GetDigest()
				// TODO: initialize previousDigest quering blockchain StateCheck
				color.Cyan("* debug previous digest %s =? %s", digest, previousDigest[mirrorAccount.CABank])
				if ONLY_ON_CHANGES && digest != previousDigest[mirrorAccount.CABank] {
					digestBytes, err := hex.DecodeString(digest)
					if err != nil {
						fmt.Println(err)
						continue
					}
					hashBytes := keccak256.New().Hash(digestBytes)
					hash := fmt.Sprintf("%x", hashBytes)
					color.Cyan("*** Submitting hash for %s", mirrorAccount.Holder)
					color.Cyan("* Digest %s", digest)
					color.Cyan("* Hash %s", hash)
					color.Cyan("* RecipientBankAddress %s", recipientBankAddress)
					err = blockchainconnector.SubmitHash(recipientBankAddress, hash)
					if err != nil {
						fmt.Println(err)
						continue
					}
					previousDigest[mirrorAccount.CABank] = digest
				}
				color.Cyan("** done!")
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
			color.White("\n------ Tick at %s", t)
			CAAccounts, _ := account.GetAllAccounts(account.CA)
			for _, CAAccount := range CAAccounts {
				color.White("--- Checking CA %s", CAAccount.Holder)
				originatorBankAddress := COUNTERPART_BANKS[CAAccount.CABank]
				//TODO: right now they need to discover each other, could instead save address.....
				pending, err := blockchainconnector.GetPendingSubmissions(originatorBankAddress)
				if err != nil {
					fmt.Println(err)
					continue
				}
				color.White("- Pending submissions: %v", pending)
				// Either no pending, or pending older that what we can provide
				if len(pending) == 0 || int(pending[len(pending)-1].Int64()) < FIRST_BLOCK_NUMBER {
					color.White("- No pending submissions for this CA")
					continue
				}
				color.White("- debug %v", bankinterop.DigestHistory[CAAccount.CABank])
				digest, err := bankinterop.PickLatestDigestPriorToResquestedBlockNumber(CAAccount.CABank, pending[len(pending)-1]) //CAAccount.GetDigest()
				if err != nil {
					fmt.Println(err)
					continue
				}
				color.White("- Submitting preimage for %s", CAAccount.Holder)
				color.White("- Digest %s", digest)
				color.White("- originatorBankAddress %s", originatorBankAddress)
				err = blockchainconnector.SubmitPreimage(originatorBankAddress, digest, pending[len(pending)-1])
				if err != nil {
					fmt.Println(err)
					continue
				}
				color.White("-- done!")
			}
		}
	}
}
