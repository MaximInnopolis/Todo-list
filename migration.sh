goose -dir ./migrations postgres "postgres://postgres:password@localhost:5432/todolist?sslmode=disable" status

goose -dir ./migrations postgres "postgres://postgres:password@localhost:5432/todolist?sslmode=disable" up