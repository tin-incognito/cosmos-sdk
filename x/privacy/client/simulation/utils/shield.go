package utils

import (
	"fmt"
	"strconv"
)

func Shield(privateKey, paymentAddress, cosmosAccount string, amount uint64, skipWaiting bool) {
	if !skipWaiting {
		// Check balance before shield
		getBalanceIncognito(privateKey, skipWaiting)
		getBalanceCosmos(cosmosAccount, skipWaiting)
	}

	execCmd([]string{
		"tx", "privacy", "shield", paymentAddress, strconv.Itoa(int(amount)),
		"--from", cosmosAccount, "--chain-id", "my-test-chain", "-y",
	}, true)
	if !skipWaiting {
		fmt.Scanln()
		fmt.Println("Press enter to continue")
	}

	if !skipWaiting {
		// Check balance after shield
		getBalanceIncognito(privateKey, skipWaiting)
		getBalanceCosmos(cosmosAccount, skipWaiting)
	}
}
