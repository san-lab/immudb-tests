package blockchainconnector

import (
	"context"
	"encoding/hex"
	"fmt"
	"math/big"
	"os"
	"sync"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"

	. "github.com/san-lab/immudb-tests/datastructs"
)

var nonceMutex sync.Mutex

// var blockNumberMutex sync.Mutex

// var TIMEOUT = 30
var NONCE int

var NETWORK string
var CHAIN_ID string
var VERIFIER_ADDRESS string
var PRIV_KEY_FILE string

type StateCheck struct {
	SubmittedHash     []byte
	SubmittedPreimage []byte
	Verified          bool
	BlockNumber       *big.Int
}

func GetStateCheckByBlockNumber(originatorBank, recipientBank string, blockNumber *big.Int) (*StateCheck, error) {
	client, err := ethclient.Dial(NETWORK)
	if err != nil {
		return nil, err
	}

	address := common.HexToAddress(VERIFIER_ADDRESS)
	instance, err := NewOnChainVerifier(address, client)
	if err != nil {
		return nil, err
	}

	originatorBankAddress := common.HexToAddress(originatorBank)
	recipientBankAddress := common.HexToAddress(recipientBank)

	stateCheckSC, err := instance.GetStateCheckByBlockNumber(&bind.CallOpts{}, originatorBankAddress, recipientBankAddress, blockNumber)
	stateCheck := &StateCheck{
		SubmittedHash:     stateCheckSC.SubmittedHash[:],
		SubmittedPreimage: stateCheckSC.SubmittedPreimage[:],
		Verified:          stateCheckSC.Verified,
		BlockNumber:       stateCheckSC.BlockNumber,
	}

	return stateCheck, err
}

func GetStateCheckByIndex(originatorBank, recipientBank string, index *big.Int) (*StateCheck, error) {
	client, err := ethclient.Dial(NETWORK)
	if err != nil {
		return nil, err
	}

	address := common.HexToAddress(VERIFIER_ADDRESS)
	instance, err := NewOnChainVerifier(address, client)
	if err != nil {
		return nil, err
	}

	originatorBankAddress := common.HexToAddress(originatorBank)
	recipientBankAddress := common.HexToAddress(recipientBank)

	stateCheckSC, err := instance.GetStateCheckByIndex(&bind.CallOpts{}, originatorBankAddress, recipientBankAddress, index)
	stateCheck := &StateCheck{
		SubmittedHash:     stateCheckSC.SubmittedHash[:],
		SubmittedPreimage: stateCheckSC.SubmittedPreimage[:],
		Verified:          stateCheckSC.Verified,
		BlockNumber:       stateCheckSC.BlockNumber,
	}

	return stateCheck, err
}

func GetPendingSubmissions(originatorBank string) ([]*big.Int, error) {
	client, err := ethclient.Dial(NETWORK)
	if err != nil {
		return nil, err
	}

	address := common.HexToAddress(VERIFIER_ADDRESS)
	instance, err := NewOnChainVerifier(address, client)
	if err != nil {
		return nil, err
	}

	// Recipient must be ThisBank
	originatorBankAddress := common.HexToAddress(originatorBank)
	recipientBankAddress := common.HexToAddress(THIS_BANK.Address)

	pendingSubmissions, err := instance.GetPendingSubmissions(&bind.CallOpts{From: recipientBankAddress}, originatorBankAddress, recipientBankAddress)
	return pendingSubmissions, err

}

func SubmitHash(recipientBank string, hash string) error {
	/*
		fmt.Println("debug hash lock")
		blockNumberMutex.Lock()

		defer fmt.Println("debug hash unlock")
		defer blockNumberMutex.Unlock()
	*/

	client, err := ethclient.Dial(NETWORK)
	if err != nil {
		return err
	}

	address := common.HexToAddress(VERIFIER_ADDRESS)
	instance, err := NewOnChainVerifier(address, client)
	if err != nil {
		return err
	}

	// Originator must be ThisBank
	originatorBankAddress := common.HexToAddress(THIS_BANK.Address)
	recipientBankAddress := common.HexToAddress(recipientBank)

	// Set signer parameters...
	auth, err := getAuth(client)
	if err != nil {
		return err
	}

	hashBytes, err := hex.DecodeString(hash)
	if err != nil {
		return err
	}
	_, err = instance.SubmitHash(auth, originatorBankAddress, recipientBankAddress, [32]byte(hashBytes))
	if err != nil { // TODO: check for the appropriate error
		fmt.Println("debug submit hash error?", err)
		handleNonceError()
	}
	// Poll for receipt (when nonce is updated) until TIMEOUT
	/*
		stay := true
		timeout := time.After(time.Duration(TIMEOUT) * time.Second)
		var receipt *types.Receipt
		for stay {
			receipt, err = client.TransactionReceipt(context.Background(), tx.Hash())
			fmt.Println("debug hash receipt", err, receipt)
			if receipt != nil && err == nil {
				fmt.Println("debug preimage receipt found it", err, receipt)
				// Found it
				break
			}
			select {
			case <-timeout:
				fmt.Println("debug hash receipt timed out")
				stay = false
			default:
			}
		}
	*/
	return err
}

