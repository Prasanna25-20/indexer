package state

import (
	"context"

	"github.com/jackc/pgx/v5"
)

func ApplySyncEvent(ctx context.Context, db *pgx.Conn, e SyncEvent) error {
	_, err := db.Exec(ctx, `
		INSERT INTO pair_reserves (
			pair_address, reserve0, reserve1, block_number, tx_hash
		)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (pair_address) DO UPDATE SET
			reserve0 = EXCLUDED.reserve0,
			reserve1 = EXCLUDED.reserve1,
			block_number = EXCLUDED.block_number,
			tx_hash = EXCLUDED.tx_hash
	`,
		e.PairAddress,
		e.Reserve0,
		e.Reserve1,
		e.BlockNumber,
		e.TxHash,
	)
	return err
}
