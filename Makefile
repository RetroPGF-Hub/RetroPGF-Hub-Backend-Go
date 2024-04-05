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
	
migrate-users:
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


build-docker:
	docker build -f ./build/Dockerfile -t retro-pgf-hub:latest .

tag:
	docker image tag retro-pgf-hub:latest tgrziminiar/retro-pgf-hub:latest

push:
	docker push tgrziminiar/retro-pgf-hub:latest


kube-svc-project:
	kubectl apply -f ./build/project/project-service.yml
	kubectl apply -f ./build/project/project-deployment.yml

	
kube-svc-datacenter:
	kubectl apply -f ./build/datacenter/datacenter-service.yml
	kubectl apply -f ./build/datacenter/datacenter-deployment.yml


kube-svc-users:
	kubectl apply -f ./build/users/users-service.yml
	kubectl apply -f ./build/users/users-deployment.yml


kube-create-configmap:
	kubectl create configmap app-env --from-file=./env/prod/.env
	kubectl get configmaps


kube-install-ingress:
	kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/controller-v1.8.2/deploy/static/provider/cloud/deploy.yaml
	kubectl scale deployment users-deployment --replicas=1