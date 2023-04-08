package rosettahandlers

import (
	"context"
	"rosetta_exporter/pkg/config"

	"github.com/coinbase/rosetta-sdk-go/fetcher"
	"github.com/coinbase/rosetta-sdk-go/types"
)

type RosettaHandler struct {
	// PrimaryNetwork identifier for blockchain network and subnetworks
	PrimaryNetwork *types.NetworkIdentifier
	// fetcher request handler for retrieving data from chain
	fetcher *fetcher.Fetcher
}

// NewRosettaHandler constructor for rosetta connection and configuration handler
func NewRosettaHandler(cfg *config.Config) (*RosettaHandler, error) {
	ctx := context.Background()

	// initialize the rosetta handler
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

// GetStatus retrieve the status from rosetta /network/status endpoint
func (r *RosettaHandler) GetStatus() (*types.NetworkStatusResponse, error) {
	ctx := context.Background()

	// Initialize the fetcher's asserter
	networkStatus, e := r.fetcher.NetworkStatus(ctx, r.PrimaryNetwork, nil)
	if e != nil {
		return nil, e.Err
	}

	return networkStatus, nil
}

// GetBlock retrieve the block from the rosetta /network/block endpoint
func (r *RosettaHandler) GetBlock(blockIdentifier *types.BlockIdentifier) (*types.Block, error) {
	ctx := context.Background()

	// fetch the block with retry
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
