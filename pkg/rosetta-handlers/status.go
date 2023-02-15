package rosettastatus

import (
	"context"
	"rosetta_exporter/pkg/config"

	"github.com/coinbase/rosetta-sdk-go/fetcher"
	"github.com/coinbase/rosetta-sdk-go/types"
)

func GetStatus(cfg *config.Config) (*types.NetworkIdentifier, *types.NetworkStatusResponse, error) {
	ctx := context.Background()

	// Get Rosetta URL
	serverURL, err := cfg.GetRosettaUrl()
	if err != nil {
		return nil, nil, err
	}

	// Create a new fetcher
	newFetcher := fetcher.New(
		serverURL.String(),
	)

	// Initialize the fetcher's asserter
	primaryNetwork, networkStatus, error := newFetcher.InitializeAsserter(ctx, nil, "")
	if error != nil {
		return nil, nil, error.Err
	}

	// Set network values
	if cfg.GetNewtork() == "" {
		cfg.SetNetwork(primaryNetwork.Network)
	}

	// Set block chain values
	if cfg.GetBlockchain() == "" {
		cfg.SetBlockchain(primaryNetwork.Blockchain)
	}

	return primaryNetwork, networkStatus, nil
}
