### ==== CONFIG ==== ###
APP_NAME=auth-service
DB_USER=app_user
DB_PASSWORD=121212
DB_NAME=app
DB_HOST=localhost
DB_PORT=5432
DB_URL=postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable

MIGRATIONS_DIR=./migrations


### ==== GO COMMANDS ==== ###

run:
	go run ./cmd/server

build:
	go build -o bin/${APP_NAME} ./cmd/server

tidy:
	go mod tidy


### ==== MIGRATIONS (LOCAL) ==== ###

migrate-up:
	migrate -path ${MIGRATIONS_DIR} -database "${DB_URL}" up

migrate-down:
	migrate -path ${MIGRATIONS_DIR} -database "${DB_URL}" down 1

migrate-down-all:
	migrate -path ${MIGRATIONS_DIR} -database "${DB_URL}" down

migrate-force:
	migrate -path ${MIGRATIONS_DIR} -database "${DB_URL}" force ${version}

migrate-version:
	migrate -path ${MIGRATIONS_DIR} -database "${DB_URL}" version

migrate-create:
	migrate create -ext sql -dir ${MIGRATIONS_DIR} -seq ${name}



### ==== DOCKER ==== ###

docker-build:
	docker build -t ${APP_NAME}:latest .

docker-up:
	docker compose up -d

docker-down:
	docker compose down

docker-logs:
	docker logs -f ${APP_NAME}



### ==== QUALITY ==== ###

lint:
	golangci-lint run


### ==== CLEAN ==== ###
clean:
	rm -rf bin
