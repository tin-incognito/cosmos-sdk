echo "Init params ... "
MY_VALIDATOR_ADDRESS=$(echo 12345678 | simd keys show my_validator -a )
echo -e "12345678\n12345678" | simd keys add accA_N 2>&1 &> /dev/null
echo -e "12345678\n12345678" | simd keys add accB_N 2>&1 &> /dev/null
echo -e "12345678\n12345678" | simd keys add accC_N 2>&1 &> /dev/null
accA_N=$(echo 12345678 | simd keys show accA_N -a )
accB_N=$(echo 12345678 | simd keys show accB_N -a )
accC_N=$(echo 12345678 | simd keys show accC_N -a )

accA_Priv=112t8roMiguHp5SZxXpnoTHiV8A7XWkm79Y4hWwrqjtvjDH52H4y6z5RBt8HqqmKXSZ33CGjd2uM4ZiG32hqMzytboMA7zrBGscciRRVG4rH
accA_Pay=12scjGeftnVsH4Xa2CsRpkqctZaY77asQQMq84Gqx3itNRZaaQxfCDARXfrVZrsSK63pGPC2DzYwdEhLAjAyrUKErGeZfcL2v7HXXTVLee6Gwvr5NsruJRCqiHnQ9aaGsYGDKy8mgTzu1pJfPdv8
accB_Priv=112t8rqkMBoQDtSFggSfJyji71BxE2HV1kTSarohB7rUgy5qYMecCqvLoWM8vN3KMVr8RQskdQvJLSK1qducBPnG2v9xLFmkvyaKUNp97y8X
accB_Pay=12soNT5whtDTZSh3Twak286DmwiXAm19edjqkCX7kF9koickEaeCPkUAQUwHEzhkSn7YuBm1jdWR14t8Q1UeCBLycwwyawvw2iJkre572Y2ZAXFSPNkjfJhio6kJP1s6Gj3fEjrq4RL7YCpPfUWQ
accC_Priv=112t8rnq9oRPURjwuzyjCRbQqZk3iQbHCXRBLRGUp9EBjVdd7HMpRCSTWUrw7gKRNtqVThaygDiqCRBryeH35kznUGHT1FK9wMZq44myvC7v
accC_Pay=12sjT1YQTyB2yQSd325JAYR9QvNq44cQkoWjphiu5HnSrMdxUiLMZCq6B2fKhi87HfRVZAZHM3vc2yCjKwTcRAy1rDAevmswbBsqRY7as8rksFbAAVhUm9BTRwSQfaWXDh3F3AdWvnVRfSUMcaaT

echo "Send coin to A ... "
echo 12345678 | simd tx bank send $MY_VALIDATOR_ADDRESS $accA_N 5000prv --chain-id my-test-chain -y


checkbalance() {
  echo -e "\nChecking balances"
  sleep 5s
  echo "Balance of A (nonprivacy):"
  simd query bank balances $accA_N --chain-id my-test-chain | grep amount
  echo "Balance of A (privacy):" . `simd query privacy balance $accA_Priv`
  echo "Balance of B (nonprivacy):"
  simd query bank balances $accB_N --chain-id my-test-chain | grep amount
  echo "Balance of B (privacy):" . `simd query privacy balance $accB_Priv`
  echo "Balance of C (nonprivacy):"
  simd query bank balances $accC_N --chain-id my-test-chain | grep amount
  echo "Balance of C (privacy):" . `simd query privacy balance $accC_Priv`
}

checkbalance

echo 12345678 | simd tx privacy shield $accB_Pay 500 --from accA_N --chain-id my-test-chain -y

checkbalance

echo 12345678 | simd tx privacy transfer $accB_Priv $accC_Pay-300 0prv --chain-id my-test-chain -y

checkbalance

echo 12345678 | simd tx privacy unshield $accC_Priv $accC_N 200 0prv --chain-id my-test-chain -y

checkbalance