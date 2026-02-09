CREATE SCHEMA IF NOT EXISTS subscriptions;

CREATE TABLE subscriptions.subscription (
    id SERIAL PRIMARY KEY,
    service_name VARCHAR(255) NOT NULL,
    price INTEGER NOT NULL CHECK (price >= 0),
    user_id UUID NOT NULL,
    start_date VARCHAR(10) NOT NULL,
    end_date VARCHAR(10),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    deleted_at TIMESTAMP NULL
);

CREATE INDEX idx_subscription_user_id ON subscriptions.subscription(user_id);
CREATE INDEX idx_subscription_service_name ON subscriptions.subscription(service_name);
