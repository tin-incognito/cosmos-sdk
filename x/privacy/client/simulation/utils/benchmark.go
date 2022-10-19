package utils

import (
	"encoding/json"
	"io/ioutil"
)

func Benchmark(privateKey0, privateKey1, paymentAddress0, paymentAddress1, cosmosAccount0, cosmosAccount1 string, skipWaiting bool) {
	//Airdrop(privateKey0, skipWaiting)
	//Shield(privateKey0, paymentAddress0, cosmosAccount0, 200, skipWaiting)
	//Transfer(privateKey0, privateKey1, paymentAddress1, skipWaiting)
	Unshield(privateKey0, cosmosAccount0, 200, skipWaiting)

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
