CREATE TABLE blocks (
    number BIGINT PRIMARY KEY,
    hash TEXT,
    parent_hash TEXT
);

CREATE TABLE events (
    id SERIAL PRIMARY KEY,
    block_number BIGINT,
    tx_hash TEXT,
    event_type TEXT,
    data JSONB
);

CREATE TABLE sync_events (
    pair_address TEXT,
    reserve0 NUMERIC,
    reserve1 NUMERIC,
    block_number BIGINT
);

CREATE TABLE pair_reserves (
    pair_address TEXT PRIMARY KEY,
    reserve0 NUMERIC,
    reserve1 NUMERIC,
    block_number BIGINT
);
