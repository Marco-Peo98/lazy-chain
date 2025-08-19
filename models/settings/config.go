package settings

import (
	"encoding/json"
	"fmt"
	"os"
)

// NetworkInfo holds RPC endpoints and tokens.
type NetworkInfo struct {
	Name         string `json:"name"`
	AlgodURL     string `json:"algod_url"`
	AlgodPort    string `json:"algod_port"`
	AlgodToken   string `json:"algod_token"`
	IndexerURL   string `json:"indexer_url"`
	IndexerPort  string `json:"indexer_port"`
	IndexerToken string `json:"indexer_token"`
}

// Config represents the configuration settings for the application.
type Config struct {
	Network        string        `json:"network"`
	WalletAddr     string        `json:"wallet_addr"`
	CustomNetworks []NetworkInfo `json:"custom_networks"`
}

func ConfigPath() string {
	home, _ := os.UserHomeDir()
	return fmt.Sprintf("%s/.lazy-chain/config.json", home)
}

func LoadConfig() (Config, error) {
	path := ConfigPath()
	buff, err := os.ReadFile(path)
	if err != nil {
		// Default config
		return Config{
			Network:        "testnet",
			WalletAddr:     "",
			CustomNetworks: []NetworkInfo{},
		}, nil
	}

	var cfg Config
	if err := json.Unmarshal(buff, &cfg); err != nil {
		return Config{}, err
	}

	// Ensure slice is non-nil
	if cfg.CustomNetworks == nil {
		cfg.CustomNetworks = []NetworkInfo{}
	}
	return cfg, nil
}

func SaveConfig(cfg Config) error {
	// Ensure directory exists
	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get user home directory: %w", err)
	}
	configDir := fmt.Sprintf("%s/.lazy-chain", home)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Ensure non-nil
	if cfg.CustomNetworks == nil {
		cfg.CustomNetworks = []NetworkInfo{}
	}

	// Pretty JSON
	buff, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(ConfigPath(), buff, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}
	return nil
}
