include .env
export

CONN_STRING = postgresql://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_MODELESS)

MIGRATION_DIRS = internal/db/migrations

ENV=.env
PRODUCTION_COMPOSE=docker-compose.production.yml
DEVELOP_COMPOSE=docker-compose.develop.yml
LOCALHOST_COMPOSE=docker-compose.localhost.yml

# Import database
import:
	docker exec -i mini-ecommerce-db psql -U huydo -d project-mini-ecommerce < ./backupdb-project-mini-ecommerce.sql

# Export database
export:
	docker exec -i mini-ecommerce-db pg_dump -U huydo -d project-mini-ecommerce > ./backupdb-project-mini-ecommerce.sql

# Run server
server:
	cd cmd/api && go run main.go

sqlc:
	sqlc generate

# Create a new migration (make migrate-create NAME=profiles)
migrate-create:
	migrate create -ext sql -dir $(MIGRATION_DIRS) -seq $(NAME)

# Run all pending migration (make migrate-up)
migrate-up:
	migrate -path $(MIGRATION_DIRS) -database "$(CONN_STRING)" up

# Rollback the last migration
migrate-down:
	migrate -path $(MIGRATION_DIRS) -database "$(CONN_STRING)" down 1

# Rollback N migrations
migrate-down-n:
	migrate -path $(MIGRATION_DIRS) -database "$(CONN_STRING)" down $(N)

# Force migration version (use with caution example: make migrate-force VERSION=1)
migrate-force:
	migrate -path $(MIGRATION_DIRS) -database "$(CONN_STRING)" force $(VERSION)

# Drop everything (include schema migration)
migrate-drop:
	migrate -path $(MIGRATION_DIRS) -database "$(CONN_STRING)" drop

# Apply specific migration version (make migrate-goto VERSION=1)
migrate-goto:
	migrate -path $(MIGRATION_DIRS) -database "$(CONN_STRING)" goto $(VERSION)

# Localhost
localhost:
	docker compose -f $(LOCALHOST_COMPOSE) down
	docker compose -f $(LOCALHOST_COMPOSE) --env-file $(ENV) up -d --build
stop-local:
	docker compose -f $(LOCALHOST_COMPOSE) down
logs-local:
	docker compose -f $(LOCALHOST_COMPOSE) logs -f --tail=100

# Dev
develop:
	docker compose -f $(DEVELOP_COMPOSE) down
	docker compose -f $(DEVELOP_COMPOSE) --env-file $(ENV) up --build

# Production
production:
	docker compose -f $(PRODUCTION_COMPOSE) down
	docker compose -f $(PRODUCTION_COMPOSE) --env-file $(ENV) up -d --build
stop-prod:
	docker compose -f $(PRODUCTION_COMPOSE) down
logs-prod:
	docker compose -f $(PRODUCTION_COMPOSE) logs -f --tail=100

.PHONY: importdb exportdb server migrate-create migrate-up migrate-down migrate-force migrate-drop migrate-goto migrate-down-n sqlc localhost stop-local logs-local production stop-prod logs-prod develop
