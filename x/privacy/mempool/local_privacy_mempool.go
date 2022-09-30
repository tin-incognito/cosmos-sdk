package mempool

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/types"
	common2 "github.com/cosmos/cosmos-sdk/x/privacy/common"
	"github.com/cosmos/cosmos-sdk/x/privacy/repos"
	"github.com/cosmos/cosmos-sdk/x/privacy/repos/coin"
	"github.com/incognitochain/go-incognito-sdk-v2/common"
)

/*
Privacy mempool
- Support check double spend in local mempool
- Support creating tx with unspent coin (not in current DB and mempool)
*/

type PrivacyMempool struct {
	coinMap map[string]common2.Hash
	msgMap  map[common2.Hash]PrivacyMessage
}

type PrivacyMessage interface {
	GetProof() repos.PaymentProof
	GetFee() uint64
}

func NewPrivacyMempool() *PrivacyMempool {
	return &PrivacyMempool{
		make(map[string]common2.Hash),
		make(map[common2.Hash]PrivacyMessage),
	}
}

func (m *PrivacyMempool) IsDoubleSpendWithCurrentMempool(newMsg types.Msg) (bool, PrivacyMessage) {
	if !newMsg.IsPrivacy() {
		return false, nil
	}
	prf := newMsg.(PrivacyMessage).GetProof()

	isDoubleSpend := false
	for _, iCoin := range prf.InputCoins() {
		key := fmt.Sprintf("%v-%v", common.PRVCoinID.String(), string(iCoin.GetKeyImage().ToBytesS()))
		if h, ok := m.coinMap[key]; ok {
			isDoubleSpend = true
			if oldMsg, ok := m.msgMap[h]; (ok) && (oldMsg != nil) {
				if newMsg.(PrivacyMessage).GetFee() > oldMsg.GetFee() {
					return true, oldMsg
				}
			}
		}
	}

	for _, iCoin := range prf.OutputCoins() {
		key := fmt.Sprintf("%v-%v", common.PRVCoinID.String(), string(iCoin.GetKeyImage().ToBytesS()))
		if h, ok := m.coinMap[key]; ok {
			isDoubleSpend = true
			if oldMsg, ok := m.msgMap[h]; (ok) && (oldMsg != nil) {
				if newMsg.(PrivacyMessage).GetFee() > oldMsg.GetFee() {
					return true, oldMsg
				}
			}
		}
	}

	return isDoubleSpend, nil
}

func (m *PrivacyMempool) trackCoin(newMsg PrivacyMessage) {
	prf := newMsg.GetProof()
	h := common2.HashH(prf.Bytes())

	for _, iCoin := range prf.InputCoins() {
		key := fmt.Sprintf("%v-%v", common.PRVCoinID.String(), string(iCoin.GetKeyImage().ToBytesS()))
		m.coinMap[key] = h
	}

	for _, iCoin := range prf.OutputCoins() {
		key := fmt.Sprintf("%v-%v", common.PRVCoinID.String(), string(iCoin.GetKeyImage().ToBytesS()))
		m.coinMap[key] = h
	}

	m.msgMap[h] = newMsg.(PrivacyMessage)
}

func (m *PrivacyMempool) unTrackCoin(oldMsg PrivacyMessage) {
	prf := oldMsg.(PrivacyMessage).GetProof()
	h := common2.HashH(prf.Bytes())
	for _, iCoin := range prf.InputCoins() {
		key := fmt.Sprintf("%v-%v", common.PRVCoinID.String(), string(iCoin.GetKeyImage().ToBytesS()))
		delete(m.coinMap, key)
	}

	for _, iCoin := range prf.OutputCoins() {
		key := fmt.Sprintf("%v-%v", common.PRVCoinID.String(), string(iCoin.GetKeyImage().ToBytesS()))
		delete(m.coinMap, key)
	}
	delete(m.msgMap, h)
}

func (m *PrivacyMempool) AddToPool(newMsg types.Msg) bool {
	doubleSpend, oldMsg := m.IsDoubleSpendWithCurrentMempool(newMsg)

	if !doubleSpend {
		m.trackCoin(newMsg.(PrivacyMessage))
		return true
	}

	if oldMsg != nil {
		// remove old message coin
		m.unTrackCoin(oldMsg)
		m.trackCoin(newMsg.(PrivacyMessage))
		return true
	}

	return false
}

func (m *PrivacyMempool) RemoveFromPool(oldMsg types.Msg) error {
	m.unTrackCoin(oldMsg.(PrivacyMessage))
	return nil
}

func (m *PrivacyMempool) HasCoin(coin coin.Coin) bool {
	key := fmt.Sprintf("%v-%v", common.PRVCoinID.String(), string(coin.GetKeyImage().ToBytesS()))
	if _, ok := m.coinMap[key]; ok {
		return true
	}
	return false
}
