-- +goose Up
-- +goose StatementBegin
SELECT
    'up SQL query';

CREATE TABLE
    IF NOT EXISTS subscription (
        id BIGSERIAL PRIMARY KEY,
        service_name VARCHAR(255) NOT NULL,
        user_id UUID NOT NULL,
        price SERIAL NOT NULL,
        start_date DATE NOT NULL,
        end_date DATE NULL,
        created_at TIMESTAMP NOT NULL DEFAULT NOW (),
        updated_at TIMESTAMP NOT NULL DEFAULT NOW ()
    );

CREATE INDEX idx_subscriptions_user_service_date ON subscription (user_id, service_name, start_date) INCLUDE (price);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
SELECT
    'down SQL query';

DROP TABLE subscription;

-- +goose StatementEnd