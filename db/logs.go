func InsertLogs(ctx context.Context, logs []RawLog) error {
	for _, l := range logs {
		_, err := db.Exec(ctx, `
			INSERT INTO logs (
				block_number,
				tx_hash,
				tx_index,
				log_index,
				address,
				topics,
				data
			)
			VALUES ($1,$2,$3,$4,$5,$6,$7)
			ON CONFLICT DO NOTHING
		`,
			l.BlockNumber,
			l.TxHash,
			l.TxIndex,
			l.LogIndex,
			l.Address,
			l.Topics,
			l.Data,
		)
		if err != nil {
			return err
		}
	}
	return nil
}
