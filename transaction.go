package main

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
	"math/big"
)

func main() {
	client, err := ethclient.Dial("https://ethereum-sepolia-rpc.publicnode.com")
	if err != nil {
		log.Fatal("init eth client is fail", err)
	}
	//加载私钥
	privateKey, err := crypto.HexToECDSA("") //todo 这里输入私钥
	if err != nil {
		log.Fatal("load hex key fail", err)
	}
	fmt.Println(privateKey)
	//获取公钥
	publicKey := privateKey.Public()
	fmt.Println(publicKey)
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}
	//根据公钥获取发送方的以太坊地址
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	//获取发送方的nonce，构建交易凭证
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	fmt.Println(nonce)                //下一笔交易的随机数
	sendValue := big.NewInt(10000000) //Wei
	gasLimit := uint64(21000)         //gas上限
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal("Failed to get gas price ", err)
	}
	toAddress := common.HexToAddress("") //todo 收款人的钱包地址
	//旧版交易
	//types.NewTransaction(nonce, toAddress, sendValue, gasLimit, gasPrice, nil)
	chainId, err := client.ChainID(context.Background())
	if err != nil {
		log.Fatal("get chainId is fail ", err)
	}
	//新版交易
	//这是分开写法
	//tx := types.NewTx(&types.LegacyTx{
	//	Nonce:    nonce,
	//	To:       &toAddress,
	//	Value:    sendValue,
	//	Gas:      gasLimit,
	//	GasPrice: gasPrice,
	//})
	//fmt.Println(tx)
	//types.SignTx(tx, types.NewEIP155Signer(chainId),privateKey)
	signTx, err := types.SignNewTx(privateKey, types.NewEIP155Signer(chainId), &types.LegacyTx{
		Nonce:    nonce,
		To:       &toAddress,
		Value:    sendValue,
		Gas:      gasLimit,
		GasPrice: gasPrice,
	})
	if err != nil {
		log.Fatal("Failed to sign transaction ", err)
	}
	err = client.SendTransaction(context.Background(), signTx)
	if err != nil {
		log.Fatal("Failed to send transaction", err)
	}
	fmt.Println("Transaction sent hash : ", signTx.Hash().Hex())
}
