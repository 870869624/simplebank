postgres:
	docker run --name simplebank-db-1 dropdb -p 5432:5432 -e POSTGRES_PASSWORD=123456 -d simplebank

createdb:
	docker exec -it simplebank-simplebanl-1 createdb --username=postgres simple_bank

dropdb:
	docker exec -it simplebank-simplebanl-1 dropdb --username=postgres simple_bank

migrateup:
	migrate -path ./db/migration -database "postgresql://postgres:123456@localhost/simple_bank?sslmode=disable" -verbose up

migrateup1:
	migrate -path ./db/migration -database "postgresql://postgres:123456@localhost/simple_bank?sslmode=disable" -verbose up 1

migratedown:
	migrate -path ./db/migration -database "postgresql://postgres:123456@localhost/simple_bank?sslmode=disable" -verbose down

migratedown1:
	migrate -path ./db/migration -database "postgresql://postgres:123456@localhost/simple_bank?sslmode=disable" -verbose down 1

sqlc:
	sqlc generate

migrate:
	migrate create -ext sql -dir db/migrations -seq create_users_table

test:
	go test -v -cover ./...
server:
	go run main.go
mock:
	mockgen -destination db/mock/store.go simplebank/db/sqlc Store

.PHONY: postgres dropdb createdb migrateup migratedown sqlc test migrate server mock migratedown1 migrateup1
