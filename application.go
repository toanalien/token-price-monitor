package token_price_monitor

import (
	"context"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/joho/godotenv"
	"github.com/toanalien/token-price-monitor/abi"
	"log"
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

	price01, price10 := abis.GetReservesPrice(contract, rClient)
	log.Println(price01, price10)
	g
	return nil
}
