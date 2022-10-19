package utils

import (
	"bytes"
	"fmt"
	"os/exec"
)

var rawTxs []string

func execCmd(command []string, isRawTx bool) {
	fmt.Println("execFileDir:", execFileDir)
	fmt.Println("command:", command)
	cmd := exec.Command(execFileDir, command...)

	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
		return
	}
	fmt.Println("Result: " + out.String())
	if isRawTx {
		rawTxs = append(rawTxs, out.String())
	}
}

func getBalanceIncognito(privateKey string, skipWaiting bool) {
	fmt.Println("balance of", privateKey)
	args := []string{"query", "privacy", "balance", privateKey}
	execCmd(args, false)
	if !skipWaiting {
		fmt.Println("Press enter to continue")
		fmt.Scanln()
	}
}

func getBalanceCosmos(cosmosAccount string, skipWaiting bool) {
	args := []string{"query", "bank", "balances", cosmosAccount, "--chain-id", "my-test-chain"}
	execCmd(args, false)
	if !skipWaiting {
		fmt.Println("Press enter to continue")
		fmt.Scanln()
	}
}
