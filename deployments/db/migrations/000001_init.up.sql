-- +goose Up
CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    login VARCHAR(50) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL
);

CREATE TABLE orders (
    id BIGSERIAL PRIMARY KEY,
    number BIGINT NOT NULL UNIQUE,
    user_id BIGSERIAL NOT NULL,
    sum DECIMAL(10, 2) DEFAULT 0,
    accrual DECIMAL(10, 2) DEFAULT 0,
    uploaded_at TIMESTAMP DEFAULT current_timestamp,
    status VARCHAR(20) NOT NULL CHECK (status IN ('NEW', 'PROCESSING', 'INVALID', 'PROCESSED')),
    FOREIGN KEY (user_id) REFERENCES users(id)
);

CREATE TABLE balances (
    id BIGSERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    current DECIMAL(10, 2) DEFAULT 0,
    withdrawn DECIMAL(10, 2) DEFAULT 0,
    FOREIGN KEY (user_id) REFERENCES users(id)
);

-- +goose Down
DROP TABLE users;
DROP TABLE orders;
DROP TABLE balances;