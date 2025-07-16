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
