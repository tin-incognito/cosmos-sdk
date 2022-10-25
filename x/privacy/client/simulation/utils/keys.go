package utils

import (
	"encoding/json"
	"io/ioutil"
	"strconv"
	"strings"
)

func AddKeys(numAcc int) {
	var addresses []string
	for i := 0; i < numAcc; i++ {
		receipient := "receipient" + strconv.Itoa(i)
		temp := execCmd([]string{"keys", "add", receipient}, false)
		str := strings.Split(temp, "\n")
		s := strings.Split(str[2], " ")
		addresses = append(addresses, s[3])
	}
	// commit rawtxs
	data, err := json.Marshal(addresses)
	if err != nil {
		panic(err)
	}
	err = ioutil.WriteFile("cosmos_accounts.json", data, 0644)
	if err != nil {
		panic(err)
	}
}
