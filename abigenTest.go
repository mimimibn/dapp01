package main

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
	"main/count"
	"math/big"
)

func main() {
	//部署
	deployContract("") //私钥
	//测试
	test("", "") //私钥,合约地址
}
func test(hexKey string, contract string) {
	client, err := ethclient.Dial("https://ethereum-sepolia-rpc.publicnode.com")
	if err != nil {
		log.Fatal("init eth client is fail", err)
	}
	countContract, err := count.NewCount(common.HexToAddress(contract), client)
	if err != nil {
		log.Fatal(err)
	}
	chainId, err := client.NetworkID(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	log.Println("NetworkID is ", chainId)
	chainIds, err := client.ChainID(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	log.Println("chainIds is : ", chainIds)
	privateKey, err := crypto.HexToECDSA(hexKey)
	if err != nil {
		log.Fatal(err)
	}
	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainId)
	if err != nil {
		log.Fatal(err)
	}
	addOne, err := countContract.AddOne(auth)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("addOne is ", addOne)
	opts := &bind.CallOpts{Context: context.Background()}
	num, err := countContract.GetCount(opts)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("num is ", num)
}
func deployContract(hexKey string) {
	client, err := ethclient.Dial("https://ethereum-sepolia-rpc.publicnode.com")
	if err != nil {
		log.Fatal("init eth client is fail", err)
	}
	//defer client.Close()
	privateKey, err := crypto.HexToECDSA(hexKey)
	if err != nil {
		log.Fatal(err)
	}
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatal(err)
	}
	//获取实时gas价格
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal("Failed to get gas price ", err)
	}
	//这里试试client.NetworkID，和client.ChainId
	chainId, err := client.NetworkID(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainId)
	if err != nil {
		log.Fatal(err)
	}
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)     // in wei
	auth.GasLimit = uint64(300000) // in units
	auth.GasPrice = gasPrice
	address, tx, c, err := count.DeployCount(auth, client)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(address.Hex())
	fmt.Println(tx.Hash().Hex())
	_ = c
}
