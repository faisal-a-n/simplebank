createdb:
	docker exec -it postgres createdb --username=postgres --owner=postgres simple_bank

dropdb:
	docker exec -it postgres dropdb --username=postgres simple_bank

migrateup:
	migrate -path db/migrations -database "postgres://postgres:password@localhost:5432/simple_bank?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migrations -database "postgres://postgres:password@localhost:5432/simple_bank?sslmode=disable" -verbose down 1

migrateforce:
	migrate -path db/migrations -database "postgres://postgres:password@localhost:5432/simple_bank?sslmode=disable" -verbose force 1


sqlc:
	sqlc generate

test:
	go clean -testcache
	go test -v -cover ./...

server:
	go run .

mock:
	mockgen -package mock_db -destination db/mock/store.go github.com/faisal-a-n/simplebank/db/sqlc Store

.PHONY: postgres createdb dropdb migrateup migratedown sqlc server mock