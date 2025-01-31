package scripts

import (
	"context"
	"encoding/hex"
	"fmt"

	types "avail-alt-da-server/types"

	"github.com/availproject/avail-go-sdk/metadata"
	daPallet "github.com/availproject/avail-go-sdk/metadata/pallets/data_availability"
	SDK "github.com/availproject/avail-go-sdk/sdk"
)

func SubmitDataAndWatch(specs *types.AvailDASpecs, ctx context.Context, data []byte) (types.AvailBlockRef, error) {
	sdk, err := SDK.NewSDK(specs.ApiURL)
	if err != nil {
		panic(err)
	}

	accountId, err := metadata.NewAccountIdFromAddress("5GrwvaEF5zXb26Fz9rcQpDWS57CtERHpNehXCPcNoHGKutQY")
	if err != nil {
		return types.AvailBlockRef{}, fmt.Errorf("unable to create account id from address: %v, error: %w", "5GrwvaEF5zXb26Fz9rcQpDWS57CtERHpNehXCPcNoHGKutQY", err)
	}

	nonce, err := SDK.Account.Nonce(sdk.Client, accountId)
	if err != nil {
		return types.AvailBlockRef{}, fmt.Errorf("unable to get nonce for account id: %v, error: %w", accountId, err)
	}

	tx := sdk.Tx.DataAvailability.SubmitData(data)
	res, err := tx.ExecuteAndWatchInclusion(specs.KeyringPair, SDK.NewTransactionOptions().WithAppId(uint32(specs.AppID)))
	if err != nil {
		err = fmt.Errorf("cannot submit data: %w", err)
		return types.AvailBlockRef{}, err
	}
	if isSuc, err := res.IsSuccessful(); err != nil {
		err = fmt.Errorf("cannot check if data was submitted: %w", err)
		return types.AvailBlockRef{}, err
	} else if !isSuc {
		err = fmt.Errorf("data was not submitted: %w", err)
		return types.AvailBlockRef{}, err
	}

	println(fmt.Sprintf(`Block Hash: %v, Block Index: %v, Tx Hash: %v, Tx Index: %v`, res.BlockHash.ToHexWith0x(), res.BlockNumber, res.TxHash.ToHexWith0x(), res.TxIndex))
	events := res.Events.Unwrap()
	event := SDK.EventFindFirst(events, daPallet.EventDataSubmitted{}).Unwrap()

	return types.AvailBlockRef{BlockHash: res.BlockHash.ToHexWith0x(), Sender: hex.EncodeToString([]byte(specs.KeyringPair.SS58Address(42))), Nonce: int64(nonce), Commitment: event.DataHash.ToHexWith0x()}, nil

}