func SubmitPreimage(originatorBank string, preimage string, blockNumber *big.Int) error {
	/*
		fmt.Println("debug preimage lock")
		blockNumberMutex.Lock()

		defer fmt.Println("debug preimage unlock")
		defer blockNumberMutex.Unlock()
	*/
	client, err := ethclient.Dial(NETWORK)
	if err != nil {
		return err
	}

	address := common.HexToAddress(VERIFIER_ADDRESS)
	instance, err := NewOnChainVerifier(address, client)
	if err != nil {
		return err
	}

	// Recipient must be ThisBank
	originatorBankAddress := common.HexToAddress(originatorBank)
	recipientBankAddress := common.HexToAddress(THIS_BANK.Address)

	// Set signer parameters...
	auth, err := getAuth(client)
	if err != nil {
		return err
	}

	preimageBytes, err := hex.DecodeString(preimage)
	if err != nil {
		return err
	}
	_, err = instance.SubmitPreimage(auth, originatorBankAddress, recipientBankAddress, [32]byte(preimageBytes), blockNumber)
	if err != nil { // TODO: check for the appropriate error
		fmt.Println("debug submit preimage error?", err)
		handleNonceError()
	}
	// Poll for receipt (when nonce is updated) until TIMEOUT
	/*
		stay := true
		timeout := time.After(time.Duration(TIMEOUT) * time.Second)
		var receipt *types.Receipt
		for stay {
			receipt, err = client.TransactionReceipt(context.Background(), tx.Hash())
			fmt.Println("debug preimage receipt", err, receipt)
			if receipt != nil && err == nil {
				fmt.Println("debug preimage receipt found it", err)
				// Found it
				break
			}
			select {
			case <-timeout:
				fmt.Println("debug preimage receipt timed out")
				stay = false
			default:
			}
		}
	*/
	return err
}

func handleNonceError() {
	// TODO: Handle nonce recovery properly
	// Need mutex to prevent method calling getLocalNonce while it is being updated here
	// fmt.Println("debug nonce error waiting lock...")
	nonceMutex.Lock()

	// defer fmt.Println("debug nonce error unlocked.")
	defer nonceMutex.Unlock()

	nonce, err := GetBlockchainNonce()
	if err != nil {
		fmt.Println(err)
	}
	NONCE = nonce
}

func GetBlockNumber() (int, error) {
	// fmt.Println("debug getBlockNumber lock")
	// blockNumberMutex.Lock()

	// defer fmt.Println("debug getBlockNumber unlock")
	// defer blockNumberMutex.Unlock()

	client, err := ethclient.Dial(NETWORK)
	if err != nil {
		return 0, err
	}

	header, err := client.HeaderByNumber(context.Background(), nil)
	if err != nil {
		return 0, err
	}

	return int(header.Number.Int64()), nil
}

func Version() (string, error) {
	client, err := ethclient.Dial(NETWORK)
	if err != nil {
		return "", err
	}

	address := common.HexToAddress(VERIFIER_ADDRESS)
	instance, err := NewOnChainVerifier(address, client)
	if err != nil {
		return "", err
	}

	version, err := instance.Version(&bind.CallOpts{})
	return version, err
}

func WatchHashSubmittedEvent(sink chan *OnChainVerifierHashSubmitted) error {
	client, err := ethclient.Dial(NETWORK)
	if err != nil {
		return err
	}

	address := common.HexToAddress(VERIFIER_ADDRESS)
	instance, err := NewOnChainVerifier(address, client)
	if err != nil {
		return err
	}

	opts := &bind.WatchOpts{}
	instance.OnChainVerifierFilterer.WatchHashSubmitted(opts, sink)
	return err
}

// Optimistic approach
func GetAndIncreaseLocalNonce() int {
	// fmt.Println("debug waiting lock...")
	nonceMutex.Lock()

	// defer fmt.Println("debug unlocked.")
	defer nonceMutex.Unlock()

	currentNonce := NONCE
	NONCE = NONCE + 1
	// fmt.Println("debug nonce", currentNonce)
	return currentNonce
}

// For the nonce to advance we should wait for the tx receipt,
// in which case we are no better than network latency
// If we dont wait, we may get desynced

func GetBlockchainNonce() (int, error) {
	client, err := ethclient.Dial(NETWORK)
	if err != nil {
		return 0, err
	}
	nonce, err := client.PendingNonceAt(context.Background(), common.HexToAddress(THIS_BANK.Address))
	return int(nonce), err
}

func getAuth(client *ethclient.Client) (*bind.TransactOpts, error) {
	// Read key from file
	privKeyBytes, err := os.ReadFile(PRIV_KEY_FILE)
	if err != nil {
		return nil, err
	}

	privateKey, err := crypto.HexToECDSA(string(privKeyBytes))
	if err != nil {
		return nil, err
	}

	/*
		publicKey := privateKey.Public()
		publicKeyECDSA, _ := publicKey.(*ecdsa.PublicKey)
		fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
		nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
		if err != nil {
			return nil, err
		}
	*/

	nonce := GetAndIncreaseLocalNonce()

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		return nil, err
	}

	chainID, ok := new(big.Int).SetString(CHAIN_ID, 10)
	if !ok {
		return nil, err
	}
	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		return nil, err
	}
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)     // in wei
	auth.GasLimit = uint64(300000) // in units
	auth.GasPrice = gasPrice

	return auth, nil
}

func PrintStateCheck(state *StateCheck) {
	fmt.Printf("{ Submitted Hash: %x\nSubmitted Preimage: %x\nVerified?: %t\nBlock Number: %d }\n",
		state.SubmittedHash, state.SubmittedPreimage, state.Verified, state.BlockNumber)
}
