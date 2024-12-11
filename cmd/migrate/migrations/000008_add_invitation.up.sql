CREATE TABLE IF NOT EXISTS user_invitations(
    token bytea PRIMARY KEY,
    user_id bigint NOT NULL,
    expiry timestamp(0) with time zone NOT NULL
);