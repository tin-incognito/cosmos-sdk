package main

import (
	"os"

	"github.com/cosmos/cosmos-sdk/x/privacy/client/simulation/utils"
)

func main() {
	scenario := os.Args[1]
	switch scenario {
	case "airdrop":
		utils.Airdrop("112t8roMiguHp5SZxXpnoTHiV8A7XWkm79Y4hWwrqjtvjDH52H4y6z5RBt8HqqmKXSZ33CGjd2uM4ZiG32hqMzytboMA7zrBGscciRRVG4rH", false)
	case "transfer":
		utils.Transfer("112t8roMiguHp5SZxXpnoTHiV8A7XWkm79Y4hWwrqjtvjDH52H4y6z5RBt8HqqmKXSZ33CGjd2uM4ZiG32hqMzytboMA7zrBGscciRRVG4rH", "112t8rqkMBoQDtSFggSfJyji71BxE2HV1kTSarohB7rUgy5qYMecCqvLoWM8vN3KMVr8RQskdQvJLSK1qducBPnG2v9xLFmkvyaKUNp97y8X", "12soNT5whtDTZSh3Twak286DmwiXAm19edjqkCX7kF9koickEaeCPkUAQUwHEzhkSn7YuBm1jdWR14t8Q1UeCBLycwwyawvw2iJkre572Y2ZAXFSPNkjfJhio6kJP1s6Gj3fEjrq4RL7YCpPfUWQ", false)
	case "shield":
		utils.Shield("112t8roMiguHp5SZxXpnoTHiV8A7XWkm79Y4hWwrqjtvjDH52H4y6z5RBt8HqqmKXSZ33CGjd2uM4ZiG32hqMzytboMA7zrBGscciRRVG4rH", "12scjGeftnVsH4Xa2CsRpkqctZaY77asQQMq84Gqx3itNRZaaQxfCDARXfrVZrsSK63pGPC2DzYwdEhLAjAyrUKErGeZfcL2v7HXXTVLee6Gwvr5NsruJRCqiHnQ9aaGsYGDKy8mgTzu1pJfPdv8", "incog1q9kcvyf89eewavtd8lgu3zh6qx3k67y8tlqkk8", 500, false)
	case "unshield":
		utils.Unshield("112t8roMiguHp5SZxXpnoTHiV8A7XWkm79Y4hWwrqjtvjDH52H4y6z5RBt8HqqmKXSZ33CGjd2uM4ZiG32hqMzytboMA7zrBGscciRRVG4rH", "incog1q9kcvyf89eewavtd8lgu3zh6qx3k67y8tlqkk8", 200, false)
	case "benchmark":
		utils.Benchmark(
			"112t8roMiguHp5SZxXpnoTHiV8A7XWkm79Y4hWwrqjtvjDH52H4y6z5RBt8HqqmKXSZ33CGjd2uM4ZiG32hqMzytboMA7zrBGscciRRVG4rH",
			"112t8rqkMBoQDtSFggSfJyji71BxE2HV1kTSarohB7rUgy5qYMecCqvLoWM8vN3KMVr8RQskdQvJLSK1qducBPnG2v9xLFmkvyaKUNp97y8X",
			"12scjGeftnVsH4Xa2CsRpkqctZaY77asQQMq84Gqx3itNRZaaQxfCDARXfrVZrsSK63pGPC2DzYwdEhLAjAyrUKErGeZfcL2v7HXXTVLee6Gwvr5NsruJRCqiHnQ9aaGsYGDKy8mgTzu1pJfPdv8",
			"12soNT5whtDTZSh3Twak286DmwiXAm19edjqkCX7kF9koickEaeCPkUAQUwHEzhkSn7YuBm1jdWR14t8Q1UeCBLycwwyawvw2iJkre572Y2ZAXFSPNkjfJhio6kJP1s6Gj3fEjrq4RL7YCpPfUWQ",
			"incog1q9kcvyf89eewavtd8lgu3zh6qx3k67y8tlqkk8",
			"incog187jvy7vxu33savdjz7pxwecr4qz55run7jwuj0",
			true,
		)
	}
}
