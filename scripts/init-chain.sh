if [ ! -f /root/.simapp/my_validator.info ]; then
  echo "Create my_validator key"
  echo -e "12345678\n12345678" | simd keys add my_validator
fi

MY_VALIDATOR_ADDRESS=$(echo 12345678 | simd keys show my_validator -a )
echo $MY_VALIDATOR_ADDRESS

if [ ! -f /root/.simapp/config/genesis.json ]; then
  simd init name --chain-id my-test-chain
else
  echo "Chain already init genesis"
fi

simd add-genesis-account $MY_VALIDATOR_ADDRESS 100000000000prv
echo 12345678 |  simd gentx my_validator 100000000prv --chain-id my-test-chain
simd collect-gentxs

(
cd ~/.simapp/config
sed -i 's/stake/prv/g' app.toml
sed -i 's/stake/prv/g' genesis.json
sed -i '/\[api\]/,+3 s/enable = false/enable = true/' app.toml
jq '.app_state.gov.voting_params.voting_period = "600s"' genesis.json > temp.json && mv temp.json genesis.json
jq '.app_state.mint.minter.inflation = "0.300000000000000000"' genesis.json > temp.json && mv temp.json genesis.json
)

simd start --log_level info
