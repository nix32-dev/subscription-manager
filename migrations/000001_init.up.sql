CREATE TABLE subscriptions (
    id SERIAL PRIMARY KEY,
    user_id UUID NOT NULL,
    service_name VARCHAR(20) NOT NULL,
    price INTEGER NOT NULL CHECK (price>0),
    start_date VARCHAR(7) NOT NULL,
    end_date VARCHAR(7)
)