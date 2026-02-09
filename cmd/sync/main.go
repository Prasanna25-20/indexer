package main

import (
	"context"
	"log"
	"indexer/internal/db"
)

func main() {
	conn, _ := db.Connect()
	defer conn.Close(context.Background())

	_, err := conn.Exec(context.Background(),
		`INSERT INTO sync_events VALUES
		 ('pair1',17000001,100,200),
		 ('pair1',17000002,150,250)`)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Sync events decoded")
}
