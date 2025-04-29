package algo

import (
	"context"
	"fmt"
	"time"

	"github.com/algorand/go-algorand-sdk/client/v2/algod"
	"github.com/algorand/go-algorand-sdk/client/v2/indexer"
	"github.com/algorand/go-algorand-sdk/crypto"
	"github.com/algorand/go-algorand-sdk/future"
	"github.com/algorand/go-algorand-sdk/types"
	"github.com/algorand/go-algorand-sdk/v2/mnemonic"
)

// AlgoClient contiene le connessioni e l'account
type AlgoClient struct {
	algod   *algod.Client
	indexer *indexer.Client
	account *crypto.Account // signer (privato)
}

// NewClient inizializza il client per algod e indexer
func NewClient(algodURL, algodToken, indexerURL string) (*AlgoClient, error) {
	algodClient, err := algod.MakeClient(algodURL, algodToken)
	if err != nil {
		return nil, fmt.Errorf("algod init error: %w", err)
	}

	indexerClient, err := indexer.MakeClient(indexerURL, "")
	if err != nil {
		return nil, fmt.Errorf("indexer init error: %w", err)
	}

	return &AlgoClient{
		algod:   algodClient,
		indexer: indexerClient,
	}, nil
}

// SetAccountFromMnemonic carica un account dal mnemonic
func (c *AlgoClient) SetAccountFromMnemonic(mn string) error {
	k, err := mnemonic.ToPrivateKey(mn)

	if err != nil {
		return fmt.Errorf("invalid mnemonic: %w", err)
	}

	
	recovered, err := crypto.AccountFromPrivateKey(k)
	if err != nil {
		return fmt.Errorf("invalid mnemonic: %w", err)
	}

	c.account = &crypto.Account{
		PublicKey:	recovered.PublicKey,
		PrivateKey: k,
		Address:    recovered.Address,
	}
	return nil
}

// GetAccountBalance restituisce il saldo in microAlgos
func (c *AlgoClient) GetAccountBalance(address string) (uint64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	info, err := c.algod.AccountInformation(address).Do(ctx)
	if err != nil {
		return 0, err
	}
	return info.Amount, nil
}

// SendAlgos invia Algos a un destinatario
func (c *AlgoClient) SendAlgos(to string, amount uint64, note string) (string, error) {
	if c.account == nil {
		return "", fmt.Errorf("signer not set")
	}

	fromAddr := c.account.Address.String()

	_, err := types.DecodeAddress(to)
	if err != nil {
		return "", fmt.Errorf("invalid recipient address: %w", err)
	}

	params, err := c.algod.SuggestedParams().Do(context.Background())
	if err != nil {
		return "", err
	}

	txn, err := future.MakePaymentTxn(
		fromAddr,
		to,
		amount,
		nil,
		note,
		params,
	)
	if err != nil {
		return "", err
	}

	txID, signedTxn, err := crypto.SignTransaction(c.account.PrivateKey, txn)
	if err != nil {
		return "", err
	}

	_, err = c.algod.SendRawTransaction(signedTxn).Do(context.Background())
	if err != nil {
		return "", err
	}

	return txID, nil
}

func (c *AlgoClient) CreateAsset(name, unitName string, total uint64, decimals uint32) (string, error) {
	if c.account == nil {
		return "", fmt.Errorf("signer not set")
	}

	fromAddr := c.account.Address.String()

	params, err := c.algod.SuggestedParams().Do(context.Background())
	if err != nil {
		return "", err
	}

	txn, err := future.MakeAssetCreateTxn(
		fromAddr,
		nil,
		params,
		total,
		decimals,
		false,
		fromAddr, // Manager
		fromAddr, // Reserve
		fromAddr, // Freeze
		fromAddr, // Clawback
		unitName,
		name,
		"",
		"",
	)
	if err != nil {
		return "", err
	}

	txID, signedTxn, err := crypto.SignTransaction(c.account.PrivateKey, txn)
	if err != nil {
		return "", err
	}

	_, err = c.algod.SendRawTransaction(signedTxn).Do(context.Background())
	if err != nil {
		return "", err
	}

	return txID, nil
}

func (c *AlgoClient) SendAsset(to string, assetID uint64, amount uint64) (string, error) {
	if c.account == nil {
		return "", fmt.Errorf("signer not set")
	}

	fromAddr := c.account.Address.String()

	_, err := types.DecodeAddress(to)
	if err != nil {
		return "", fmt.Errorf("invalid recipient address: %w", err)
	}

	params, err := c.algod.SuggestedParams().Do(context.Background())
	if err != nil {
		return "", err
	}

	txn, err := future.MakeAssetTransferTxn(
		fromAddr,
		to,
		0,      // closeRemainderTo
		nil,    // note
		amount, // amount
		params,
		assetID,
	)
	if err != nil {
		return "", err
	}

	txID, signedTxn, err := crypto.SignTransaction(c.account.PrivateKey, txn)
	if err != nil {
		return "", err
	}

	_, err = c.algod.SendRawTransaction(signedTxn).Do(context.Background())
	if err != nil {
		return "", err
	}

	return txID, nil
}

func (c *AlgoClient) OptInAsset(assetID uint64) (string, error) {
	if c.account == nil {
		return "", fmt.Errorf("signer not set")
	}

	fromAddr := c.account.Address.String()

	params, err := c.algod.SuggestedParams().Do(context.Background())
	if err != nil {
		return "", err
	}

	txn, err := future.MakeAssetAcceptanceTxn(
		fromAddr,
		nil,
		params,
		assetID,
	)
	if err != nil {
		return "", err
	}

	txID, signedTxn, err := crypto.SignTransaction(c.account.PrivateKey, txn)
	if err != nil {
		return "", err
	}

	_, err = c.algod.SendRawTransaction(signedTxn).Do(context.Background())
	if err != nil {
		return "", err
	}

	return txID, nil
}
