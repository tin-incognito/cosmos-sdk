package utils

import "fmt"

func Transfer(privateKey0, privateKey1, paymentAddress string, skipWaiting bool) {
	if !skipWaiting {
		//check balance first time
		getBalanceIncognito(privateKey0, skipWaiting)
		getBalanceIncognito(privateKey1, skipWaiting)
	}

	// airdrop
	/*Airdrop(privateKey0, skipWaiting)*/
	/*Airdrop(privateKey1, skipWaiting)*/

	if !skipWaiting {
		// check balance second time
		getBalanceIncognito(privateKey0, skipWaiting)
		getBalanceIncognito(privateKey1, skipWaiting)
	}

	// transfer
	temp := paymentAddress + "-100"
	args := []string{"tx", "privacy", "transfer", privateKey0, temp, "0prv", "--from", "my_validator", "--chain-id", "my-test-chain", "-y"}
	execCmd(args, true)
	if !skipWaiting {
		fmt.Println("Press enter to continue")
		fmt.Scanln()
		// check balance third time
		getBalanceIncognito(privateKey0, skipWaiting)
		getBalanceIncognito(privateKey1, skipWaiting)
	}
}
