package utils

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/x/privacy/repos/key"
	"github.com/incognitochain/go-incognito-sdk-v2/wallet"
)

func Airdrop(privateKey string) {
	getBalanceIncognito(privateKey)
	fmt.Println("Airdrop 100000 token to account with privateKey", privateKey)
	keyWallet, err := wallet.Base58CheckDeserialize(privateKey)
	if err != nil {
		panic(err)
	}
	keySet := key.KeySet{}
	err = keySet.InitFromPrivateKeyByte(keyWallet.KeySet.PrivateKey)
	if err != nil {
		panic(err)
	}

	args := []string{"tx", "privacy", "airdrop", privateKey, "100000", "--from", "my_validator", "--chain-id", "my-test-chain", "-y"}
	execCmd(args, true)
	fmt.Println("Press enter to continue")
	fmt.Scanln()
	getBalanceIncognito(privateKey)
}
