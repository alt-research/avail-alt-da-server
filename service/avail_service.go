package avail

import (
	"context"
	"fmt"

	"time"

	"avail-alt-da-server/scripts"
	"avail-alt-da-server/types"

	"github.com/ethereum/go-ethereum/log"
)

type AvailService struct {
	Seed    string              `json:"seed"`
	ApiURL  string              `json:"api_url"`
	AppID   int                 `json:"app_id"`
	Timeout time.Duration       `json:"timeout"`
	Specs   *types.AvailDASpecs `json:"availDASpecs"`
	log     log.Logger
}

func NewAvailService(apiURL string, seed string, appID int, timeout time.Duration, log log.Logger) *AvailService {

	availSpecs, err := types.NewAvailDASpecs(apiURL, appID, seed, timeout)
	if err != nil {
		panic("failed avail initialisation")
	}

	return &AvailService{
		Seed:    seed,
		ApiURL:  apiURL,
		AppID:   appID,
		Timeout: timeout,
		Specs:   availSpecs,
		log:     log,
	}
}

func (s *AvailService) Get(ctx context.Context, comm []byte) ([]byte, error) {
	avail_blk_ref := types.AvailBlockRef{}
	err := avail_blk_ref.UnmarshalFromBinary(comm)
	if err != nil {
		s.log.Error("failed to unmarshal the ethereum tx data to avail block reference", "error", err)
		return []byte{}, fmt.Errorf("failed to unmarshal the ethereum tx data to avail block reference, error: %w", err)
	}

	input, err := scripts.GetBlockExtrinsicData(*s.Specs, avail_blk_ref, s.log)

	if err != nil {
		s.log.Error("failed to get block extrinsic data", "error", err)
		return []byte{}, fmt.Errorf("failed to get block extrinsic data: %w", err)
	}

	return input, nil
}

func (s *AvailService) Put(ctx context.Context, value []byte) ([]byte, error) {

	if len(value) >= 512000 {
		return nil, fmt.Errorf("the length of input cannot be greater than 512kb")
	}

	avail_Blk_Ref, err := scripts.SubmitDataAndWatch(s.Specs, ctx, value, s.log)

	if err != nil {
		s.log.Error("cannot submit data", "error", err)
		return nil, fmt.Errorf("cannot submit data:%w", err)
	}

	comm, err := avail_Blk_Ref.MarshalToBinary()

	if err != nil {
		s.log.Error("cannot get the binary form of avail block reference", "error", err)
		return nil, fmt.Errorf("cannot get the binary form of avail block reference:%w", err)
	}

	return comm, nil
}
