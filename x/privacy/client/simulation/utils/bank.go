package utils

import "strconv"

func BankTransfer(from, to string, amount uint64, skipWaiting bool) {
	execCmd([]string{
		"tx", "bank", "send", from, to, strconv.Itoa(int(amount)) + "prv", "--chain-id", "my-test-chain", "-y",
	}, true)
}
