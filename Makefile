PRODUCTION=false
PORT=8080

VAULT_TOKEN=censor-token
VAULT_HOST=http://localhost:8200
VAULT_MOUNT_PATH=secret
VAULT_DB_CREDS_PATH=censor/db
VAULT_RABBITMQ_CREDS_PATH=censor/rabbitmq

TRACER_URL=http://localhost:14268/api/traces
TRACER_NAME=censor

DEFAULT_ENV=PRODUCTION=$(PRODUCTION) PORT=$(PORT)

VAULT_ENV = VAULT_TOKEN=$(VAULT_TOKEN) VAULT_HOST=$(VAULT_HOST) VAULT_MOUNT_PATH=$(VAULT_MOUNT_PATH) VAULT_DB_CREDS_PATH=$(VAULT_DB_CREDS_PATH) VAULT_RABBITMQ_CREDS_PATH=$(VAULT_RABBITMQ_CREDS_PATH)

TRACER_ENV = TRACER_URL=$(TRACER_URL) TRACER_NAME=$(TRACER_NAME)

ENVIRONMENT = $(DEFAULT_ENV) $(VAULT_ENV) $(TRACER_ENV)

run:
	$(ENVIRONMENT) go run .

test:
	go test ./...