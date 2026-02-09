package main

import (
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

/*
========================
DATABASE MODELS
========================
*/

type RawLog struct {
	BlockNumber int64
	TxHash      []byte
	Data        []byte
	Topics      []byte
}

type SwapEvent struct {
	BlockNumber int64
	Amount0In   float64
	Amount0Out  float64
	Amount1In   float64
	Amount1Out  float64
}

type SyncEvent struct {
	BlockNumber int64
	Reserve0    float64
	Reserve1    float64
}

type PoolReserves struct {
	Reserve0 float64
	Reserve1 float64
}

/*
========================
MAIN
========================
*/

func main() {
	ctx := context.Background()

	// DATABASE_URL env variable or fallback
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://postgres:password@localhost:5432/indexer?sslmode=disable"
	}

	// Create DB pool
	pool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		log.Fatal("Failed to create pool:", err)
	}
	defer pool.Close()

	log.Println("Starting Day 14 replay (Sync is authoritative)...")

	/*
	========================
	LOAD CURRENT RESERVES
	========================
	*/
	reserves := PoolReserves{
		Reserve0: 0,
		Reserve1: 0,
	}

	err = pool.QueryRow(ctx, `
		SELECT reserve0, reserve1
		FROM pools
		WHERE pool_id = 1
	`).Scan(&reserves.Reserve0, &reserves.Reserve1)

	if err != nil {
		log.Println("No existing reserves found, starting from zero")
	}

	/*
	========================
	REPLAY RAW LOGS
	========================
	*/
	batchSize := int64(500)
	offset := int64(0)

	for {
		rows, err := pool.Query(ctx, `
			SELECT block_number, tx_hash, data, topics
			FROM logs
			ORDER BY block_number, tx_hash
			OFFSET $1 LIMIT $2
		`, offset, batchSize)
		if err != nil {
			log.Fatal("Failed to fetch logs:", err)
		}

		count := 0

		for rows.Next() {
			var raw RawLog
			if err := rows.Scan(
				&raw.BlockNumber,
				&raw.TxHash,
				&raw.Data,
				&raw.Topics,
			); err != nil {
				log.Println("Scan error:", err)
				continue
			}

			eventType, swap, sync := decodeEvent(raw)

			switch eventType {

			case "Swap":
				// Tentative reserve update
				reserves.Reserve0 = reserves.Reserve0 + swap.Amount0In - swap.Amount0Out
				reserves.Reserve1 = reserves.Reserve1 + swap.Amount1In - swap.Amount1Out

				log.Printf(
					"[Block %d] Swap applied tentatively",
					swap.BlockNumber,
				)

			case "Sync":
				// Authoritative overwrite
				reserves.Reserve0 = sync.Reserve0
				reserves.Reserve1 = sync.Reserve1

				log.Printf(
					"[Block %d] Sync applied — reserves overwritten (authoritative)",
					sync.BlockNumber,
				)
			}

			count++
		}

		rows.Close()

		if count == 0 {
			break
		}

		/*
		========================
		SAVE FINAL RESERVES
		========================
		*/
		_, err = pool.Exec(ctx, `
			UPDATE pools
			SET reserve0 = $1,
			    reserve1 = $2
			WHERE pool_id = 1
		`, reserves.Reserve0, reserves.Reserve1)

		if err != nil {
			log.Println("Failed to update reserves:", err)
		}

		offset += batchSize
		log.Printf("Processed %d logs, moving to next batch...\n", count)
	}

	log.Println("replay completed successfully!")
}

/*
========================
EVENT DECODER (MOCK)
========================
NOTE:
This is intentionally simple.
Replace with real ABI decoding later.
*/

func decodeEvent(raw RawLog) (string, SwapEvent, SyncEvent) {
	dataStr := string(raw.Data)

	// MOCK RULE:
	// If data == "SYNC" → Sync event
	// Else → Swap event
	if dataStr == "SYNC" {
		return "Sync",
			SwapEvent{},
			SyncEvent{
				BlockNumber: raw.BlockNumber,
				Reserve0:    1000,
				Reserve1:    500,
			}
	}

	return "Swap",
		SwapEvent{
			BlockNumber: raw.BlockNumber,
			Amount0In:   100,
			Amount0Out:  0,
			Amount1In:   0,
			Amount1Out:  50,
		},
		SyncEvent{}
}
