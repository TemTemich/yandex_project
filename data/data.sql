CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS expressions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    expression TEXT,
    created_at TIMESTAMP,
    done_at TIMESTAMP,
    status TEXT,
    result TEXT
);

CREATE TABLE IF NOT EXISTS operations (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    expression_id UUID REFERENCES expressions(id),
    operation TEXT,
    status TEXT,
    created_at TIMESTAMP,
    done_at TIMESTAMP,
    result TEXT
);
