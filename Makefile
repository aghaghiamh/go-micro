FRONTEND_BINARY=frontendApp
BROKER_BINARY=brokerApp
AUTH_BINARY=authApp
LOGGER_BINARY=loggerApp
MAILER_BINARY=mailerApp

## up: start all containers in the backgroung without forcing build
up:
	@echo "Starting Docker images..."
	docker-compose up -d
	@echo "Docker images started!"

## down: stop docker compose
down:
	@echo "Stopping docker compose..."
	docker-compose down
	@echo "Done!"

## up-build: stops docker-compose (if running), builds all projects images and starts docker-compose
up-build: build_brocker build_auth build_logger build_mailer
	@echo "Stopping docker containers (if running...)"
	docker-compose down
	@echo "Building (when required) and starting docker images..."
	docker-compose up --build -d
	@echo "Docker images built and started!"

## build_brocker: builds the broker binary as a linux executable
build_brocker:
	@echo "Building broker binary..."
	cd ./brocker && env GOOS=linux CGO_ENABLED=0 go build -o ${BROKER_BINARY} ./cmd/api
	@echo "Built"

## build_auth: builds the auth binary as a linux executable
build_auth:
	@echo "Building auth binary..."
	cd ./authentication && env GOOS=linux CGO_ENABLED=0 go build -o ${AUTH_BINARY} ./cmd/api
	@echo "Built"

## build_logger: builds the logger binary as a linux executable
build_logger:
	@echo "Building logger binary..."
	cd ./logger && env GOOS=linux CGO_ENABLED=0 go build -o ${LOGGER_BINARY} ./cmd/api
	@echo "Built"

## build_mailer: builds the mailer binary as a linux executable
build_mailer:
	@echo "Building mailer binary..."
	cd ./mailer && env GOOS=linux CGO_ENABLED=0 go build -o ${MAILER_BINARY} ./cmd/api
	@echo "Built"

## build_frontend: builds the frontend binary as a linux executable
build_frontend:
	@echo "Building broker binary..."
	cd ./front-end && env CGO_ENABLED=0 go build -o ${FRONTEND_BINARY} ./cmd/web
	@echo "Built"

## start: starts the frontend
start: build_frontend
	@echo "Starting the frontend"
	cd ./front-end && ./${FRONTEND_BINARY} &

## stop: stop the front end
stop:
	@echo "Stopping front end..."
	@-pkill -SIGTERM -f "./${FRONT_END_BINARY}"
	@echo "Stopped front end!"