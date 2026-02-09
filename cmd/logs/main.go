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
		`INSERT INTO events VALUES
		 (17000001,'tx1',0,'pair1','SYNC','0xdata'),
		 (17000002,'tx2',0,'pair1','SYNC','0xdata')`)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Logs indexed")
}
