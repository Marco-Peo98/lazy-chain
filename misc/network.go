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


const (
	LOCALNET 	= "localnet"
	MAINNET 	= "mainnet"
	TESTNET 	= "testnet"
	FNET 		= "fnet"
)

var (
	LOCALNET_ALGOD_ADDRESS = "http://localhost"
	LOCALNET_ALGOD_PORT    = "4001"
	LOCALNET_ALGOD_URL     = ""
	LOCALNET_ALGOD_TOKEN   = strings.Repeat("a", 64)

	LOCALNET_INDEXER_ADDRESS = "http://localhost"
	LOCALNET_INDEXER_PORT    = "8980"
	LOCALNET_INDEXER_TOKEN   = strings.Repeat("a", 64)
	LOCALNET_INDEXER_URL     = ""

	MAINNET_ALGOD_ADDRESS	= "https://mainnet-api.algonode.cloud/"
	MAINNET_ALGOD_PORT		= "443"
	MAINNET_ALGOD_TOKEN		= ""

	MAINNET_INDEXER_ADDRESS = "https://mainnet-idx.algonode.cloud/"
	MAINNET_INDEXER_PORT	= "443"
	MAINNET_INDEXER_TOKEN	= ""

	TESTNET_ALGOD_ADDRESS	= "https://testnet-api.algonode.cloud/"
	TESTNET_ALGOD_PORT		= "443"
	TESTNET_ALGOD_TOKEN		= ""

	KMD_ADDRESS = "http://localhost"
	KMD_PORT    = "4002"
	KMD_TOKEN   = strings.Repeat("a", 64)
	KMD_URL     = ""

	KMD_WALLET_NAME     = "unencrypted-default-wallet"
	KMD_WALLET_PASSWORD = ""

	LORA_WALLET_NAME 		= "lora-dev"
	LORA_WALLET_PASSWORD	= ""
)

func (c *AlgoClient) ConnectNetwork(name string) {
	switch name {
	case LOCALNET: 
		netCfg := NetworkConfig {
			LOCALNET_ALGOD_ADDRESS,
			LOCALNET_ALGOD_PORT,
			LOCALNET_ALGOD_TOKEN,
		}
		rest, err := GetAlgodClient(netCfg)
		if err != nil {
			fmt.Println("did not connect, problem at connection with network", name)
		}
		c.algod = rest
		break
	case MAINNET:
		netCfg := NetworkConfig {
			MAINNET_ALGOD_ADDRESS,
			MAINNET_ALGOD_PORT,
			MAINNET_ALGOD_TOKEN,
		}
		rest, err := GetAlgodClient(netCfg)
		if err != nil {
			fmt.Println("did not connect, problem at connection with network", name)
		}
		c.algod = rest
		break
	case TESTNET:
		netCfg := NetworkConfig {
			TESTNET_ALGOD_ADDRESS,
			TESTNET_ALGOD_PORT,
			TESTNET_ALGOD_TOKEN,
		}
		rest, err := GetAlgodClient(netCfg)
		if err != nil {
			fmt.Println("did not connect, problem at connection with network", name)
		}
		c.algod = rest
		break
	}

}

func GetAlgodClient(network NetworkConfig) (*algod.Client, error) {
	url := fmt.Sprintf("%s%s", network.NodeAddress, network.NodePort)
	token := network.NodeToken
	algod, err := algod.MakeClient(
		url,
		token,
	)
	return algod, err
}

func GetIndexerClient(network NetworkConfig) (*indexer.Client, error) {
	url := fmt.Sprintf("%s%s", network.NodeAddress, network.NodePort)
	token := network.NodeToken
	indexer, err := indexer.MakeClient(
		url,
		token,
	)
	return indexer, err
}

func GetKmdClient(network NetworkConfig) (kmd.Client, error) {
	url := fmt.Sprintf("%s%s", network.NodeAddress, network.NodePort)
	token := network.NodeToken
	kmd, err := kmd.MakeClient(
		url,
		token,
	)
	return kmd, err
}

func GetKmdAccounts() ([]crypto.Account, error) {
	localnet := NetworkConfig{
		LOCALNET_ALGOD_ADDRESS,
		LOCALNET_ALGOD_TOKEN,
		LOCALNET_ALGOD_PORT	}

	client, err := GetKmdClient(localnet)

	resp, err := client.ListWallets()
	if err != nil {
		return nil, fmt.Errorf("Failed to list wallets: %+v", err)
	}

	var walletId []string = make([]string, len(resp.Wallets))
	for i, wallet := range resp.Wallets {
		walletId[i] = wallet.ID
	}

	if len(walletId) == 0 {
		return nil, fmt.Errorf("No wallet named %s", KMD_WALLET_NAME)
	}

	
	var accts []crypto.Account

	for _, id := range walletId {
		whResp, err := client.InitWalletHandle(id, KMD_WALLET_PASSWORD)
		if err != nil {
			return nil, fmt.Errorf("Failed to init wallet handle: %+v", err)
		}

		addrResp, err := client.ListKeys(whResp.WalletHandleToken)
		if err != nil {
			return nil, fmt.Errorf("Failed to list keys: %+v", err)
		}

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
	}

	return accts, nil
}
