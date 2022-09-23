package main

import (
	"os"

	"github.com/cosmos/cosmos-sdk/x/privacy/client/simulation/utils"
)

func main() {
	scenario := os.Args[1]
	switch scenario {
	case "airdrop":
		utils.Airdrop("112t8rnXSbfsy9vSbst5HGVY7XbKK4exEFPkWPho41a6NC9PN5WNsS9ieJJcDyEFrkzfjaneU52WzJH4WFDKoz9Vv7qH82wjQfY3MpeUpJUc")
	case "transfer":
		utils.Transfer("112t8rnXSbfsy9vSbst5HGVY7XbKK4exEFPkWPho41a6NC9PN5WNsS9ieJJcDyEFrkzfjaneU52WzJH4WFDKoz9Vv7qH82wjQfY3MpeUpJUc", "112t8rnXXKcvUtpwcK1HwRsmJvTxp8mCB12Zma7dKWFSNNDkdhjdgg8UCftueJy7t5rStvqYYDvddPiJSTcKpmTXJkHZW9smo3K5moCFZiei", "12skvRi6rzvj8UwUhxcr3xe8Z5YDQyaYvDf7e3FQUhyDMNKLgf1MyXh6GnWPbdvGH898PU1duHYzBJK3Qs7dx75VtttJXjm8aadp3ozDM4XnohccCi4dEdgte8o8n6ff29RHnqE3zcaLirTTcM25")
	}
}
