.PHONY: up down restart logs build clean init-kafka ps

up:
	docker-compose up -d

down:
	docker-compose down

restart:
	docker-compose down && docker-compose up -d

logs:
	docker-compose logs -f --tail=100

ps:
	docker-compose ps

build:
	docker-compose build --no-cache

clean:
	docker-compose down -v --remove-orphans

init-kafka:
	docker exec kafka bash -c "\
	kafka-topics --bootstrap-server localhost:9092 --list | grep -q '^logs.ingest$$' || \
	kafka-topics --create --topic logs.ingest \
	--bootstrap-server localhost:9092 --replication-factor 1 --partitions 1"

db-migrate:
	go run ./cmd db --migrate

db-rollback:
	go run ./cmd db --rollback

db-version:
	go run ./cmd db --version

db-drop:
	go run ./cmd db --drop

run-generator:
	go run ./cmd/log-generator/main.go

run-ingestor:
	go run ./cmd/log-ingestor/main.go

run-api:
	go run ./cmd/api-server/main.go

test-pipeline: up init-kafka db-migrate
	@echo "🚀 Starting full pipeline test..."
	@echo "📊 Run 'make run-ingestor' in one terminal"
	@echo "📝 Run 'make run-generator' in another terminal"
	@echo "🔍 Check PostgreSQL with: docker exec -it postgres psql -U loguser -d logdb -c 'SELECT * FROM log_raw;'"