CREATE TABLE stats (
    stat_date DATE PRIMARY KEY DEFAULT CURRENT_DATE,
    total_calories INT NOT NULL DEFAULT 0,
    total_meals INT NOT NULL DEFAULT 0,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);