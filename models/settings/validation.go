package settings

import "fmt"

// ValidateNetworkConfig validates network configuration before connection.
func ValidateNetworkConfig(networkInfo NetworkInfo) error {
	if networkInfo.Name == "" {
		return fmt.Errorf("network name cannot be empty")
	}
	if networkInfo.AlgodURL == "" {
		return fmt.Errorf("algod URL cannot be empty")
	}
	if networkInfo.AlgodPort == "" {
		return fmt.Errorf("algod port cannot be empty")
	}
	if networkInfo.IndexerURL != "" && networkInfo.IndexerPort == "" {
		return fmt.Errorf("indexer port required when indexer URL is provided")
	}
	return nil
}
