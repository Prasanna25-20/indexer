package db

import "database/sql"

func GetLastIndexedBlock(db *sql.DB) (int64, error) {
	var last int64
	err := db.QueryRow(
		`SELECT last_block FROM indexer_state WHERE id = TRUE`,
	).Scan(&last)
	return last, err
}

func UpdateLastIndexedBlock(db *sql.DB, blockNumber int64) error {
	_, err := db.Exec(
		`UPDATE indexer_state SET last_block = $1 WHERE id = TRUE`,
		blockNumber,
	)
	return err
}
