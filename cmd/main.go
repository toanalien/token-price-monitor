package main

import token_price_monitor "github.com/toanalien/token-price-monitor"

func main() {
	_ = token_price_monitor.Monitor(nil, token_price_monitor.PubSubMessage{})
}
