package eth_connector

import (
	"context"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

const (
	RpcEndpointMainnet = "wss://ethereum-rpc.publicnode.com"
	RpcEndpointHolesky = "wss://ethereum-holesky-rpc.publicnode.com"

	USDTContractAddressMainnet = "0xdAC17F958D2ee523a2206206994597C13D831ec7"
	USDTContractAddressHolesky = "0xb27e39Fb20333aC358E7fB37a9994F44f1a7F66B"
	USDCContractAddressMainnet = "0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48"
	USDCContractAddressHolesky = "0xB6Fd173B7d71fae7413A7Aa880d4cbD57d29908D"
)

func TestBalance(t *testing.T) {
	client, err := ethclient.Dial(RpcEndpointHolesky)
	if err != nil {
		t.Error(err)
		return
	}
	erc20 := NewErc20(NewEthereum(client))
	contractAddress := common.HexToAddress(USDTContractAddressHolesky)
	account := common.HexToAddress("0x646D15cCC9157EE02a51747a6fd5D8B914F655F0")
	balance, err := erc20.BalanceOf(contractAddress, account)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(balance.String())
}

func TestTransfer(t *testing.T) {
	client, err := ethclient.Dial(RpcEndpointHolesky)
	if err != nil {
		t.Error(err)
		return
	}
	erc20 := NewErc20(NewEthereum(client))

	contractAddress := common.HexToAddress(USDCContractAddressHolesky)
	decimal := int64(18)

	from := common.HexToAddress("0x94b6D081b604953FE0720046d4D8023291A91656")
	to := common.HexToAddress("0x646D15cCC9157EE02a51747a6fd5D8B914F655F0")
	amount := big.NewInt(100)
	precision := new(big.Int).Exp(big.NewInt(10), big.NewInt(decimal), nil)
	amount = new(big.Int).Mul(amount, precision)
	tx, err := erc20.Transfer(contractAddress, from, to, amount)
	if err != nil {
		t.Error(err)
		return
	}

	hexKey := "YOUR_PRIVATE_KEY"
	privateKey, err := crypto.HexToECDSA(hexKey)
	signedTx, err := erc20.SignTx(tx, privateKey)
	if err != nil {
		t.Error(err)
		return
	}

	txId := signedTx.Hash().Hex()
	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		t.Error(err)
	}
	t.Logf("Transaction sent: %s", txId)
}

func TestApprove(t *testing.T) {
	client, err := ethclient.Dial(RpcEndpointHolesky)
	if err != nil {
		t.Error(err)
		return
	}
	erc20 := NewErc20(NewEthereum(client))

	base := big.NewInt(2)
	exponent := big.NewInt(256)
	twoPow256 := new(big.Int).Exp(base, exponent, nil)
	one := big.NewInt(1)
	unlimitedApproveAmount := new(big.Int).Sub(twoPow256, one)

	contractAddress := common.HexToAddress(USDTContractAddressHolesky)
	from := common.HexToAddress("0x646D15cCC9157EE02a51747a6fd5D8B914F655F0")
	spender := common.HexToAddress("0x38c108CEbf53edD0B025D6390fF7Eb473d98bAbE")
	tx, err := erc20.Approve(contractAddress, from, spender, unlimitedApproveAmount)
	if err != nil {
		t.Error(err)
		return
	}

	hexKey := "YOUR_PRIVATE_KEY"
	privateKey, err := crypto.HexToECDSA(hexKey)
	signedTx, err := erc20.SignTx(tx, privateKey)
	if err != nil {
		t.Error(err)
		return
	}

	txId := signedTx.Hash().Hex()
	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		t.Error(err)
	} else {
		t.Logf("Transaction sent: %s", txId)
	}
}

func TestAllowance(t *testing.T) {
	client, err := ethclient.Dial(RpcEndpointHolesky)
	if err != nil {
		t.Error(err)
		return
	}
	erc20 := NewErc20(NewEthereum(client))
	contractAddress := common.HexToAddress(USDTContractAddressHolesky)
	owner := common.HexToAddress("0x646D15cCC9157EE02a51747a6fd5D8B914F655F0")
	spender := common.HexToAddress("0x38c108CEbf53edD0B025D6390fF7Eb473d98bAbE")
	allowance, err := erc20.Allowance(contractAddress, owner, spender)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(allowance.String())
}

func TestTransferFrom(t *testing.T) {
	client, err := ethclient.Dial(RpcEndpointHolesky)
	if err != nil {
		t.Error(err)
		return
	}
	erc20 := NewErc20(NewEthereum(client))

	contractAddress := common.HexToAddress(USDTContractAddressHolesky)
	decimal := int64(18)

	from := common.HexToAddress("0x38c108CEbf53edD0B025D6390fF7Eb473d98bAbE")
	sender := common.HexToAddress("0x646D15cCC9157EE02a51747a6fd5D8B914F655F0")
	receiver := common.HexToAddress("0x38c108CEbf53edD0B025D6390fF7Eb473d98bAbE")
	amount := big.NewInt(100)
	precision := new(big.Int).Exp(big.NewInt(10), big.NewInt(decimal), nil)
	amount = new(big.Int).Mul(amount, precision)
	tx, err := erc20.TransferFrom(contractAddress, from, sender, receiver, amount)
	if err != nil {
		t.Error(err)
		return
	}

	hexKey := "YOUR_PRIVATE_KEY"
	privateKey, err := crypto.HexToECDSA(hexKey)
	signedTx, err := erc20.SignTx(tx, privateKey)
	if err != nil {
		t.Error(err)
		return
	}

	txId := signedTx.Hash().Hex()
	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		t.Error(err)
	} else {
		t.Logf("Transaction sent: %s", txId)
	}
}
