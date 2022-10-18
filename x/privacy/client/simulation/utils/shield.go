package utils

import (
	"fmt"
	"strconv"
)

func Shield(privateKey, paymentAddress, cosmosAccount string, amount uint64) {
	// Check balance before shield
	getBalanceIncognito(privateKey)
	getBalanceCosmos(cosmosAccount)

	execCmd([]string{
		"tx", "privacy", "shield", paymentAddress, strconv.Itoa(int(amount)),
		"--from", cosmosAccount, "--chain-id", "my-test-chain", "-y",
	}, true)
	fmt.Scanln()
	fmt.Println("Press enter to continue")

	// Check balance after shield
	getBalanceIncognito(privateKey)
	getBalanceCosmos(cosmosAccount)
}
