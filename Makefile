include .env

run:
	docker-compose up --build

producer:
	go run cmd/test_producer.go

consumer:
	go run cmd/app/main.go

.PHONY: run producer consumer