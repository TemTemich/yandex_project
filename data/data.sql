CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    login TEXT UNIQUE,
    password TEXT,
    created_at TIMESTAMP
);

CREATE TABLE IF NOT EXISTS expressions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES users(id),
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


