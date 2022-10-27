package utils

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"time"
)

func Benchmark(validatorPrivateKey, validatorPaymentAddress string, skipWaiting bool) {
	validatorCosmosAccount := os.Args[2]

	benchmarkWithStrategy(validatorPrivateKey, validatorPaymentAddress, validatorCosmosAccount, skipWaiting)

	// commit rawtxs
	data, err := json.Marshal(rawTxs)
	if err != nil {
		panic(err)
	}
	err = ioutil.WriteFile("output.json", data, 0644)
	if err != nil {
		panic(err)
	}
}

func benchmarkWithStrategy(validatorPrivateKey, validatorPaymentAddress, validatorCosmosAccount string, skipWaiting bool) {
	var cosmosAccounts []CosmosAccount
	var privacyAccounts []PrivacyAccount

	data, err := ioutil.ReadFile("cosmos_accounts.json")
	if err != nil {
		panic(err)
	}

	if err := json.Unmarshal(data, &cosmosAccounts); err != nil {
		panic(err)
	}

	data, err = ioutil.ReadFile("privacy_accounts.json")
	if err != nil {
		panic(err)
	}

	if err := json.Unmarshal(data, &privacyAccounts); err != nil {
		panic(err)
	}

	step := os.Args[3]

	// Airdrop to cosmos accounts
	// should run only one time
	if step == "1" {
		for i, v := range cosmosAccounts {
			if i == 0 {
				continue
			}
			BankTransfer(validatorCosmosAccount, string(v), 100000, true)
			time.Sleep(time.Second * 5)
		}
	}

	if step == "2" {
		for i, v := range cosmosAccounts {
			if i == 0 {
				continue
			}
			privacyAccount := privacyAccounts[i]
			Shield(privacyAccount.PrivateKey, privacyAccount.PaymentAddress, string(v), 90000, skipWaiting)
		}
	}

	if step == "3" {
		for i := range cosmosAccounts {
			if i == 0 {
				continue
			}
			privacyAccount := privacyAccounts[i]
			//Shield(privacyAccount.PrivateKey, privacyAccount.PaymentAddress, string(v), 90000, skipWaiting)
			Transfer(privacyAccount.PrivateKey, validatorPrivateKey, validatorPaymentAddress, skipWaiting)
			//Unshield(privacyAccount.PrivateKey, string(v), 200, skipWaiting)
		}
	}

	if step == "4" {
		for i, v := range cosmosAccounts {
			if i == 0 {
				continue
			}
			privacyAccount := privacyAccounts[i]
			Unshield(privacyAccount.PrivateKey, string(v), 200, skipWaiting)
		}
	}

	if step == "5" {
		for i, v := range cosmosAccounts {
			if i == 0 {
				continue
			}
			BankTransfer(validatorCosmosAccount, string(v), 100000, true)
			privacyAccount := privacyAccounts[i]
			Shield(privacyAccount.PrivateKey, privacyAccount.PaymentAddress, string(v), 90000, skipWaiting)
			Transfer(privacyAccount.PrivateKey, validatorPrivateKey, validatorPaymentAddress, skipWaiting)
			Unshield(privacyAccount.PrivateKey, string(v), 200, skipWaiting)
		}
	}
}
