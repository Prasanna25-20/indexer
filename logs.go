package main

import (
	"context"
	"fmt"
)

type Log struct {
	BlockNumber int64
	TxHash      string
	TxIndex     int
	LogIndex    int
	Address     string
	Topics      []string
	Data        string
}

func GetLogsRPC(fromBlock, toBlock int64) ([]Log, error) {
	fmt.Printf("Simulating GetLogsRPC from %d to %d\n", fromBlock, toBlock)
	var logs []Log
	for b := fromBlock; b <= toBlock; b++ {
		logs = append(logs, Log{
			BlockNumber: b,
			TxHash:      fmt.Sprintf("0xTxHash%d", b),
			TxIndex:     0,
			LogIndex:    0,
			Address:     "0xAddress",
			Topics:      []string{"0xTopic1"},
			Data:        "0xData",
		})
	}
	return logs, nil
}

func InsertLogs(ctx context.Context, logs []Log) error {
	fmt.Printf("Inserting %d logs into DB\n", len(logs))
	return nil
}

func IndexLogs(ctx context.Context, fromBlock, toBlock int64) error {
	fmt.Printf("Indexing logs from %d to %d\n", fromBlock, toBlock)
	for start := fromBlock; start <= toBlock; start += LOG_RANGE_SIZE {
		end := start + LOG_RANGE_SIZE - 1
		if end > toBlock {
			end = toBlock
		}
		logs, _ := GetLogsRPC(start, end)
		InsertLogs(ctx, logs)
	}
	return nil
}
