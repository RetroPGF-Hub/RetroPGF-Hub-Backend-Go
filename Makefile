.PHONY: users db-up db-down

users:
	go run main.go ./env/.env.users

project:
	go run main.go ./env/.env.project

datacenter:
	go run main.go ./env/.env.datacenter

db-up:
	docker compose -f docker-compose.db.yaml up -d

db-stop:
	docker compose -f docker-compose.db.yaml stop


# migrate project is includes 3 database migration
migrate-project:
	go run ./pkg/database/script/migration.go ./env/.env.project
	
migrate-user:
	go run ./pkg/database/script/migration.go ./env/.env.users
	
	

grpc-path:
	export PATH="$PATH:$(go env GOPATH)/bin"

grpc-fav: 
	protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		modules/favorite/favoritePb/favoritePb.proto
		
grpc-datacenter: 
	protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		modules/datacenter/datacenterPb/datacenterPb.proto

grpc-users: 
	protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		modules/users/usersPb/usersPb.proto
# db-down:
# 	docker compose -f docker-compose.db.yaml down

