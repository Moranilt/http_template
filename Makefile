PRODUCTION=false
PORT=8080

DB_NAME=test
DB_HOST=localhost
DB_USER=root
DB_PASSWORD=123456
DB_SSL_MODE=disable

RABBITMQ_HOST=localhost:5672
RABBITMQ_USERNAME=test
RABBITMQ_PASSWORD=1234

REDIS_HOST=localhost:6379
REDIS_PASSWORD=""

TRACER_URL=http://localhost:14268/api/traces
TRACER_NAME=test

DEFAULT_ENV=PRODUCTION=$(PRODUCTION) PORT=$(PORT)

DB_ENV=DB_NAME=$(DB_NAME) DB_HOST=$(DB_HOST) DB_USER=$(DB_USER) DB_PASSWORD=$(DB_PASSWORD) DB_SSL_MODE=$(DB_SSL_MODE)

RABBITMQ_ENV=RABBITMQ_HOST=$(RABBITMQ_HOST) RABBITMQ_USERNAME=$(RABBITMQ_USERNAME) RABBITMQ_PASSWORD=$(RABBITMQ_PASSWORD)

REDIS_ENV=REDIS_HOST=$(REDIS_HOST) REDIS_PASSWORD=$(REDIS_PASSWORD)

TRACER_ENV = TRACER_URL=$(TRACER_URL) TRACER_NAME=$(TRACER_NAME)

ENVIRONMENT = $(DEFAULT_ENV) $(DB_ENV) $(TRACER_ENV) $(RABBITMQ_ENV) $(REDIS_ENV)

run:
	$(ENVIRONMENT) go run .

run-race:
	$(ENVIRONMENT) go run -race .

test:
	go test ./...

test-cover:
	go test ./... -coverprofile=./cover

cover-html:
	go tool cover -html=./cover

docker-up:
	docker compose up -d

docker-down:
	docker compose down


migrate_version = latest
host = localhost
user = root
pass = 123456
dbname = test
sslmode = disable


migrate: ./cmd
	@if [ "$(filter up down,$(MAKECMDGOALS))" = "" ]; then \
		go run ./cmd run -dbname=$(dbname) -pass=$(pass) -user=$(user) -host=$(host) -version=$(migrate_version) -sslmode=$(sslmode); \
	else \
		go run ./cmd $(filter up down,$(MAKECMDGOALS)) -dbname=$(dbname) -pass=$(pass) -user=$(user) -host=$(host) -version=$(migrate_version) -sslmode=$(sslmode); \
	fi
		
up down:
	@: