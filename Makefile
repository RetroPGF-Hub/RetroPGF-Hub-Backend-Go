.PHONY: users db-up db-down

users:
	go run main.go ./env/.env.users

project:
	go run main.go ./env/.env.project

db-up:
	docker compose -f docker-compose.db.yaml up -d

db-stop:
	docker compose -f docker-compose.db.yaml stop

# db-down:
# 	docker compose -f docker-compose.db.yaml down

