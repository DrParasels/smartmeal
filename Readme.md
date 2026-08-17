Поднять контейнер.

Применение миграций в проекте для постгри:
1) go install -tags "postgres" github.com/golang-migrate/migrate/v4/cmd/migrate@latest
2) migrate -path db/migrations -database "postgres://smartuser:smartpass@localhost:5433/smartmeal?sslmode=disable" up

Войти в постгрю в нужную таблицу:
1) docker exec -it smartmeal-postgres psql -U smartuser -d smartmeal

Генерация ogen http сервер + методы к нему:
1) Установка go install github.com/ogen-go/ogen/cmd/ogen@latest
2) Сама генерация ogen -target api/ogen -package ogen --clean api/openapi/openapi.yaml
3) Подтянуть все зависимости go mod tidy

Хендлер и сервер:
1) Описать методы в хендлере и завести структуру
2) Создать сервер и добавить в переменные хендлер
3) Так же прописать условие остановки сервереа listenAndServe.

Подключить postgresql и настроить в хендлере реальные запросы:
1) Написать запросы в файле queries sql
2) Потом настроить файл sqlc.yaml 
3) Сгенерировать код: sqlc generate

Делаю второй сервис stats с nats и protobuf:
1) создал миграцию
2) создал запрос

3) сгенерил sqlc

4) создать и сгенерить protobuf 
5) команда для генерации protoc --go_out=. --go_opt=module=github.com/yourname/smartmeal proto/events/meal_created.proto

6) добавил nats в docker compose
7) Написать минимальный код для подключения и подписку на события в nats

Нужно понять логику куда класть сообщение (NATS), какое сообщение (я его описал - MealCreated) и когда его отправлять (после create в бд)
8) Пишу запросы для второго сервиса в бд

Подписка на событие, передача и получение данных между сервисами:
1) Использование marshal unmarshal


Не забыть доделать(done):
1) явные ошибки 404/400 (done)
2) конфиг из env (done, но пока go не подхватывает env, поэтому в файле config задал переменные явно)
3) slog (done в cmd и без хендлеров)
4) Graceful shutdown (done, но позже ещё раз пройтись по нему)