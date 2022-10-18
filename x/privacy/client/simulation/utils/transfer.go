package utils

import "fmt"

func Transfer(privateKey0, privateKey1, paymentAddress string) {
	// check balance first time
	getBalanceIncognito(privateKey0)
	getBalanceIncognito(privateKey1)

	// airdrop
	Airdrop(privateKey0)
	Airdrop(privateKey1)

	// check balance second time
	getBalanceIncognito(privateKey0)
	getBalanceIncognito(privateKey1)

	// transfer
	fmt.Println("Start transfer")
	temp := paymentAddress + "-100"
	args := []string{"tx", "privacy", "transfer", privateKey0, temp, "0prv", "--from", "my_validator", "--chain-id", "my-test-chain", "-y"}
	execCmd(args, true)
	fmt.Println("Press enter to continue")
	fmt.Scanln()

	// check balance third time
	getBalanceIncognito(privateKey0)
	getBalanceIncognito(privateKey1)
}
