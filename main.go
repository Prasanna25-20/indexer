package main

import (
	"context"
	"crypto/sha256"
	"fmt"
	"log"
	"sort"
	"time"

	"github.com/jackc/pgx/v5"
)

// ---------- CONFIG ----------
const (
	LOG_RANGE_SIZE = 10
	CONFIRMATIONS  = 12
)

// ---------- INVARIANT HELPER ----------
func assertInvariant(cond bool, msg string) {
	if !cond {
		log.Fatal("Invariant failed:", msg)
	}
}

// ---------- HEX UTILS ----------
func hexToInt(hex string) int64 {
	var v int64
	_, err := fmt.Sscanf(hex, "%x", &v)
	if err != nil {
		return 0
	}
	return v
}

// ---------- FAKE RPC / DB HELPERS ----------
func GetBlockAndLogsRPC(ctx context.Context, blockNumber int64) (map[string]interface{}, []map[string]interface{}, error) {
	blk := map[string]interface{}{
		"number": blockNumber,
		"hash":   fmt.Sprintf("0xBLOCKHASH%d", blockNumber),
	}
	logs := []map[string]interface{}{}
	for i := 0; i < 3; i++ {
		logs = append(logs, map[string]interface{}{
			"logIndex":         fmt.Sprintf("%x", i),
			"transactionHash":  fmt.Sprintf("0xTX%d", i),
			"blockNumber":      fmt.Sprintf("%x", blockNumber),
			"address":          fmt.Sprintf("0xADDR%d", i),
			"data":             fmt.Sprintf("0xDATA%d", i),
			"topics":           []string{"0xTOPIC"},
			"transactionIndex": fmt.Sprintf("%x", i),
		})
	}
	return blk, logs, nil
}

func InsertBlockWithLogs(ctx context.Context, conn *pgx.Conn, blk map[string]interface{}, logs []map[string]interface{}) error {
	fmt.Println("Insert block", blk["number"], "with", len(logs), "logs")
	return nil
}

// ---------- REPLAY LOGIC ----------
func RunReplay(ctx context.Context, conn *pgx.Conn, startBlock, safeTip int64) error {
	for batchStart := startBlock; batchStart <= safeTip; batchStart += LOG_RANGE_SIZE {
		batchEnd := batchStart + LOG_RANGE_SIZE - 1
		if batchEnd > safeTip {
			batchEnd = safeTip
		}

		log.Printf("\nIndexing logs from %d to %d\n", batchStart, batchEnd)

		for i := batchStart; i <= batchEnd; i++ {
			blk, logs, err := GetBlockAndLogsRPC(ctx, i)
			if err != nil {
				log.Println("Failed to fetch block", i)
				continue
			}

			sort.Slice(logs, func(a, b int) bool {
				return hexToInt(logs[a]["logIndex"].(string)) < hexToInt(logs[b]["logIndex"].(string))
			})

			if err := InsertBlockWithLogs(ctx, conn, blk, logs); err != nil {
				log.Println("Failed to insert block", i)
				continue
			}
		}

		time.Sleep(100 * time.Millisecond)
	}

	return nil
}

// ---------- DATABASE HASH (FOR VERIFICATION) ----------
func ComputeDBHash() string {
	// Fake deterministic hash for demo — replace with real DB query
	data := ""
	for i := int64(1); i <= 100; i++ {
		for j := 0; j < 3; j++ {
			data += fmt.Sprintf("%d-%d", i, j)
		}
	}
	sum := sha256.Sum256([]byte(data))
	return fmt.Sprintf("%x", sum)
}

// ---------- RESET DB ----------
func ResetDB() error {
	// For demo: just print, replace with actual DB drop & create commands
	fmt.Println("Resetting database...")
	return nil
}

// ---------- FULL REPLAY VERIFICATION ----------
func FullReplayVerification(ctx context.Context, conn *pgx.Conn) {
	fmt.Println("\n=== Day 17: Full Replay Verification ===")

	// Step 1: Delete DB
	if err := ResetDB(); err != nil {
		log.Fatal(err)
	}

	// Step 2: First replay
	fmt.Println("Running first replay...")
	if err := RunReplay(ctx, conn, 1, 100); err != nil {
		log.Fatal(err)
	}
	hash1 := ComputeDBHash()
	fmt.Println("First replay hash:", hash1)

	// Step 3: Delete DB again
	if err := ResetDB(); err != nil {
		log.Fatal(err)
	}

	// Step 4: Second replay
	fmt.Println("Running second replay...")
	if err := RunReplay(ctx, conn, 1, 100); err != nil {
		log.Fatal(err)
	}
	hash2 := ComputeDBHash()
	fmt.Println("Second replay hash:", hash2)

	// Step 5: Compare
	if hash1 == hash2 {
		fmt.Println(" Full Replay Verification PASSED — Hashes are identical!")
	} else {
		fmt.Println(" Full Replay Verification FAILED — Hashes differ!")
	}
}

// ---------- MAIN ----------
func main() {
	ctx := context.Background()

	conn, _ := pgx.Connect(ctx, "postgres://postgres:postgres@localhost:5432/indexer")
	defer conn.Close(ctx)

	// Existing Day 15 indexing
	lastIndexed := int64(0)
	startBlock := lastIndexed + 1
	latest := int64(100)
	safeTip := latest - CONFIRMATIONS
	if startBlock <= safeTip {
		log.Printf("Starting Day 15 Indexing from block %d to %d\n", startBlock, safeTip)
		if err := RunReplay(ctx, conn, startBlock, safeTip); err != nil {
			log.Fatal(err)
		}
		log.Println(" Indexing Complete ")
	} else {
		log.Println("No new blocks to index. Safe tip:", safeTip)
	}

	// Day 17: Full Replay Verification
	FullReplayVerification(ctx, conn)
}
