To benchmark privacy cosmos

### Prerequisites:

1. Init `privacy_accounts.json` file by `incognito cli` repo (use this branch `benchmark/init-privacy-accounts`)
2. Init `cosmos_accounts.json` file by run this command `go run main.go keys {{amount}}` in `x/privacy/client/simulation`
3. With amount is the amount of cosmos account number

### Change password:

Choose one in two options

1. Change password in file keyring/keyring.go to your os password
2. Change password in your os to "11234566"

### Flow to run:

Run by these commands in exactly order

1. Send as many as possible token to validator_address
1. `go run main.go {{validator_address}} 1` (bank)
1. `go run main.go {{validator_address}} 2` (shield)
1. After exec cmd file `output.json` will appear
1. Run `simd tx privacy benchmark {{path_to_output.json}}` to execute raw txs above
1. `go run main.go {{validator_address}} 3` (transfer)
1. Run `simd tx privacy benchmark {{path_to_output.json}}` to execute raw txs above
1. `go run main.go {{validator_address}} 4` (unshield)
1. Run `simd tx privacy benchmark {{path_to_output.json}}` to execute raw txs above
1. `go run main.go {{validator_address}} 5` (compose)
1. Run `simd tx privacy benchmark {{path_to_output.json}}` to execute raw txs above
