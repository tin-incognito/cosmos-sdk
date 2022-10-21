package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

func Benchmark(privateKey0, privateKey1, paymentAddress0, paymentAddress1, cosmosAccount0, cosmosAccount1 string, skipWaiting bool) {

	benchmarkWithStrategy(skipWaiting)

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

func benchmarkWithStrategy(skipWaiting bool) {
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

	// Airdrop to cosmos accounts
	// should run only one time
	/*for i, v := range cosmosAccounts {*/
	/*if i == 0 {*/
	/*continue*/
	/*}*/
	/*BankTransfer(string(cosmosAccounts[0]), string(v), 100000, true)*/
	/*time.Sleep(time.Second * 10)*/
	/*}*/

	for i, v := range cosmosAccounts {
		if i == 0 {
			continue
		}
		fmt.Println(v)
		rootPrivacyAccount := privacyAccounts[0]
		privacyAccount := privacyAccounts[i]
		//Shield(privacyAccount.PrivateKey, privacyAccount.PaymentAddress, string(v), 20000, skipWaiting)
		Transfer(privacyAccount.PrivateKey, rootPrivacyAccount.PrivateKey, rootPrivacyAccount.PaymentAddress, skipWaiting)
		/*Unshield(privacyAccount.PrivateKey, string(v), 200, skipWaiting)*/
	}
}
