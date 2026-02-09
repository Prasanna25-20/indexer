package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {
	// -------------------------------
	// STEP 1: RPC URL
	// -------------------------------
	rpcURL := "https://eth-mainnet.g.alchemy.com/v2/jsL7QNrtDFZydcjcfskuX"
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		log.Fatal("Failed to connect:", err)
	}

	// -------------------------------
	// STEP 2: Load UniswapV2Pair ABI
	// -------------------------------
	abiFile, err := os.ReadFile("UniswapV2Pair.json")
	if err != nil {
		log.Fatal("ABI file not found")
	}

	uniswapABI, err := abi.JSON(strings.NewReader(string(abiFile)))
	if err != nil {
		log.Fatal(err)
	}

	// -------------------------------
	// STEP 3: Uniswap V2 Pair address
	// -------------------------------
	pairAddress := common.HexToAddress("0xB4e16d0168e52d35CaCD2c6185b44281Ec28C9Dc")

	// -------------------------------
	// STEP 4: Determine latest N blocks
	// -------------------------------
	latestBlock, err := client.BlockNumber(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	N := int64(50) // number of recent blocks to fetch
	startBlock := int64(latestBlock) - N + 1
	if startBlock < 0 {
		startBlock = 0
	}
	endBlock := int64(latestBlock)

	fmt.Printf("Fetching logs from block %d to %d in 10-block batches...\n", startBlock, endBlock)

	// -------------------------------
	// STEP 5: Fetch in 10-block batches
	// -------------------------------
	batchSize := int64(10) // Free-tier limit
	for b := startBlock; b <= endBlock; b += batchSize {
		from := big.NewInt(b)
		to := big.NewInt(b + batchSize - 1)
		if to.Int64() > endBlock {
			to = big.NewInt(endBlock)
		}

		fmt.Printf("Fetching blocks %d → %d\n", from.Int64(), to.Int64())

		query := ethereum.FilterQuery{
			Addresses: []common.Address{pairAddress},
			FromBlock: from,
			ToBlock:   to,
		}

		logs, err := client.FilterLogs(context.Background(), query)
		if err != nil {
			log.Println("Error fetching logs:", err)
			continue
		}

		// -------------------------------
		// Decode Swap + Sync events
		// -------------------------------
		for _, vLog := range logs {
			switch vLog.Topics[0].Hex() {
			case uniswapABI.Events["Swap"].ID.Hex():
				var swapEvent struct {
					Sender     common.Address
					Amount0In  *big.Int
					Amount1In  *big.Int
					Amount0Out *big.Int
					Amount1Out *big.Int
					To         common.Address
				}
				if err := uniswapABI.UnpackIntoInterface(&swapEvent, "Swap", vLog.Data); err != nil {
					log.Println("Swap decode error:", err)
					continue
				}
				fmt.Println("=== Swap Event ===")
				printJSON(swapEvent)

			case uniswapABI.Events["Sync"].ID.Hex():
				var syncEvent struct {
					Reserve0 *big.Int
					Reserve1 *big.Int
				}
				if err := uniswapABI.UnpackIntoInterface(&syncEvent, "Sync", vLog.Data); err != nil {
					log.Println("Sync decode error:", err)
					continue
				}
				fmt.Println("=== Sync Event ===")
				printJSON(syncEvent)
			}
		}
	}
}

// Helper: pretty print struct as JSON
func printJSON(v interface{}) {
	b, _ := json.MarshalIndent(v, "", "  ")
	fmt.Println(string(b))
}
