package state

import (
    "context"

    "github.com/jackc/pgx/v5"
)

type SyncEvent struct {
    Pair        []byte
    Token0      []byte
    Token1      []byte
    Reserve0    string
    Reserve1    string
    BlockNumber int64
    TxHash      []byte
    LogIndex    int
}

func FetchSyncEvents(ctx context.Context, db *pgx.Conn, fromBlock int64) ([]SyncEvent, error) {
    rows, err := db.Query(ctx, `
        SELECT
            pair_address,
            token0,
            token1,
            reserve0,
            reserve1,
            block_number,
            tx_hash,
            log_index
        FROM events
        WHERE event_name = 'Sync'
          AND block_number >= $1
        ORDER BY
            block_number,
            tx_index,
            log_index
    `, fromBlock)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var events []SyncEvent
    for rows.Next() {
        var ev SyncEvent
        if err := rows.Scan(
            &ev.Pair,
            &ev.Token0,
            &ev.Token1,
            &ev.Reserve0,
            &ev.Reserve1,
            &ev.BlockNumber,
            &ev.TxHash,
            &ev.LogIndex,
        ); err != nil {
            return nil, err
        }
        events = append(events, ev)
    }
    return events, nil
}

func ApplySyncEvent(ctx context.Context, db *pgx.Conn, ev SyncEvent) error {
    _, err := db.Exec(ctx, `
        INSERT INTO pair_reserves (
            pair_address,
            token0,
            token1,
            reserve0,
            reserve1,
            block_number,
            tx_hash,
            log_index
        ) VALUES ($1,$2,$3,$4,$5,$6,$7,$8)
        ON CONFLICT (pair_address) DO UPDATE SET
            reserve0 = EXCLUDED.reserve0,
            reserve1 = EXCLUDED.reserve1,
            block_number = EXCLUDED.block_number,
            tx_hash = EXCLUDED.tx_hash,
            log_index = EXCLUDED.log_index,
            updated_at = now()
    `,
        ev.Pair,
        ev.Token0,
        ev.Token1,
        ev.Reserve0,
        ev.Reserve1,
        ev.BlockNumber,
        ev.TxHash,
        ev.LogIndex,
    )
    return err
}
