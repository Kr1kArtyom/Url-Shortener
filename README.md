Запуск Url-shortener с использованием локального хранилища:
1) Запускаем приложение - go run cmd/main.go

Запуск Url-shortener с использованием бд (PosgreSQL):
1) Поднимаем базу данных - docker run --name=url-shortener-db -e POSTGRES_PASSWORD='postgres' -p 5436:5432 -d --rm postgres
2) Применяем мигрцацию - migrate -path ./migrations -database 'postgres://postgres:postgres@localhost:5436/postgres?sslmode=disable' up
3) Запускаем приложение - go run cmd/main.go -d
