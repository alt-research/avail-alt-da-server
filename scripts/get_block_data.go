package scripts

import (
	"fmt"

	"avail-alt-da-server/types"

	"github.com/availproject/avail-go-sdk/metadata"
	"github.com/availproject/avail-go-sdk/primitives"

	SDK "github.com/availproject/avail-go-sdk/sdk"
)

func GetBlockExtrinsicData(specs types.AvailDASpecs, avail_blk_ref types.AvailBlockRef) ([]byte, error) {

	Hash := avail_blk_ref.BlockHash
	Address := avail_blk_ref.Sender
	Nonce := avail_blk_ref.Nonce

	avail_blk, err := fetchBlock(specs.ApiURL, Hash)
	if err != nil {
		panic(fmt.Errorf("cannot fetch block: %w", err))
	}

	return extractExtrinsic(Address, Hash, Nonce, avail_blk)
}

func fetchBlock(apiURL string, Hash string) (SDK.Block, error) {

	sdk, err := SDK.NewSDK(apiURL)
	if err != nil {
		return SDK.Block{}, fmt.Errorf("cannot create sdk: %w", err)
	}

	blockHash, err := primitives.NewBlockHashFromHexString(Hash)
	if err != nil {
		return SDK.Block{}, fmt.Errorf("unable to convert string hash into types.hash, error:%w", err)
	}
	block, err := SDK.NewBlock(sdk.Client, blockHash)
	if err != nil {
		return SDK.Block{}, fmt.Errorf("unable to create block, error:%w", err)
	}

	return block, nil
}

func extractExtrinsic(Address string, Hash string, Nonce int64, avail_blk SDK.Block) ([]byte, error) {
	accountId, err := metadata.NewAccountIdFromAddress(Address)
	if err != nil {
		return nil, fmt.Errorf("unable to create account id from address: %v, error: %w", Address, err)
	}

	for _, blob := range avail_blk.Block.Extrinsics {

		ext_Nonce := blob.Signed.Unwrap().Nonce

		if sameSignature(&blob, accountId) && ext_Nonce == uint32(Nonce) {
			extrinsic, ok := SDK.NewDataSubmission(&blob)
			if !ok {
				return []byte{}, fmt.Errorf("unable to create data submission from extrinsic")
			}
			return extrinsic.Data, nil
		}
	}

	return []byte{}, fmt.Errorf("didn't find any extrinsic data for address:%v in block having hash:%v", Address, Hash)
}

func sameSignature(tx *primitives.DecodedExtrinsic, accountId metadata.AccountId) bool {
	txAccountId := tx.Signed.Unwrap().Address.Id.Unwrap()
	if accountId.Value != txAccountId {
		return false
	}

	return true
}
