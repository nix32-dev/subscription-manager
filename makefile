include .env
export

service-run:
	go run cmd/main.go

service-build:
	go build cmd/main.go

service-build-and-run:
	go build cmd/main.go
	go run main.go

migrate-up:
	migrate -path migrations -database ${DB_CONN} up

migrate-down:
	migrate -path migrations -database ${DB_CONN} down

deploy-up:
	docker compose up -d subscriptions-service
	make migrate-up

deploy-down:
	docker-compose down