CREATE TABLE meals (
    id BIGSERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    calories INT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);