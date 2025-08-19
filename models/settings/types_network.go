package settings

import "time"

// NetworkStatus represents the current network status.
type NetworkStatus struct {
	NetworkName    string    `json:"network_name"`
	AlgodURL       string    `json:"algod_url"`
	IndexerURL     string    `json:"indexer_url"`
	Connected      bool      `json:"connected"`
	LastRound      uint64    `json:"last_round"`
	GenesisID      string    `json:"genesis_id"`
	GenesisHash    string    `json:"genesis_hash"`
	IndexerHealthy bool      `json:"indexer_healthy"`
	Timestamp      time.Time `json:"timestamp"`
}
