# Blockchain Indexer & State Replay Engine
Project Overview

This project is a high-performance blockchain indexer and state replay engine that:

-Processes raw blockchain logs and events.

-Detects and handles chain reorganizations (reorgs).

-Reconstructs blockchain state deterministically via replay logic.

-Validates invariants to ensure correctness and integrity of the state.

-Use Case: This engine is suitable for auditing, analytics, or backend applications that require precise on-chain state reconstruction.

# Tech Stack

-Language: Go (Golang)

-Data: JSON-formatted blockchain logs

-Tools: CLI-based Go commands for indexing, replay, and invariant validation

-Version Control: Git

 # Key Features

- Reorg Detection & Handling: Automatically detects chain reorganizations and rolls back/apply state changes.

- Deterministic Replay: Rebuilds the exact blockchain state from raw logs.

- Invariant Validation: Ensures rules like no negative balances, total liquidity consistency, and swap integrity are always enforced.

- Full State Audit: Provides transparent verification of blockchain events for correctness.


 
1. Run Indexer
go run main.go logs.go reorg.go

2. Replay State
go run ./cmd/replay_state

3. Validate Invariants
go run ./cmd/invariants

4. Sample Output
Block 17000001 -> State Hash: 2378a57a2e5e24b25f29b94c70d4ec88766da1f6
STATE pair1 150 250



