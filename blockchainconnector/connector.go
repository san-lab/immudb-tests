package blockchainconnector

import (
	"context"
	"crypto/ecdsa"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"

	. "github.com/san-lab/immudb-tests/datastructs"
)

var NETWORK string
var CHAIN_ID string
var VERIFIER_ADDRESS string
var PRIV_KEY_FILE string

func GetStateCheckByBlockNumber(originatorBank, recipientBank string, blockNumber *big.Int) (*StateCheck, error) {
	client, err := ethclient.Dial(NETWORK)
	if err != nil {
		return nil, err
	}

	address := common.HexToAddress(VERIFIER_ADDRESS)
	instance, err := NewOnchainVerifier(address, client)

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
	instance, err := NewOnchainVerifier(address, client)

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
	instance, err := NewOnchainVerifier(address, client)

	// Recipient must be ThisBank
	originatorBankAddress := common.HexToAddress(originatorBank)
	recipientBankAddress := common.HexToAddress(ThisBank.Address)

	pendingSubmissions, err := instance.GetPendingSubmissions(&bind.CallOpts{}, originatorBankAddress, recipientBankAddress)
	return pendingSubmissions, err

}

func SubmitHash(recipientBank string, hash []byte) error {
	client, err := ethclient.Dial(NETWORK)
	if err != nil {
		return err
	}

	address := common.HexToAddress(VERIFIER_ADDRESS)
	instance, err := NewOnchainVerifier(address, client)

	// Originator must be ThisBank
	originatorBankAddress := common.HexToAddress(ThisBank.Address)
	recipientBankAddress := common.HexToAddress(recipientBank)

	// Set signer parameters...
	auth, err := getAuth(client)
	if err != nil {
		return err
	}

	_, err = instance.SubmitHash(auth, originatorBankAddress, recipientBankAddress, [32]byte(hash))
	return err
}

func SubmitPreimage(originatorBank string, preimage []byte, blockNumber *big.Int) error {
	client, err := ethclient.Dial(NETWORK)
	if err != nil {
		return err
	}

	address := common.HexToAddress(VERIFIER_ADDRESS)
	instance, err := NewOnchainVerifier(address, client)

	// Recipient must be ThisBank
	originatorBankAddress := common.HexToAddress(originatorBank)
	recipientBankAddress := common.HexToAddress(ThisBank.Address)

	// Set signer parameters...
	auth, err := getAuth(client)
	if err != nil {
		return err
	}

	_, err = instance.SubmitPreimage(auth, originatorBankAddress, recipientBankAddress, [32]byte(preimage), blockNumber)
	return err
}

func GetBlockNumber() (*big.Int, error) {
	client, err := ethclient.Dial(NETWORK)
	if err != nil {
		return nil, err
	}

	header, err := client.HeaderByNumber(context.Background(), nil)
	if err != nil {
		return nil, err
	}

	return header.Number, nil
}

func getAuth(client *ethclient.Client) (*bind.TransactOpts, error) {
	// Read key from file
	privKeyBytes, err := os.ReadFile(PRIV_KEY_FILE)
	if err != nil {
		return nil, err
	}

	/*
		privKey := new(datastructs.PrivateKey)
		err = json.Unmarshal(privKeyBytes, privKey)
		if err != nil {
			return nil, err
		}
	*/

	privateKey, err := crypto.HexToECDSA(string(privKeyBytes))
	if err != nil {
		return nil, err
	}
	publicKey := privateKey.Public()
	publicKeyECDSA, _ := publicKey.(*ecdsa.PublicKey)

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		return nil, err
	}

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
