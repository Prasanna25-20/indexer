package main

import (
	"context"
	"log"

	"indexer/internal/db"
)

func main() {
	ctx := context.Background()

	conn, err := db.Connect()
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close(ctx)

	rows, err := conn.Query(ctx,
		`SELECT pair_address, reserve0, reserve1 FROM pair_reserves`)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
	 hookup := ""
	 var r0, r1 int64

	 if err := rows.Scan(&hookup, &r0, &r1); err != nil {
		 log.Fatal(err)
	 }

		log.Println("STATE", hookup, r0, r1)
	}

	log.Println("state derived ONLY from events")
}
