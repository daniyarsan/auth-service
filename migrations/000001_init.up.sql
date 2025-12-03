CREATE TABLE users
(
    id            BIGSERIAL PRIMARY KEY,
    email         TEXT      NOT NULL UNIQUE,
    password_hash TEXT      NOT NULL,
    created_at    TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE refresh_tokens
(
    user_id    BIGINT NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    token      TEXT   NOT NULL UNIQUE,
    expires_at BIGINT NOT NULL,
    PRIMARY KEY (user_id, token)
);