package main

import (
	"context"
	"log"

	"indexer/internal/db"
)

type SyncEvent struct {
	Pair  string
	Block int64
	R0    int64
	R1    int64
}

func main() {
	ctx := context.Background()

	conn, err := db.Connect()
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close(ctx)

	//  READ all sync events first
	rows, err := conn.Query(ctx, `
		SELECT
			pair_address,
			block_number,
			reserve0,
			reserve1
		FROM sync_events
		ORDER BY block_number
	`)
	if err != nil {
		log.Fatal(err)
	}

	var events []SyncEvent

	for rows.Next() {
		var ev SyncEvent
		if err := rows.Scan(&ev.Pair, &ev.Block, &ev.R0, &ev.R1); err != nil {
			log.Fatal(err)
		}
		events = append(events, ev)
	}
	rows.Close() 

	//  APPLY events to build reserves
	for _, ev := range events {
		_, err := conn.Exec(ctx, `
			INSERT INTO pair_reserves
				(pair_address, reserve0, reserve1, block_number)
			VALUES ($1,$2,$3,$4)
			ON CONFLICT (pair_address) DO UPDATE
			SET
				reserve0 = EXCLUDED.reserve0,
				reserve1 = EXCLUDED.reserve1,
				block_number = EXCLUDED.block_number
		`, ev.Pair, ev.R0, ev.R1, ev.Block)
		if err != nil {
			log.Fatal(err)
		}
	}

	log.Println("Reserves built")
}
