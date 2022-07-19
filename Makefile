postgres:
	docker run --name postgres14 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:14.4-alpine

createdb:
	docker exec -it postgres14 createdb --username=root --owner=root aqualog

dropdb:
	docker exec -it postgres14 dropdb aqualog

migrateup:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/aqualog?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/aqualog?sslmode=disable" -verbose down

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

.PHONY: postgres createdb dropdb migrateup migratedown sqlc