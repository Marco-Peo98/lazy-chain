package settings

import (
	"context"
	"fmt"
	"time"

	"github.com/algorand/go-algorand-sdk/client/v2/algod"
	"github.com/algorand/go-algorand-sdk/client/v2/indexer"
)

// NetworkManager owns the live connections and current selection.
type NetworkManager struct {
	algodClient    *algod.Client
	indexerClient  *indexer.Client
	currentNetwork NetworkInfo
	connected      bool
}

// NewNetworkManager creates a new network manager.
func NewNetworkManager() *NetworkManager {
	return &NetworkManager{connected: false}
}

// TestNetworkConnection tests connection to a network without storing it.
func (nm *NetworkManager) TestNetworkConnection(networkInfo NetworkInfo) error {
	algodClient, err := createAlgodClient(networkInfo)
	if err != nil {
		return fmt.Errorf("failed to create algod client: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if _, err = algodClient.Status().Do(ctx); err != nil {
		return fmt.Errorf("algod connection test failed: %w", err)
	}

	// Indexer optional
	if networkInfo.IndexerURL != "" {
		indexerClient, err := createIndexerClient(networkInfo)
		if err != nil {
			return fmt.Errorf("failed to create indexer client: %w", err)
		}
		if _, err = indexerClient.SearchForApplications().Limit(1).Do(ctx); err != nil {
			return fmt.Errorf("indexer connection test failed: %w", err)
		}
	}
	return nil
}

// ConnectToNetwork establishes connection to the specified network.
func (nm *NetworkManager) ConnectToNetwork(networkInfo NetworkInfo) error {
	algodClient, err := createAlgodClient(networkInfo)
	if err != nil {
		return fmt.Errorf("failed to connect to algod: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if _, err = algodClient.Status().Do(ctx); err != nil {
		return fmt.Errorf("failed to get network status: %w", err)
	}

	var indexerClient *indexer.Client
	if networkInfo.IndexerURL != "" {
		indexerClient, err = createIndexerClient(networkInfo)
		if err == nil {
			if _, err = indexerClient.SearchForApplications().Limit(1).Do(ctx); err != nil {
				fmt.Printf("Warning: indexer connection test failed: %v\n", err)
				indexerClient = nil
			}
		} else {
			fmt.Printf("Warning: failed to connect to indexer: %v\n", err)
		}
	}

	nm.algodClient = algodClient
	nm.indexerClient = indexerClient
	nm.currentNetwork = networkInfo
	nm.connected = true
	return nil
}

// GetNetworkStatus returns current network connection status.
func (nm *NetworkManager) GetNetworkStatus() (NetworkStatus, error) {
	if !nm.connected || nm.algodClient == nil {
		return NetworkStatus{}, fmt.Errorf("not connected to any network")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	status, err := nm.algodClient.Status().Do(ctx)
	if err != nil {
		return NetworkStatus{}, fmt.Errorf("failed to get network status: %w", err)
	}

	genesisResponse, err := nm.algodClient.GetGenesis().Do(ctx)
	if err != nil {
		return NetworkStatus{}, fmt.Errorf("failed to get genesis info: %w", err)
	}

	// Check indexer health if available
	indexerHealthy := false
	if nm.indexerClient != nil {
		_, err := nm.indexerClient.SearchForApplications().Limit(1).Do(ctx)
		indexerHealthy = (err == nil)
	}

	// Extract genesis info - original code used placeholders
	genesisID := "unknown"
	genesisHash := "unknown"
	if len(genesisResponse) > 0 {
		genesisID = "genesis-available"
		genesisHash = nm.currentNetwork.Name
	}

	return NetworkStatus{
		NetworkName:    nm.currentNetwork.Name,
		AlgodURL:       nm.currentNetwork.AlgodURL,
		IndexerURL:     nm.currentNetwork.IndexerURL,
		Connected:      true,
		LastRound:      status.LastRound,
		GenesisID:      genesisID,
		GenesisHash:    genesisHash,
		IndexerHealthy: indexerHealthy,
		Timestamp:      time.Now(),
	}, nil
}

func (nm *NetworkManager) IsConnected() bool                 { return nm.connected && nm.algodClient != nil }
func (nm *NetworkManager) GetCurrentNetwork() NetworkInfo    { return nm.currentNetwork }
func (nm *NetworkManager) GetAlgodClient() *algod.Client     { return nm.algodClient }
func (nm *NetworkManager) GetIndexerClient() *indexer.Client { return nm.indexerClient }

func (nm *NetworkManager) Disconnect() {
	nm.algodClient = nil
	nm.indexerClient = nil
	nm.connected = false
	nm.currentNetwork = NetworkInfo{}
}
