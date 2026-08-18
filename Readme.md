Поднять контейнер.

Применение миграций в проекте для постгри:
1) go install -tags "postgres" github.com/golang-migrate/migrate/v4/cmd/migrate@latest
2) migrate -path db/migrations -database "postgres://smartuser:smartpass@localhost:5433/smartmeal?sslmode=disable" up

Войти в постгрю в нужную таблицу:
1) docker exec -it smartmeal-postgres psql -U smartuser -d smartmeal

Генерация ogen http сервер + методы к нему:
1) Установка go install github.com/ogen-go/ogen/cmd/ogen@latest
3) Подтянуть все зависимости go mod tidy

Подключить postgresql и настроить в хендлере реальные запросы:
1) Написать запросы в файле queries sql
2) Потом настроить файл sqlc.yaml

Хендлер и сервер:
1) Описать методы в хендлере и завести структуру
2) Создать сервер и добавить в переменные хендлер
3) Так же прописать условие остановки сервереа listenAndServe.

Делаю второй сервис stats с nats и protobuf:
1) создал миграцию
2) создал запрос
3) добавил nats в docker compose
4) Написать минимальный код для подключения и подписку на события в nats

Нужно понять логику куда класть сообщение (NATS), какое сообщение (я его описал - MealCreated) и когда его отправлять (после create в бд)
8) Пишу запросы для второго сервиса в бд

Подписка на событие, передача и получение данных между сервисами:
1) Использование marshal unmarshal



Правки:
1) go:generate (done)
2) Манифест sqlc рядом со схемой (done)
3) Манифест protocol buffers (done)
4) вместо logger -> zerolog (done, пока без NATS)
5) отдельно config не надо (done)
6) в api сервер отдельно вынести (done)