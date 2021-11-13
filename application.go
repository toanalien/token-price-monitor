package token_price_monitor

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/joho/godotenv"
	token "github.com/toanalien/token-price-monitor/abi"
	"log"
	"math/big"
	"os"
)

type PubSubMessage struct {
	Data []byte `json:"data"`
}

func Monitor(ctx context.Context, m PubSubMessage) error {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}
	contract := os.Getenv("CONTRACT")
	rClient, _ := ethclient.Dial("https://bsc-dataseed.binance.org/")
	ins, err := token.NewToken(common.HexToAddress(contract), rClient)
	reverses, err := ins.GetReserves(&bind.CallOpts{})
	basePrice := new(big.Float).Quo(new(big.Float).SetInt(reverses.Reserve0), new(big.Float).SetInt(reverses.Reserve1))
	fmt.Println(fmt.Sprintf("Price 0/1: %f", basePrice))
	fmt.Println(fmt.Sprintf("Price 1/0: %f", new(big.Float).Quo(big.NewFloat(1), basePrice)))
	if err != nil {
		log.Fatal(err)
	}
	return nil
}
