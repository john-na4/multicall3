// Multicall3 example using go-ethereum: https://github.com/ethereum/go-ethereum
//
// This example demonstrates how to use Multicall3 to batch multiple Ethereum calls.
// It fetches DAI token information and balances for Vitalik's address, similar to the Rust example.
//
// Run `go run main.go` to run the example. Make sure to set the MAINNET_RPC_URL environment variable.
package main

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/joho/godotenv"
)

// Multicall3 ABI - only the functions we need
const multicall3ABI = `[
	{
		"inputs": [
			{
				"components": [
					{"internalType": "address", "name": "target", "type": "address"},
					{"internalType": "bytes", "name": "callData", "type": "bytes"}
				],
				"internalType": "struct Multicall3.Call[]",
				"name": "calls",
				"type": "tuple[]"
			}
		],
		"name": "aggregate",
		"outputs": [
			{"internalType": "uint256", "name": "blockNumber", "type": "uint256"},
			{"internalType": "bytes[]", "name": "returnData", "type": "bytes[]"}
		],
		"stateMutability": "payable",
		"type": "function"
	}
]`

// DAI ABI - only the functions we need
const daiABI = `[
	{"constant": true, "inputs": [], "name": "symbol", "outputs": [{"internalType": "string", "name": "", "type": "string"}], "payable": false, "stateMutability": "view", "type": "function"},
	{"constant": true, "inputs": [], "name": "decimals", "outputs": [{"internalType": "uint8", "name": "", "type": "uint8"}], "payable": false, "stateMutability": "view", "type": "function"},
	{"constant": true, "inputs": [{"internalType": "address", "name": "", "type": "address"}], "name": "balanceOf", "outputs": [{"internalType": "uint256", "name": "", "type": "uint256"}], "payable": false, "stateMutability": "view", "type": "function"}
]`

// Known contract addresses
var (
	multicall3Address = common.HexToAddress("0xcA11bde05977b3631167028862bE2a173976CA11")
	daiAddress        = common.HexToAddress("0x6B175474E89094C44Da98b954EedeAC495271d0F")
	vitalikAddress    = common.HexToAddress("0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045")
)

// Call represents a single call in the multicall
type Call struct {
	Target   common.Address `json:"target"`
	CallData []byte         `json:"callData"`
}

// AggregateResult represents the result of a multicall aggregate
type AggregateResult struct {
	BlockNumber *big.Int
	ReturnData  [][]byte
}

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Get RPC URL from environment
	rpcURL := os.Getenv("MAINNET_RPC_URL")
	if rpcURL == "" {
		log.Fatal("MAINNET_RPC_URL environment variable is required")
	}

	// Connect to Ethereum client
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}
	defer client.Close()

	// Parse ABIs
	multicallABI, err := abi.JSON(strings.NewReader(multicall3ABI))
	if err != nil {
		log.Fatalf("Failed to parse Multicall3 ABI: %v", err)
	}

	daiABIParsed, err := abi.JSON(strings.NewReader(daiABI))
	if err != nil {
		log.Fatalf("Failed to parse DAI ABI: %v", err)
	}

	// Prepare calls
	var calls []Call
	symbolData, err := daiABIParsed.Pack("symbol")
	if err != nil {
		log.Fatalf("Failed to pack symbol call: %v", err)
	}
	calls = append(calls, Call{Target: daiAddress, CallData: symbolData})

	decimalsData, err := daiABIParsed.Pack("decimals")
	if err != nil {
		log.Fatalf("Failed to pack decimals call: %v", err)
	}
	calls = append(calls, Call{Target: daiAddress, CallData: decimalsData})

	balanceData, err := daiABIParsed.Pack("balanceOf", vitalikAddress)
	if err != nil {
		log.Fatalf("Failed to pack balanceOf call: %v", err)
	}
	calls = append(calls, Call{Target: daiAddress, CallData: balanceData})

	// Pack the multicall
	multicallData, err := multicallABI.Pack("aggregate", calls)
	if err != nil {
		log.Fatalf("Failed to pack multicall: %v", err)
	}

	// Execute the multicall
	msg := ethereum.CallMsg{
		To:   &multicall3Address,
		Data: multicallData,
	}
	result, err := client.CallContract(context.Background(), msg, nil)
	if err != nil {
		log.Fatalf("Failed to execute multicall: %v", err)
	}

	// Unpack the result
	var aggregateResult AggregateResult
	err = multicallABI.UnpackIntoInterface(&aggregateResult, "aggregate", result)
	if err != nil {
		log.Fatalf("Failed to unpack multicall result: %v", err)
	}

	var symbol string
	err = daiABIParsed.UnpackIntoInterface(&symbol, "symbol", aggregateResult.ReturnData[0])
	if err != nil {
		log.Fatalf("Failed to unpack symbol: %v", err)
	}

	var decimals uint8
	err = daiABIParsed.UnpackIntoInterface(&decimals, "decimals", aggregateResult.ReturnData[1])
	if err != nil {
		log.Fatalf("Failed to unpack decimals: %v", err)
	}

	var daiBalance *big.Int
	err = daiABIParsed.UnpackIntoInterface(&daiBalance, "balanceOf", aggregateResult.ReturnData[2])
	if err != nil {
		log.Fatalf("Failed to unpack balance: %v", err)
	}
	// Convert DAI balance to human readable format
	daiBalanceFloat := new(big.Float).Quo(new(big.Float).SetInt(daiBalance), new(big.Float).SetInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(decimals)), nil)))

	// display results
	fmt.Printf("Block Number: %s\n", aggregateResult.BlockNumber.String())
	fmt.Printf("DAI Symbol: %s\n", symbol)
	fmt.Printf("DAI Decimals: %d\n", decimals)
	fmt.Printf("Vitalik's %s balance: %s\n", symbol, daiBalanceFloat.Text('f', 18))
}
