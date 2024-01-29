.PHONY: users db-up db-down

users:
	go run main.go ./env/.env.users

project:
	go run main.go ./env/.env.project

fav:
	go run main.go ./env/.env.fav

db-up:
	docker compose -f docker-compose.db.yaml up -d

db-stop:
	docker compose -f docker-compose.db.yaml stop

grpc-path:
	export PATH="$PATH:$(go env GOPATH)/bin"

grpc-fav: 
	protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		modules/favorite/favPb/favPb.proto

grpc-com: 
	protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		modules/comment/commentPb/commentPb.proto

grpc-users: 
	protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		modules/users/usersPb/usersPb.proto
# db-down:
# 	docker compose -f docker-compose.db.yaml down

