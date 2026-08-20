Применение миграций PostgreSQL:
1) go install -tags "postgres" github.com/golang-migrate/migrate/v4/cmd/migrate@latest
2) migrate -path db/migrations -database "postgres://smartuser:smartpass@localhost:5433/smartmeal?sslmode=disable" up

Генерация sqlc, ogen, protobuf:
go generate ./...



