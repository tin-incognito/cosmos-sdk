package utils

import "fmt"

func Transfer(privateKey0, privateKey1, paymentAddress string) {
	// check balance first time
	fmt.Println("balance of", privateKey0)
	args := []string{"query", "privacy", "balance", privateKey0}
	execCmd(args)
	fmt.Println("balance of", privateKey1)
	args = []string{"query", "privacy", "balance", privateKey1}
	execCmd(args)
	fmt.Println("Press enter to continue")
	fmt.Scanln()

	// airdrop
	Airdrop(privateKey0)
	Airdrop(privateKey1)

	// check balance second time
	fmt.Println("balance of", privateKey0)
	args = []string{"query", "privacy", "balance", privateKey0}
	execCmd(args)
	fmt.Println("balance of", privateKey1)
	args = []string{"query", "privacy", "balance", privateKey1}
	execCmd(args)
	fmt.Println("Press enter to continue")
	fmt.Scanln()

	// transfer
	temp := paymentAddress + "-10000"
	args = []string{"tx", "privacy", "transfer", privateKey0, temp, "0", "--from", "alice", "-y"}
	execCmd(args)

	// check balance third time
	fmt.Println("balance of", privateKey0)
	args = []string{"query", "privacy", "balance", privateKey0}
	execCmd(args)
	fmt.Println("balance of", privateKey1)
	args = []string{"query", "privacy", "balance", privateKey1}
	execCmd(args)
	fmt.Scanln()
}
