createdb:
	docker exec -it postgres createdb --username=postgres --owner=postgres simple_bank

dropdb:
	docker exec -it postgres dropdb --username=postgres simple_bank

migrateup:
	sudo docker run -v ~/Documents/Projects/simple-bank/db/migrations:/migrations --network host migrate/migrate -path=/migrations/ -database "postgres://postgres:password@localhost:5432/simple_bank?sslmode=disable" -verbose up

migratedown:
	sudo docker run -v ~/Documents/Projects/simple-bank/db/migrations:/migrations --network host migrate/migrate -path=/migrations/ -database "postgres://postgres:password@localhost:5432/simple_bank?sslmode=disable" -verbose down 1

migrateforce:
	sudo docker run -v ~/Documents/Projects/simple-bank/db/migrations:/migrations --network host migrate/migrate -path=/migrations/ -database "postgres://postgres:password@localhost:5432/simple_bank?sslmode=disable" -verbose force 1


sqlc:
	sqlc generate

test:
	go clean -testcache
	go test -v -cover ./...

.PHONY: postgres createdb dropdb migrateup migratedown sqlc