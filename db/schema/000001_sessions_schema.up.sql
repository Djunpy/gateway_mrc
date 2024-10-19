CREATE TABLE sessions (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    access_token TEXT,
    refresh_token TEXT,
    session_data TEXT,
    user_agent VARCHAR NOT NULL,
    client_ip VARCHAR NOT NULL,
    is_blocked BOOL DEFAULT false,
    last_active TIMESTAMPTZ,
    expires_at TIMESTAMPTZ NOT NULL,
    session_length_seconds INTEGER DEFAULT 0,
    created_at TIMESTAMPTZ DEFAULT NOW()
);