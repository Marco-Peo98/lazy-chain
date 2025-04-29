package algo

import (
	"fmt"
	"strings"

	"github.com/algorand/go-algorand-sdk/v2/client/kmd"
	"github.com/algorand/go-algorand-sdk/v2/client/v2/algod"
	"github.com/algorand/go-algorand-sdk/v2/client/v2/indexer"
	"github.com/algorand/go-algorand-sdk/v2/crypto"
)

type NetworkConfig struct {
	NodeAddress   string
	NodeToken string
	NodePort  string
}

var (
	ALGOD_ADDRESS = "http://localhost"
	ALGOD_PORT    = "4001"
	ALGOD_URL     = ""
	ALGOD_TOKEN   = strings.Repeat("a", 64)

	INDEXER_ADDRESS = "http://localhost"
	INDEXER_PORT    = "8980"
	INDEXER_TOKEN   = strings.Repeat("a", 64)
	INDEXER_URL     = ""

	KMD_ADDRESS = "http://localhost"
	KMD_PORT    = "4002"
	KMD_TOKEN   = strings.Repeat("a", 64)
	KMD_URL     = ""

	KMD_WALLET_NAME     = "unencrypted-default-wallet"
	KMD_WALLET_PASSWORD = ""

	LORA_WALLET_NAME 		= "lora-dev"
	LORA_WALLET_PASSWORD	= ""
)

func GetAlgodClient(network NetworkConfig) (*algod.Client, error) {
	ALGOD_URL = fmt.Sprintf("%s%s", network.NodeAddress, network.NodePort)
	ALGOD_TOKEN = network.NodeToken
	algod, err := algod.MakeClient(
		ALGOD_URL,
		ALGOD_TOKEN,
	)
	return algod, err
}

func GetIndexerClient(network NetworkConfig) (*indexer.Client, error) {
	INDEXER_URL = fmt.Sprintf("%s%s", network.NodeAddress, network.NodePort)
	INDEXER_TOKEN = network.NodeToken
	indexer, err := indexer.MakeClient(
		INDEXER_URL,
		INDEXER_TOKEN,
	)
	return indexer, err
}

func GetKmdClient(network NetworkConfig) (kmd.Client, error) {
	KMD_URL = fmt.Sprintf("%s%s", network.NodeAddress, network.NodePort)
	KMD_TOKEN = network.NodeToken
	kmd, err := kmd.MakeClient(
		KMD_URL,
		KMD_TOKEN,
	)
	return kmd, err
}

func GetSandboxAccounts() ([]crypto.Account, error) {
	localnet := NetworkConfig{
		ALGOD_ADDRESS,
		ALGOD_TOKEN,
		ALGOD_PORT	}

	client, err := GetKmdClient(localnet)

	resp, err := client.ListWallets()
	if err != nil {
		return nil, fmt.Errorf("Failed to list wallets: %+v", err)
	}

	var walletId string
	for _, wallet := range resp.Wallets {
		if wallet.Name == KMD_WALLET_NAME {
			walletId = wallet.ID
		}
	}

	if walletId == "" {
		return nil, fmt.Errorf("No wallet named %s", KMD_WALLET_NAME)
	}

	whResp, err := client.InitWalletHandle(walletId, KMD_WALLET_PASSWORD)
	if err != nil {
		return nil, fmt.Errorf("Failed to init wallet handle: %+v", err)
	}

	addrResp, err := client.ListKeys(whResp.WalletHandleToken)
	if err != nil {
		return nil, fmt.Errorf("Failed to list keys: %+v", err)
	}

	var accts []crypto.Account
	for _, addr := range addrResp.Addresses {
		expResp, err := client.ExportKey(whResp.WalletHandleToken, KMD_WALLET_PASSWORD, addr)
		if err != nil {
			return nil, fmt.Errorf("Failed to export key: %+v", err)
		}

		acct, err := crypto.AccountFromPrivateKey(expResp.PrivateKey)
		if err != nil {
			return nil, fmt.Errorf("Failed to create account from private key: %+v", err)
		}

		accts = append(accts, acct)
	}

	return accts, nil
}

func GetLoraAccounts() ([]crypto.Account, error) {
	localnet := NetworkConfig{
		ALGOD_ADDRESS,
		ALGOD_TOKEN,
		ALGOD_PORT	}

	client, err := GetKmdClient(localnet)

	resp, err := client.ListWallets()
	if err != nil {
		return nil, fmt.Errorf("Failed to list wallets: %+v", err)
	}

	var walletId string
	for _, wallet := range resp.Wallets {
		if wallet.Name == LORA_WALLET_NAME {
			walletId = wallet.ID
		}
	}

	if walletId == "" {
		return nil, fmt.Errorf("No wallet named %s", LORA_WALLET_NAME)
	}

	whResp, err := client.InitWalletHandle(walletId, LORA_WALLET_PASSWORD)
	if err != nil {
		return nil, fmt.Errorf("Failed to init wallet handle: %+v", err)
	}

	addrResp, err := client.ListKeys(whResp.WalletHandleToken)
	if err != nil {
		return nil, fmt.Errorf("Failed to list keys: %+v", err)
	}

	var accts []crypto.Account
	for _, addr := range addrResp.Addresses {
		expResp, err := client.ExportKey(whResp.WalletHandleToken, LORA_WALLET_PASSWORD, addr)
		if err != nil {
			return nil, fmt.Errorf("Failed to export key: %+v", err)
		}

		acct, err := crypto.AccountFromPrivateKey(expResp.PrivateKey)
		if err != nil {
			return nil, fmt.Errorf("Failed to create account from private key: %+v", err)
		}

		accts = append(accts, acct)
	}

	return accts, nil
}