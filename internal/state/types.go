package state

type SyncEvent struct {
	PairAddress string
	Reserve0    string
	Reserve1    string
	BlockNumber int64
	TxHash      string
}
