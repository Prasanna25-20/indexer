package state

import (
	"context"

	"github.com/jackc/pgx/v5"
)

func FetchSyncEvents(ctx context.Context, db *pgx.Conn, fromBlock int64) ([]SyncEvent, error) {
	rows, err := db.Query(ctx, `
		SELECT pair_address, reserve0, reserve1, block_number, tx_hash
		FROM sync_events
		WHERE block_number >= $1
		ORDER BY block_number, log_index
	`, fromBlock)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []SyncEvent
	for rows.Next() {
		var e SyncEvent
		if err := rows.Scan(
			&e.PairAddress,
			&e.Reserve0,
			&e.Reserve1,
			&e.BlockNumber,
			&e.TxHash,
		); err != nil {
			return nil, err
		}
		events = append(events, e)
	}
	return events, nil
}
