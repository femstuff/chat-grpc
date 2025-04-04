.PHONY: build up down cli logs clean migrate-up migrate-down run-and-cli

build:
	docker-compose build

up:
	docker-compose up -d

down:
	docker-compose down

cli:
	docker-compose run --service-ports --rm chat-cli

logs:
	docker-compose logs -f

clean:
	docker-compose down --rmi all -v
	docker system prune -f

migrate-up:
	docker run --rm --network host -v $(PWD)/migrations:/migrations migrate/migrate -path=/migrations -path=/migrations -database "postgres://auth_user:auth_pass@localhost:5432/auth_db?sslmode=disable" up

migrate-down:
	docker run --rm --network host -v $(PWD)/migrations:/migrations migrate/migrate -path=/migrations -database "postgres://auth_user:auth_pass@localhost:5432/auth_db?sslmode=disable" down

run-and-cli: up migrate-up cli
