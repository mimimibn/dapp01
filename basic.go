package main

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {
	//这里使用测试网
	client, err := ethclient.Dial("https://ethereum-sepolia-rpc.publicnode.com")
	if err != nil {
		log.Fatal("init eth client is fail", err)
	}
	//获取指定块
	blockNumber := big.NewInt(5671744)
	//返回当前规范链的区块头。如果 `number` 为 nil，则返回最新的已知区块头。
	//newHeader, _ := client.HeaderByNumber(context.Background(), nil)
	//fmt.Println("number is : ", newHeader.Number)
	//获取当前链的区块头
	header, err := client.HeaderByNumber(context.Background(), blockNumber)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("header hash is : ", header.Hash().Hex())              //hash
	fmt.Println("header number : ", header.Number.Uint64())            //查询区块号
	fmt.Println("header date is : ", header.Time)                      //时间
	fmt.Println("header date is : ", time.Unix(int64(header.Time), 0)) //时间

	//chainID, err := client.ChainID(context.Background())
	//if err != nil {
	//	log.Fatal(err)
	//}

	//获取当前区块
	block, err := client.BlockByNumber(context.Background(), blockNumber)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("block hash is : ", block.Hash().Hex())        //hash
	fmt.Println("block number is : ", block.Number().Uint64()) //块id
	fmt.Println("block number is : ", block.Time())            //获取区块的时间

	//transactions := block.Transactions()
	//for _, tx := range transactions {
	//	sender, err := types.Sender(types.NewEIP155Signer(chainID), tx)
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//	fmt.Println(sender.Hex())
	//	break
	//}

	//获取当前区块的交易数量
	count, err := client.TransactionCount(context.Background(), block.Hash())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(count)
}
