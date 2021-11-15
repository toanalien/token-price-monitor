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
	"math"
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
	providerUrl := os.Getenv("PROVIDER_URL")
	rClient, _ := ethclient.Dial(providerUrl)

	ins, _ := token.NewToken(common.HexToAddress(contract), rClient)
	reverses, _ := ins.GetReserves(&bind.CallOpts{})

	token0, _ := ins.Token0(&bind.CallOpts{})
	token1, _ := ins.Token1(&bind.CallOpts{})
	token0Contract, _ := token.NewToken(common.HexToAddress(token0.Hex()), rClient)
	token1Contract, _ := token.NewToken(common.HexToAddress(token1.Hex()), rClient)
	decimal0Uint, _ := token0Contract.Decimals(&bind.CallOpts{})
	decimal1Uint, _ := token1Contract.Decimals(&bind.CallOpts{})

	decimal0, _ := new(big.Float).SetInt64(int64(decimal0Uint)).Int64()
	decimal1, _ := new(big.Float).SetInt64(int64(decimal1Uint)).Int64()

	decimal0Pow10 := math.Pow(10, float64(decimal0))
	decimal1Pow10 := math.Pow(10, float64(decimal1))
	log.Println("decimals 0: ", decimal0)
	log.Println("decimals 1: ", decimal1)

	basePrice0 := new(big.Float).Quo(new(big.Float).SetInt(reverses.Reserve0), new(big.Float).SetFloat64(decimal0Pow10))
	basePrice1 := new(big.Float).Quo(new(big.Float).SetInt(reverses.Reserve1), new(big.Float).SetFloat64(decimal1Pow10))

	basePrice := new(big.Float).Quo(basePrice0, basePrice1)
	fmt.Println(fmt.Sprintf("Price 0/1: %f", basePrice))
	fmt.Println(fmt.Sprintf("Price 1/0: %f", new(big.Float).Quo(big.NewFloat(1), basePrice)))
	if err != nil {
		log.Fatal(err)
	}

	return nil
}
