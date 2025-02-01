CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE outbox (
    id UUID DEFAULT uuid_generate_v4(),  -- UUID type with default value using uuid_generate_v4
    event_type VARCHAR(255) NOT NULL,                -- Type of message (e.g., 'send_notification')
    status VARCHAR(50) NOT NULL,                     -- Status of the message (e.g., PENDING, SENT, FAILED)
    attempt INT DEFAULT 0,                         -- Retry attempts (initially 0)
    destination_type VARCHAR(50) NOT NULL,           -- Type of destination (e.g., asyncq, rabbitmq, kafka)
    sent_at TIMESTAMP WITH TIME ZONE,                -- Timestamp when the message was sent (nullable)
    execute_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP, -- Timestamp for the next run
    payload JSONB NOT NULL DEFAULT '{}'::jsonb,      -- JSONB type with default empty object
    error_message TEXT,                               -- Error message (nullable)
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP, -- Timestamp when created
    PRIMARY KEY (id, execute_at)
) PARTITION BY RANGE (execute_at);

CREATE INDEX idx_outbox_filter ON outbox (execute_at asc, status);

-- Create partitions for the "outbox" table
DO $$ 
DECLARE 
    start_date DATE := '2024-01-01';
    end_date DATE := '2026-12-31';
    partition_start DATE;
    partition_end DATE;
    partition_name TEXT;
    sql_statement TEXT;
BEGIN 
    partition_start := start_date;

    WHILE partition_start <= end_date LOOP
        partition_end := partition_start + INTERVAL '1 month';
        partition_name := 'outbox_' || TO_CHAR(partition_start, 'YYYY_MM');

        -- Construct the SQL statement
        sql_statement := 'CREATE TABLE IF NOT EXISTS ' || partition_name || ' PARTITION OF outbox ' ||
                         'FOR VALUES FROM (' || quote_literal(partition_start) || ') TO (' || quote_literal(partition_end) || ');';

        -- Execute the SQL statement
        EXECUTE sql_statement;
        
        -- Move to the next month
        partition_start := partition_end;
    END LOOP;
END $$;
