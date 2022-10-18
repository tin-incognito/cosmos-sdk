package utils

import (
	"fmt"
	"strconv"
)

func Unshield(privateKey, cosmosAccount string, amount uint64) {
	// Check balance before shield
	getBalanceIncognito(privateKey)
	getBalanceCosmos(cosmosAccount)

	execCmd([]string{
		"tx", "privacy", "unshield", privateKey, cosmosAccount,
		strconv.Itoa(int(amount)), "0prv",
		"--from", "my_validator", "--chain-id", "my-test-chain", "-y",
	}, true)
	fmt.Scanln()
	fmt.Println("Press enter to continue")

	// Check balance after shield
	getBalanceIncognito(privateKey)
	getBalanceCosmos(cosmosAccount)
}
