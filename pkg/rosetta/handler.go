package rosettahandlers

import (
	"context"
	"rosetta_exporter/pkg/config"

	"github.com/coinbase/rosetta-sdk-go/fetcher"
	"github.com/coinbase/rosetta-sdk-go/types"
)

type RosettaHandler struct {
	PrimaryNetwork *types.NetworkIdentifier
	fetcher        *fetcher.Fetcher
}

func NewRosettaHandler(cfg *config.Config) (*RosettaHandler, error) {
	ctx := context.Background()

	rosettaHandler := RosettaHandler{}

	// Get Rosetta URL
	serverURL, err := cfg.GetRosettaUrl()
	if err != nil {
		return nil, err
	}

	// Create a new fetcher
	rosettaHandler.fetcher = fetcher.New(
		serverURL.String(),
	)

	// Initialize the fetcher's asserter
	primaryNetwork, _, e := rosettaHandler.fetcher.InitializeAsserter(ctx, nil, "")
	if e != nil {
		return nil, e.Err
	}

	// Set primary network
	rosettaHandler.PrimaryNetwork = primaryNetwork

	return &rosettaHandler, nil
}

func (r *RosettaHandler) GetStatus() (*types.NetworkStatusResponse, error) {
	ctx := context.Background()

	// Initialize the fetcher's asserter
	networkStatus, e := r.fetcher.NetworkStatus(ctx, r.PrimaryNetwork, nil)
	if e != nil {
		return nil, e.Err
	}

	return networkStatus, nil
}

func (r *RosettaHandler) GetBlock(blockIdentifier *types.BlockIdentifier) (*types.Block, error) {
	ctx := context.Background()

	block, e := r.fetcher.BlockRetry(
		ctx,
		r.PrimaryNetwork,
		types.ConstructPartialBlockIdentifier(
			blockIdentifier,
		),
	)
	if e != nil {
		return nil, e.Err
	}

	return block, nil
}
