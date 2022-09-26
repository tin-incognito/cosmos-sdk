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

simd add-genesis-account $MY_VALIDATOR_ADDRESS 100000000000stake
echo 12345678 |  simd gentx my_validator 100000000stake --chain-id my-test-chain
simd collect-gentxs
simd start