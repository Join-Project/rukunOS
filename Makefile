.PHONY: help setup up down logs ps restart build clean migrate backup

# Default target
help:
	@echo "Available commands:"
	@echo "  make setup      - Initial project setup"
	@echo "  make up         - Start all services"
	@echo "  make down       - Stop all services"
	@echo "  make logs       - View logs"
	@echo "  make ps         - Check service status"
	@echo "  make restart    - Restart all services"
	@echo "  make build      - Build all containers"
	@echo "  make clean      - Stop and remove containers, volumes"
	@echo "  make migrate    - Run database migrations"
	@echo "  make backup     - Backup database"
	@echo "  make shell-api  - Open shell in API container"
	@echo "  make shell-db   - Open psql in database"
	@echo "  make prod-up    - Start production services"
	@echo "  make prod-down  - Stop production services"

# Initial setup
setup:
	@chmod +x scripts/*.sh
	@./scripts/setup.sh

# Development
up:
	docker compose up -d

down:
	docker compose down

logs:
	docker compose logs -f

ps:
	docker compose ps

restart:
	docker compose restart

build:
	docker compose build

clean:
	docker compose down -v
	docker system prune -f

# Database
migrate:
	cd backend && ./scripts/run-migrations.sh docker

backup:
	./scripts/backup-db.sh

# Shell access
shell-api:
	docker compose exec api sh

shell-db:
	docker compose exec db psql -U {{PROJECT_NAME}}_user -d {{PROJECT_NAME}}_db

# Production
prod-up:
	docker compose -f docker-compose.prod.yml up -d --build

prod-down:
	docker compose -f docker-compose.prod.yml down

prod-logs:
	docker compose -f docker-compose.prod.yml logs -f

prod-restart:
	docker compose -f docker-compose.prod.yml restart

