package main

import (
    "context"
    "fmt"
    "log"
)

func FullReplayVerification(ctx context.Context) {
    fmt.Println("=== Day 17: Full Replay Verification ===")

    // Step 1: Delete DB
    fmt.Println("Deleting database...")
    if err := ResetDB(); err != nil {
        log.Fatal(err)
    }

    // Step 2: Run first replay
    fmt.Println("Running first replay...")
    if err := RunReplay(ctx); err != nil {
        log.Fatal(err)
    }

    hash1, err := ComputeDBHash(ctx)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("First replay hash:", hash1)

    // Step 3: Delete DB again
    fmt.Println("Resetting database for second replay...")
    if err := ResetDB(); err != nil {
        log.Fatal(err)
    }

    // Step 4: Run second replay
    fmt.Println("Running second replay...")
    if err := RunReplay(ctx); err != nil {
        log.Fatal(err)
    }

    hash2, err := ComputeDBHash(ctx)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("Second replay hash:", hash2)

    // Step 5: Compare hashes
    if hash1 == hash2 {
        fmt.Println("✅ Full Replay Verification PASSED — Hashes are identical!")
    } else {
        fmt.Println("❌ Full Replay Verification FAILED — Hashes differ!")
    }
}
