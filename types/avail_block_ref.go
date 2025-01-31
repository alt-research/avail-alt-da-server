package types

import (
	"encoding/json"
	"fmt"
	"time"

	"avail-alt-da-server/utils"

	SDK "github.com/availproject/avail-go-sdk/sdk"
	"github.com/ethereum/go-ethereum/log"
	"github.com/vedhavyas/go-subkey/v2"
)

type AvailBlockRef struct {
	BlockHash  string // Hash for block on avail chain
	Sender     string // sender address to filter extrinsic out sepecifically for this address
	Nonce      int64  // nonce to filter specific extrinsic
	Commitment string
}

func (a *AvailBlockRef) MarshalToBinary() ([]byte, error) {
	ref_bytes, err := json.Marshal(a)
	if err != nil {
		return []byte{}, fmt.Errorf("unable to covert the avail block referece into array of bytes and getting error:%w", err)
	}

	return []byte(ref_bytes), nil
}

func (a *AvailBlockRef) UnmarshalFromBinary(avail_blk_Ref []byte) error {
	err := json.Unmarshal(avail_blk_Ref, a)
	if err != nil {
		return fmt.Errorf("unable to convert avail_blk_Ref bytes to AvailBlockRef Struct and getting error:%w", err)
	}
	return nil
}

type AvailDASpecs struct {
	ApiURL      string
	Timeout     time.Duration
	AppID       int
	KeyringPair subkey.KeyPair
}

func NewAvailDASpecs(ApiURL string, AppID int, Seed string, Timeout time.Duration) (*AvailDASpecs, error) {

	AppID = utils.EnsureValidAppID(AppID)

	keyringPair, err := SDK.Account.NewKeyPair(Seed)
	if err != nil {
		log.Warn("⚠️ cannot create LeyPair: error:%w", err)
		return nil, err
	}

	return &AvailDASpecs{
		ApiURL:      ApiURL,
		Timeout:     Timeout,
		AppID:       AppID,
		KeyringPair: keyringPair,
	}, nil
}
