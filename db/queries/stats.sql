-- name: UpsertDailyStat :one

INSERT INTO stats (stat_date, total_calories, total_meals)
VALUES ($1, $2, 1)
ON CONFLICT (stat_date) DO UPDATE SET 
    total_calories = stats.total_calories + EXCLUDED.total_calories, 
    total_meals = stats.total_meals + EXCLUDED.total_meals,
    updated_at = now()
RETURNING *;