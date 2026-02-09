
func GetBlockByNumber(ctx context.Context, number int64) (*Block, error) {
    row := db.QueryRow(ctx,
        `SELECT number, hash, parent_hash, finalized
         FROM blocks
         WHERE number = $1`, number)

    var b Block
    err := row.Scan(&b.Number, &b.Hash, &b.ParentHash, &b.Finalized)
    if err != nil {
        return nil, err
    }
    return &b, nil
}
