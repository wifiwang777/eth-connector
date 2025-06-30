package eth_connector

import (
	"context"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"golang.org/x/crypto/sha3"
	"math/big"
)

type Erc20 struct {
	*Ethereum
}

func NewErc20(e *Ethereum) *Erc20 {
	return &Erc20{
		Ethereum: e,
	}
}

func (e *Erc20) BalanceOf(contractAddress, account common.Address) (*big.Int, error) {
	transferFnSignature := []byte("balanceOf(address)")
	keccak256 := sha3.NewLegacyKeccak256()
	keccak256.Write(transferFnSignature)
	methodId := keccak256.Sum(nil)[:4]

	var data []byte
	data = append(data, methodId...)
	data = append(data, common.LeftPadBytes(account.Bytes(), 32)...)

	msg := ethereum.CallMsg{
		To:   &contractAddress,
		Data: data,
	}

	bytes, err := e.client.CallContract(context.Background(), msg, nil)
	if err != nil {
		return nil, err
	}
	balance := new(big.Int).SetBytes(bytes)
	return balance, nil
}

func (e *Erc20) Transfer(contractAddress, from, to common.Address, amount *big.Int) (*types.Transaction, error) {
	transferFnSignature := []byte("transfer(address,uint256)")
	erc20hash := sha3.NewLegacyKeccak256()
	erc20hash.Write(transferFnSignature)
	methodId := erc20hash.Sum(nil)[:4]

	var data []byte
	data = append(data, methodId...)
	data = append(data, common.LeftPadBytes(to.Bytes(), 32)...)
	data = append(data, common.LeftPadBytes(amount.Bytes(), 32)...)

	tx, err := e.NewTx(from, contractAddress, nil, data)
	if err != nil {
		return nil, err
	}

	return tx, nil
}

func (e *Erc20) Approve(contractAddress, from, spender common.Address, amount *big.Int) (*types.Transaction, error) {
	approveFnSignature := []byte("approve(address,uint256)")
	erc20hash := sha3.NewLegacyKeccak256()
	erc20hash.Write(approveFnSignature)
	methodId := erc20hash.Sum(nil)[:4]

	var data []byte
	data = append(data, methodId...)
	data = append(data, common.LeftPadBytes(spender.Bytes(), 32)...)
	data = append(data, common.LeftPadBytes(amount.Bytes(), 32)...)

	tx, err := e.NewTx(from, contractAddress, nil, data)
	if err != nil {
		return nil, err
	}
	return tx, nil
}

func (e *Erc20) Allowance(contractAddress, owner, spender common.Address) (*big.Int, error) {
	transferFnSignature := []byte("allowance(address,address)")
	keccak256 := sha3.NewLegacyKeccak256()
	keccak256.Write(transferFnSignature)
	methodId := keccak256.Sum(nil)[:4]

	var data []byte
	data = append(data, methodId...)
	data = append(data, common.LeftPadBytes(owner.Bytes(), 32)...)
	data = append(data, common.LeftPadBytes(spender.Bytes(), 32)...)
	msg := ethereum.CallMsg{
		To:   &contractAddress,
		Data: data,
	}

	bytes, err := e.client.CallContract(context.Background(), msg, nil)
	if err != nil {
		return nil, err
	}
	balance := new(big.Int).SetBytes(bytes)
	return balance, nil
}

func (e *Erc20) TransferFrom(contractAddress, from, sender, receiver common.Address, amount *big.Int) (*types.Transaction, error) {
	transferFnSignature := []byte("transferFrom(address,address,uint256)")
	erc20hash := sha3.NewLegacyKeccak256()
	erc20hash.Write(transferFnSignature)
	methodId := erc20hash.Sum(nil)[:4]

	var data []byte
	data = append(data, methodId...)
	data = append(data, common.LeftPadBytes(sender.Bytes(), 32)...)
	data = append(data, common.LeftPadBytes(receiver.Bytes(), 32)...)
	data = append(data, common.LeftPadBytes(amount.Bytes(), 32)...)

	tx, err := e.NewTx(from, contractAddress, nil, data)
	if err != nil {
		return nil, err
	}

	return tx, nil
}
