package abis

import (
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"log"
	"math"
	"math/big"
	"os"
)

func GetReservesPrice(contract string, backend bind.ContractBackend) (*big.Float, *big.Float) {
	contractMultiCall := os.Getenv("CONTRACT_MULTICALL")
	mtCall, _ := NewMulticall(common.HexToAddress(contractMultiCall), backend)

	token0 := methodToken("token0")
	token1 := methodToken("token1")
	getReserves := methodGetReserves()

	callToken0 := Struct0{}
	callToken1 := Struct0{}
	callReserves := Struct0{}
	callReserves.Target = common.HexToAddress(contract)
	callReserves.CallData = getReserves.ID
	callToken0.Target = common.HexToAddress(contract)
	callToken0.CallData = token0.ID
	callToken1.Target = common.HexToAddress(contract)
	callToken1.CallData = token1.ID

	tokens, _ := mtCall.Aggregate(&bind.CallOpts{}, []Struct0{callToken0, callToken1, callReserves})
	log.Println(tokens.ReturnData)

	token0Address := tokens.ReturnData[0]
	token1Address := tokens.ReturnData[1]

	callDecimals0 := Struct0{}
	callDecimals1 := Struct0{}
	callDecimals0.Target = common.BytesToAddress(token0Address)
	callDecimals0.CallData = methodDecimals().ID
	callDecimals1.Target = common.BytesToAddress(token1Address)
	callDecimals1.CallData = methodDecimals().ID
	decimals, _ := mtCall.Aggregate(&bind.CallOpts{}, []Struct0{callDecimals0, callDecimals1})
	log.Println(decimals)

	uint256, _ := abi.NewType("uint256", "", nil)
	decimal0 := abi.ReadInteger(uint256, decimals.ReturnData[0]).(*big.Int)
	decimal1 := abi.ReadInteger(uint256, decimals.ReturnData[1]).(*big.Int)
	log.Println("decimals 0: ", decimal0)
	log.Println("decimals 1: ", decimal1)

	decimal0Pow10 := math.Pow(10, float64(decimal0.Int64()))
	decimal1Pow10 := math.Pow(10, float64(decimal1.Int64()))

	type Reverses struct {
		Reserve0           *big.Int
		Reserve1           *big.Int
		BlockTimestampLast uint32
	}
	abiPair := abi.ABI{
		Methods: map[string]abi.Method{
			"getReserves": getReserves},
	}
	var reverses Reverses
	if err := abiPair.UnpackIntoInterface(&reverses, "getReserves", tokens.ReturnData[2]); err != nil {
		log.Println(err)
	}
	log.Println(reverses)
	basePrice0 := new(big.Float).Quo(new(big.Float).SetInt(reverses.Reserve0), new(big.Float).SetFloat64(decimal0Pow10))
	basePrice1 := new(big.Float).Quo(new(big.Float).SetInt(reverses.Reserve1), new(big.Float).SetFloat64(decimal1Pow10))

	basePrice := new(big.Float).Quo(basePrice0, basePrice1)

	return basePrice, new(big.Float).Quo(big.NewFloat(1), basePrice)
}

func methodToken(token string) abi.Method {
	address, _ := abi.NewType("address", "", nil)
	return abi.NewMethod(
		token,
		token,
		abi.Function,
		"view",
		false,
		false,
		[]abi.Argument{},
		[]abi.Argument{{"address", address, false}})
}

func methodDecimals() abi.Method {
	uint256, _ := abi.NewType("uint256", "", nil)
	return abi.NewMethod(
		"decimals",
		"decimals",
		abi.Function,
		"view",
		false,
		false,
		[]abi.Argument{},
		[]abi.Argument{{"decimals", uint256, false}})
}

func methodGetReserves() abi.Method {
	uint256, _ := abi.NewType("uint256", "", nil)
	return abi.NewMethod(
		"getReserves",
		"getReserves",
		abi.Function,
		"view",
		false,
		false,
		[]abi.Argument{},
		[]abi.Argument{
			{"_reserve0", uint256, false},
			{"_reserve1", uint256, false},
			{"_blockTimestampLast", uint256, false},
		})
}
