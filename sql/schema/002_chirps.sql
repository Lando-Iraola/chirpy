-- +goose up
CREATE TABLE chirps(
    id UUID primary key,
    created_at TIMESTAMP not null,
    updated_at TIMESTAMP not null,
    body text not null,
    user_id UUID,
    foreign key (user_id) references users(id)
);

-- +goose down
DROP TABLE chirps;
