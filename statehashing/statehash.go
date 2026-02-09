package main

import (
    "context"
    "crypto/sha256"
    "database/sql"
    "encoding/json"
    "fmt"
    "log"
    "sort"

    _ "github.com/lib/pq"
)

// ---------------- CONFIG ----------------
const (
    RPC_URL      = "https://eth-mainnet.g.alchemy.com/v2/YOUR_KEY"
    PG_CONN      = "postgres://stateuser:mypassword@localhost:5432/blockchain?sslmode=disable"


    START_BLOCK  = 17000000
    END_BLOCK    = 17000010
)

// ---------------- STATE STRUCT ----------------
type AccountState struct {
    Address string            `json:"address"`
    Balance string            `json:"balance"`
    Nonce   uint64            `json:"nonce"`
    Storage map[string]string `json:"storage"`
}

// ---------------- MAIN ----------------
func main() {
    ctx := context.Background()

    // Connect to Postgres
    db, err := sql.Open("postgres", PG_CONN)
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    for block := START_BLOCK; block <= END_BLOCK; block++ {
        // Step 1: Fetch state (simplified for demonstration)
        state := FetchBlockState(block)

        // Step 2: Compute hash
        hash := ComputeStateHash(state)
        fmt.Printf("Block %d -> State Hash: %s\n", block, hash)

        // Step 3: Store in DB
        err := StoreStateHash(ctx, db, block, hash)
        if err != nil {
            log.Printf("Error storing state hash for block %d: %v\n", block, err)
        }
    }

    fmt.Println("✅ State hashing completed.")
}

// ---------------- FETCH STATE ----------------
// Mock function: replace with actual RPC calls to get balances/contracts
func FetchBlockState(block int) []AccountState {
    // Example: deterministic dummy data
    accounts := []AccountState{
        {Address: "0xAAA", Balance: "100", Nonce: 1, Storage: map[string]string{"key1": "val1"}},
        {Address: "0xBBB", Balance: "200", Nonce: 2, Storage: map[string]string{"key2": "val2"}},
    }

    // Sort by address for deterministic hash
    sort.Slice(accounts, func(i, j int) bool {
        return accounts[i].Address < accounts[j].Address
    })

    return accounts
}

// ---------------- COMPUTE HASH ----------------
func ComputeStateHash(state []AccountState) string {
    data, err := json.Marshal(state)
    if err != nil {
        log.Fatal(err)
    }
    hash := sha256.Sum256(data)
    return fmt.Sprintf("%x", hash)
}

// ---------------- STORE IN DB ----------------
func StoreStateHash(ctx context.Context, db *sql.DB, block int, hash string) error {
    _, err := db.ExecContext(ctx,
        `INSERT INTO state_hashes (block_number, state_hash) VALUES ($1, $2)
        ON CONFLICT (block_number) DO UPDATE SET state_hash = EXCLUDED.state_hash`,
        block, hash)
    return err
}
