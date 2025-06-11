DO $$
BEGIN
    CREATE TABLE IF NOT EXISTS interface_log (
        trace_id VARCHAR(50) NOT NULL,
        service_name VARCHAR(50) NOT NULL,
        client_name VARCHAR(50) NOT NULL,
        request_payload TEXT NOT NULL,
        response_payload TEXT NOT NULL,
        request_date TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
        response_date TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
        CONSTRAINT pk_action_log PRIMARY KEY (trace_id)
    );

    CREATE INDEX IF NOT EXISTS idx_trace_id ON interface_log (trace_id);
END $$;