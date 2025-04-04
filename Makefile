DB_URL=postgres://postgres:1111@localhost:5432/authdb?sslmode=disable
MIGRATIONS_DIR= ./migrations

install-migrate:
	go install github.com/golang-migrate/migrate/v4/cmd/migrate@latest

migrate-up:
	migrate -path $(MIGRATIONS_DIR) -database "$(DB_URL)" -verbose up

migrate-down:
	migrate -path $(MIGRATIONS_DIR) -database "$(DB_URL)" -verbose down 1

migrate-reset:
	migrate -path $(MIGRATIONS_DIR) -database "$(DB_URL)" -verbose down

migrate-drop:
	migrate -path $(MIGRATIONS_DIR) -database "$(DB_URL)" -verbose drop -f

migrate-force:
	migrate -path $(MIGRATIONS_DIR) -database "$(DB_URL)" force $(version)