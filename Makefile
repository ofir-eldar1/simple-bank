DB_CONTAINER=postgres16
DATABASE_NAME=simple_bank

postgres:
	docker run --name ${DB_CONTAINER} -e POSTGRES_PASSWORD=pass -e POSTGRES_USER=root -p 5432:5432 -d postgres:16-alpine

createdb:
	docker exec -it ${DB_CONTAINER} createdb --username=root --owner=root ${DATABASE_NAME}

dropdb:
	docker exec -it ${DB_CONTAINER} dropdb --username=root --owner=root ${DATABASE_NAME}

migrateup:
	migrate -path db/migrations -database "postgresql://root:pass@localhost:5432/${DATABASE_NAME}?sslmode=disable" -verbose up	

migratedown:
	migrate -path db/migrations -database "postgresql://root:pass@localhost:5432/${DATABASE_NAME}?sslmode=disable" -verbose down

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

.PHONY: postgres createdb dropdb migrateup migratedown sqlc test