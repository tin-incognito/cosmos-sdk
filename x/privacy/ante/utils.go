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
	Proofs             *simplelru.LRU
	TxValidateSanity   *simplelru.LRU
	TxValidateByItself *simplelru.LRU
	TxValidateByDb     *simplelru.LRU
}

func NewCache() *Cache {
	proofCache, err := simplelru.NewLRU(6*1024*1024, nil)
	if err != nil {
		panic(err)
	}
	txValidateSanityCache, err := simplelru.NewLRU(6*1024*1024, nil)
	if err != nil {
		panic(err)
	}
	txValidateByItselfCache, err := simplelru.NewLRU(6*1024*1024, nil)
	if err != nil {
		panic(err)
	}
	txValidateByDbCache, err := simplelru.NewLRU(6*1024*1024, nil)
	if err != nil {
		panic(err)
	}
	return &Cache{
		Proofs:             proofCache,
		TxValidateSanity:   txValidateSanityCache,
		TxValidateByItself: txValidateByItselfCache,
		TxValidateByDb:     txValidateByDbCache,
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

func (c *Cache) GetTxValidateSanity(key common.Hash) (*TxCache, error) {
	val, found := c.TxValidateSanity.Get(key.String())
	if !found {
		return nil, fmt.Errorf("Cannot find TxValidateSanity from cache with key %s", key.String())
	}
	res, ok := val.(*TxCache)
	if !ok {
		return nil, fmt.Errorf("TxValidateSanity value from cache with key %s is invalid format", key.String())
	}
	return res, nil
}

func (c *Cache) SetTxValidateSanity(key common.Hash, value *TxCache) error {
	c.TxValidateSanity.Add(key.String(), value)
	return nil
}

func (c *Cache) GetTxValidateByItself(key common.Hash) (*TxCache, error) {
	val, found := c.TxValidateByItself.Get(key.String())
	if !found {
		return nil, fmt.Errorf("Cannot find TxValidateByItself from cache with key %s", key.String())
	}
	res, ok := val.(*TxCache)
	if !ok {
		return nil, fmt.Errorf("TxValidateByItself value from cache with key %s is invalid format", key.String())
	}
	return res, nil
}

func (c *Cache) SetTxValidateByItself(key common.Hash, value *TxCache) error {
	c.TxValidateByItself.Add(key.String(), value)
	return nil
}

func (c *Cache) GetTxValidateByDb(key common.Hash) (*TxCache, error) {
	val, found := c.TxValidateByDb.Get(key.String())
	if !found {
		return nil, fmt.Errorf("Cannot find TxValidateByDb from cache with key %s", key.String())
	}
	res, ok := val.(*TxCache)
	if !ok {
		return nil, fmt.Errorf("TxValidateByDb value from cache with key %s is invalid format", key.String())
	}
	return res, nil
}

func (c *Cache) SetTxValidateByDb(key common.Hash, value *TxCache) error {
	c.TxValidateByDb.Add(key.String(), value)
	return nil
}

type TxCache struct {
	Fee uint64
}
