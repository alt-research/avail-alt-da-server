package scripts

import (
	"context"
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

	accountId, err := metadata.NewAccountIdFromAddress(specs.KeyringPair.SS58Address(42))

	if err != nil {
		return types.AvailBlockRef{}, fmt.Errorf("unable to create account id from address: %v, error: %w", specs.KeyringPair.SS58Address(42), err)
	}

	nonce, err := SDK.Account.Nonce(sdk.Client, accountId)
	if err != nil {
		return types.AvailBlockRef{}, fmt.Errorf("unable to get nonce for account id: %v, error: %w", accountId, err)
	}

	tx := sdk.Tx.DataAvailability.SubmitData(data)
	res, err := tx.ExecuteAndWatchInclusion(specs.KeyringPair, SDK.NewTransactionOptions().WithAppId(uint32(specs.AppID)))
	if err != nil {
		return types.AvailBlockRef{}, err
	}

	if isSuc, err := res.IsSuccessful(); err != nil {
		return types.AvailBlockRef{}, fmt.Errorf("cannot check if data was submitted: %w", err)
	} else if !isSuc {
		return types.AvailBlockRef{}, fmt.Errorf("data was not found in the block")
	}

	println(fmt.Sprintf(`Block Hash: %v, Block Index: %v, Tx Hash: %v, Tx Index: %v`, res.BlockHash.ToHexWith0x(), res.BlockNumber, res.TxHash.ToHexWith0x(), res.TxIndex))
	events := res.Events.Unwrap()
	event := SDK.EventFindFirst(events, daPallet.EventDataSubmitted{}).Unwrap()

	return types.AvailBlockRef{BlockHash: res.BlockHash.ToHexWith0x(), Sender: specs.KeyringPair.SS58Address(42), Nonce: int64(nonce), Commitment: event.DataHash.ToHexWith0x()}, nil

}
