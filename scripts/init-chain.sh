simd keys add my_validator
MY_VALIDATOR_ADDRESS=$(simd keys show my_validator -a)
simd init name --chain-id my-test-chain
simd add-genesis-account  $MY_VALIDATOR_ADDRESS 100000000000stake
simd gentx my_validator 100000000stake --chain-id my-test-chain
simd collect-gentxs
simd start