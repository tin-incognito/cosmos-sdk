package ante

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/privacy/common"
	"github.com/cosmos/cosmos-sdk/x/privacy/models"
	"github.com/cosmos/cosmos-sdk/x/privacy/repos"
	"github.com/cosmos/cosmos-sdk/x/privacy/types"
	"github.com/hashicorp/golang-lru/simplelru"
)

func IsMintTx(tx sdk.Tx) (bool, error) {
	msg := tx.GetMsgs()[0]
	switch msg := msg.(type) {
	case *types.MsgPrivacyData:
		return msg.TxType == models.TxMintType, nil
	default:
		errMsg := fmt.Sprintf("unrecognized %s message type: %T", types.ModuleName, msg)
		return false, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, errMsg)
	}
}

func IsPrivacyTx(tx sdk.Tx) (bool, error) {
	msgs := tx.GetMsgs()
	if len(msgs) != 1 {
		return false, nil
	}
	msg := tx.GetMsgs()[0]
	switch msg := msg.(type) {
	case *types.MsgPrivacyData:
		return true, nil
	default:
		errMsg := fmt.Sprintf("unrecognized %s message type: %T", types.ModuleName, msg)
		return false, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, errMsg)
	}
}

type Cache struct {
	Proofs *simplelru.LRU
}

func NewCache() *Cache {
	proofCache, err := simplelru.NewLRU(6*1024*1024, nil)
	if err != nil {
		panic(err)
	}
	return &Cache{
		Proofs: proofCache,
	}
}

func (c *Cache) AddProof(key common.Hash, proof *repos.PaymentProof) error {
	c.Proofs.Add(key.String(), proof)
	return nil
}

func (c *Cache) GetProof(key common.Hash) (*repos.PaymentProof, error) {
	val, found := c.Proofs.Get(key.String())
	if !found {
		return nil, fmt.Errorf("Cannot find proof from cache with key %s", key.String())
	}
	res, ok := val.(*repos.PaymentProof)
	if !ok {
		return nil, fmt.Errorf("Proof value from cache with key %s is invalid format", key.String())
	}
	return res, nil
}
