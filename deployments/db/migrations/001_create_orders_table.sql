-- +goose Up
CREATE TABLE orders (
    id BIGSERIAL PRIMARY KEY,
    number BIGINT NOT NULL UNIQUE,
    user_id BIGSERIAL NOT NULL,
    status VARCHAR(20) NOT NULL CHECK (status IN ('NEW', 'PROCESSING', 'INVALID', 'PROCESSED')),
    FOREIGN KEY (user_id) REFERENCES users(id)
);

-- +goose Down
DROP TABLE orders;
