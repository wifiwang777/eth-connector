package eth_connector

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"math/big"
)

type Ethereum struct {
	client *ethclient.Client
}

func NewEthereum(client *ethclient.Client) *Ethereum {
	return &Ethereum{
		client: client,
	}
}

func (e *Ethereum) NewTx(from, to common.Address, value *big.Int, data []byte) (*types.Transaction, error) {
	chainId, err := e.client.ChainID(context.Background())
	if err != nil {
		return nil, err
	}

	nonce, err := e.client.PendingNonceAt(context.Background(), from)
	if err != nil {
		return nil, err
	}

	gasLimit, err := e.client.EstimateGas(context.Background(), ethereum.CallMsg{
		From:  from,
		To:    &to,
		Value: value,
		Data:  data,
	})
	if err != nil {
		return nil, err
	}
	
	// maxPriorityFeePerGas
	gasTipCap, err := e.client.SuggestGasTipCap(context.Background())
	if err != nil {
		return nil, err
	}

	if gasTipCap.Cmp(big.NewInt(1)) < 0 {
		gasTipCap = big.NewInt(1)
	}

	header, err := e.client.HeaderByNumber(context.Background(), nil)
	if err != nil {
		return nil, err
	}

	baseFee := header.BaseFee
	if baseFee == nil {
		return nil, errors.New("base fee is required")
	}

	// Ensure base fee is at least 150% of the current base fee
	baseFee = new(big.Int).Mul(baseFee, big.NewInt(150))
	baseFee = new(big.Int).Div(baseFee, big.NewInt(100))

	// Calculate maxFeePerGas as 2 * baseFee + gasTipCap
	maxFeePerGas := new(big.Int).Add(baseFee, gasTipCap)

	tx := types.NewTx(&types.DynamicFeeTx{
		ChainID:   chainId,
		Nonce:     nonce,
		GasTipCap: gasTipCap,
		GasFeeCap: maxFeePerGas,
		Gas:       gasLimit,
		To:        &to,
		Value:     value,
		Data:      data,
	})
	return tx, nil
}

func (e *Ethereum) SignTx(tx *types.Transaction, privateKey *ecdsa.PrivateKey) (*types.Transaction, error) {
	signer := types.NewLondonSigner(tx.ChainId())
	digest := signer.Hash(tx)
	sign, err := crypto.Sign(digest.Bytes(), privateKey)
	if err != nil {
		return nil, err
	}
	signedTx, err := tx.WithSignature(signer, sign)
	if err != nil {
		return nil, err
	}
	return signedTx, nil
}
