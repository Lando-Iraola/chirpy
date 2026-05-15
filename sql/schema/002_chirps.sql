-- +goose up
CREATE TABLE chirps(
    id UUID primary key,
    created_at TIMESTAMP not null,
    updated_at TIMESTAMP not null,
    body text not null,
    user_id UUID NOT NULL references users(id) ON DELETE CASCADE
);

-- +goose down
DROP TABLE chirps;
