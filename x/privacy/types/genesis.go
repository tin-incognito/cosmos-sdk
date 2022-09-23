package types

import (
	"fmt"
)

// DefaultIndex is the default capability global index
const DefaultIndex uint64 = 1

// DefaultGenesis returns the default Capability genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		SerialNumberList:   []SerialNumber{},
		OutputCoinList:     []OutputCoin{},
		CommitmentList:     []Commitment{},
		TokenList:          []Token{},
		OnetimeAddressList: []OnetimeAddress{},
		OTACoinList:        []OTACoin{},
		OutputCoinLength:   nil,
		// this line is used by starport scaffolding # genesis/types/default
		Params: DefaultParams(),
	}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	// Check for duplicated index in serialNumber
	serialNumberIndexMap := make(map[string]struct{})

	for _, elem := range gs.SerialNumberList {
		index := string(SerialNumberKey(elem.Index))
		if _, ok := serialNumberIndexMap[index]; ok {
			return fmt.Errorf("duplicated index for serialNumber")
		}
		serialNumberIndexMap[index] = struct{}{}
	}
	// Check for duplicated index in outputCoin
	outputCoinIndexMap := make(map[string]struct{})

	for _, elem := range gs.OutputCoinList {
		index := string(OutputCoinKey(elem.Index))
		if _, ok := outputCoinIndexMap[index]; ok {
			return fmt.Errorf("duplicated index for outputCoin")
		}
		outputCoinIndexMap[index] = struct{}{}
	}
	// Check for duplicated index in token
	tokenIndexMap := make(map[string]struct{})

	for _, elem := range gs.TokenList {
		index := string(TokenKey(elem.Index))
		if _, ok := tokenIndexMap[index]; ok {
			return fmt.Errorf("duplicated index for token")
		}
		tokenIndexMap[index] = struct{}{}
	}
	// Check for duplicated index in onetimeAddress
	onetimeAddressIndexMap := make(map[string]struct{})

	for _, elem := range gs.OnetimeAddressList {
		index := string(OnetimeAddressKey(elem.Index))
		if _, ok := onetimeAddressIndexMap[index]; ok {
			return fmt.Errorf("duplicated index for onetimeAddress")
		}
		onetimeAddressIndexMap[index] = struct{}{}
	}
	// Check for duplicated index in oTACoin
	oTACoinIndexMap := make(map[string]struct{})

	for _, elem := range gs.OTACoinList {
		index := string(OTACoinKey(elem.Index))
		if _, ok := oTACoinIndexMap[index]; ok {
			return fmt.Errorf("duplicated index for oTACoin")
		}
		oTACoinIndexMap[index] = struct{}{}
	}
	// this line is used by starport scaffolding # genesis/types/validate

	return gs.Params.Validate()
}
