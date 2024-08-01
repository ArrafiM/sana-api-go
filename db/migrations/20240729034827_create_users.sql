-- +goose Up
CREATE TABLE IF NOT EXISTS users (
    id BIGSERIAL PRIMARY KEY,
    email VARCHAR(50) UNIQUE NOT NULL,
    phone VARCHAR(50) NOT NULL,
    password VARCHAR(255) NOT NULL,
    name VARCHAR(50) NOT NULL,
    role_id BIGINT,
    FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE SET NULL,
    created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at timestamp
);

-- +goose Down
DROP TABLE users;