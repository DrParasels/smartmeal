-- name: CreateDailyMeal :one
INSERT INTO meals (name, calories) 
VALUES ($1, $2) 
RETURNING *;

-- name: ListDailyMeals :many
SELECT * FROM meals
ORDER BY created_at DESC;

-- name: GetDailyMeal :one
SELECT * FROM meals
WHERE id = $1;