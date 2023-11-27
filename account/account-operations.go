package account

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

	. "github.com/san-lab/immudb-tests/datastructs"
	sdk "github.com/san-lab/immudb-tests/immudbsdk"
)

const CA = "ca"
const MIRROR = "mirror"

// CA
// IBAN: OtherBank@MyBankIBAN
// Holder: OtherBank - CA

// Mirror
// IBAN: MyBank@OtherBankIBAN
// Holder: MyBank @ OtherBank - Mirror

func CAAccountIBAN(cABank string) string {
	return strings.ReplaceAll(cABank+"@"+THIS_BANK.Name+"IBAN", " ", "")
}

func CAAccountHolder(cABank string) string {
	return cABank + " - CA"
}

func MirrorAccountIBAN(cABank string) string {
	return strings.ReplaceAll(THIS_BANK.Name+"@"+cABank+"IBAN", " ", "")
}

func MirrorAccountHolder(cABank string) string {
	return THIS_BANK.Name + " @ " + cABank + " - Mirror"
}

func CreateAccount(bic, iban, holder, currency, cABank string, balance float32, suspended, isCA, isMirror bool) error {
	// Check if IBAN already in database
	_, err := GetAndDeserializeAccount(iban)
	if err == nil {
		return errors.New("account with that IBAN is already in the database")
	}

	accountState := SetAccount(bic, iban, holder, currency, cABank, balance, suspended, isCA, isMirror)

	err = SerializeAndSetAccount(iban, accountState)
	return err
}

func CreateCAAccount(bic, currency, cABank string, balance float32) error {
	iban := CAAccountIBAN(cABank)
	holder := CAAccountHolder(cABank)
	err := CreateAccount(bic, iban, holder, currency, cABank, balance, false, true, false)
	return err
}

func CreateMirrorAccount(bic, currency, cABank string, balance float32) error {
	iban := MirrorAccountIBAN(cABank)
	holder := MirrorAccountHolder(cABank)
	err := CreateAccount(bic, iban, holder, currency, cABank, balance, false, false, true)
	return err
}

func SuspendAccount(userIban string) error {
	accountState, err := GetAndDeserializeAccount(userIban)
	if err != nil {
		return err
	}
	accountState.Suspend()
	err = SerializeAndSetAccount(userIban, accountState)
	return err
}

func UnsuspendAccount(userIban string) error {
	accountState, err := GetAndDeserializeAccount(userIban)
	if err != nil {
		return err
	}
	accountState.Unsuspend()
	err = SerializeAndSetAccount(userIban, accountState)
	return err
}

func SetAccountBalance(userIban, balanceString string) error {
	balance, err := strconv.ParseFloat(balanceString, 32)
	if err != nil {
		fmt.Println(err)
		return err
	}

	accountState, err := GetAndDeserializeAccount(userIban)
	if err != nil {
		return err
	}
	err = accountState.SetBalance(float32(balance))
	if err != nil {
		return err
	}
	err = SerializeAndSetAccount(userIban, accountState)
	return err
}

func DepositToAccount(userIban, amountString string) error {
	amount, err := strconv.ParseFloat(amountString, 32)
	if err != nil {
		fmt.Println(err)
		return err
	}

	accountState, err := GetAndDeserializeAccount(userIban)
	if err != nil {
		return err
	}
	err = accountState.Deposit(float32(amount))
	if err != nil {
		return err
	}
	err = SerializeAndSetAccount(userIban, accountState)
	return err
}

func WithdrawFromAccount(userIban, amountString string) error {
	amount, err := strconv.ParseFloat(amountString, 32)
	if err != nil {
		fmt.Println(err)
		return err
	}

	accountState, err := GetAndDeserializeAccount(userIban)
	if err != nil {
		return err
	}
	err = accountState.Withdraw(float32(amount))
	if err != nil {
		return err
	}
	err = SerializeAndSetAccount(userIban, accountState)
	return err
}

func GetAccount(userIban string) (*Account, error) {
	accountState, err := GetAndDeserializeAccount(userIban)
	if err != nil {
		return nil, err
	}
	return accountState, nil
}

func GetCABank(userIban string) (string, error) {
	accountState, err := GetAndDeserializeAccount(userIban)
	if err != nil {
		return "", err
	}
	cABank, err := accountState.GetCABank()
	return cABank, err
}

func GetAccountDigest(userIban string) (string, error) {
	accountState, err := GetAndDeserializeAccount(userIban)
	if err != nil {
		return "", err
	}
	accountDigest, err := accountState.GetDigest()
	return accountDigest, err
}

func GetAndDeserializeAccount(key string) (*Account, error) {
	// Pick state
	accountStateRaw, err := sdk.VerifiedGet(key)
	if err != nil {
		return nil, err
	}
	accountState := new(Account)
	err = json.Unmarshal(accountStateRaw.Value, accountState)
	return accountState, err
}

func SerializeAndSetAccount(key string, accountState *Account) error {
	// Marshal and set new state into the DB
	finalAccountState, err := json.Marshal(accountState)
	if err != nil {
		return err
	}
	err = sdk.VerifiedSet(key, string(finalAccountState))
	return err
}

func GetAllAccounts(filter string) ([]*Account, error) {
	entries, err := sdk.GetAllEntries()
	if err != nil {
		return nil, err
	}
	accounts := []*Account{}
	for _, entry := range entries.Entries {
		accountState := new(Account)
		err := json.Unmarshal(entry.Value, accountState)
		if err != nil {
			return accounts, err
		}

		switch filter {
		case CA:
			if accountState.IsCA {
				accounts = append(accounts, accountState)
			}
		case MIRROR:
			if accountState.IsMirror {
				accounts = append(accounts, accountState)
			}
		case "":
			accounts = append(accounts, accountState)
		default:
			fmt.Println("u shouldn be here...")
		}

	}
	return accounts, nil
}
