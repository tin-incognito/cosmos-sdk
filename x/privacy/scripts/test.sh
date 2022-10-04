MY_VALIDATOR_ADDRESS=$(echo 12345678 | simd keys show my_validator -a )
simd tx privacy shield 12skvRi6rzvj8UwUhxcr3xe8Z5YDQyaYvDf7e3FQUhyDMNKLgf1MyXh6GnWPbdvGH898PU1duHYzBJK3Qs7dx75VtttJXjm8aadp3ozDM4XnohccCi4dEdgte8o8n6ff29RHnqE3zcaLirTTcM25 100 --from $MY_VALIDATOR_ADDRESS --chain-id my-test-chain
