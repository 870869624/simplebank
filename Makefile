postgres:
	docker run --name simplebank-db-1 dropdb -p 5432:5432 -e POSTGRES_PASSWORD=123456 -d simplebank

createdb:
	docker exec -it simplebank-simplebanl-1 createdb --username=postgres simple_bank

dropdb:
	docker exec -it simplebank-simplebanl-1 dropdb --username=postgres simple_bank

migrateup:
	migrate -path ./db/migration -database "postgresql://postgres:123456@localhost/simple_bank?sslmode=disable" -verbose up

migratedown:
	migrate -path ./db/migration -database "postgresql://postgres:123456@localhost/simple_bank?sslmode=disable" -verbose down

sqlc:
	sqlc generate

migrate:
	migrate create -ext sql -dir db/migrations -seq create_users_table

test:
	go test -v -cover ./...

.PHONY: postgres dropdb createdb migrateup migratedown sqlc
