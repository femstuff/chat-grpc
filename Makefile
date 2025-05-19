.PHONY: build up down cli logs clean migrate-up migrate-down run-all auth chat saga notification

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

cli:
	docker-compose run --service-ports --rm chat-cli
  
migrate-up:
	docker run --rm --network host -v $(PWD)/migrations:/migrations migrate/migrate -path=/migrations -database "postgres://auth_user:auth_pass@localhost:5432/auth_db?sslmode=disable" up
	docker run --rm --network host -v $(PWD)/migrations:/migrations migrate/migrate -path=/migrations -database "postgres://user:user_pass@localhost:5433/users_db?sslmode=disable" up

migrate-down:
	docker run --rm --network host -v $(PWD)/migrations:/migrations migrate/migrate -path=/migrations -database "postgres://auth_user:auth_pass@localhost:5432/auth_db?sslmode=disable" down

run-all: build up migrate-up

auth:
	docker-compose up -d auth-service

chat:
	docker-compose up -d chat-service

notification:
	docker-compose up -d notification-service

saga:
	docker-compose up -d saga-orchestrator
