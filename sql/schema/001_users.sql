-- +goose up
CREATE TABLE users (
    id UUID primary key,
    created_at TIMESTAMP not null,
    updated_at TIMESTAMP not null,
    email TEXT unique not null
);

-- +goose down
DROP TABLE users;