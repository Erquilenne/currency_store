run:
	go run cmd/v1/main.go

docker_build:
	sudo docker compose up -d --build

docker_up:
	sudo docker compose up -d

docker_down:
	sudo docker compose down

psql:
	sudo docker compose exec db psql -U postgres -d currency

createdb:
	sudo docker compose exec db createdb -U postgres currency

dropdb:
	sudo docker compose exec db dropdb -U postgres currency

swagger:
	swag init -g cmd/v1/main.go --parseDependency --parseInternal