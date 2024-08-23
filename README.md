# HTTP Template

HTTP microservice template for fast start

The idea of this package is to make standart template which you can clone or use [gonew](https://go.dev/blog/gonew) tool.

From the box it contains clients for:
- [Vault](https://www.vaultproject.io/)
- [Postgresql](https://www.postgresql.org/)
- [Redis](https://redis.io/)
- [Rabbitmq](https://www.rabbitmq.com/)

Tracing:
- [Jaeger](https://www.jaegertracing.io/)
- [Opentelemetry](https://opentelemetry.io/)

## Usage
1. Install [gonew](https://go.dev/blog/gonew)
2. Run:
```bash
go install golang.org/x/tools/cmd/gonew@latest
```
3. Run:
```bash
gonew github.com/Moranilt/http_template example.com/project_name
```
4. Happy coding!

## Makefile
To fast usage you can use makefile commands
- `run` - runs you app with local environment variables
- `run-race` - runs you app with local environment variables with `-race` flag
- `test` - runs tests in all folders
- `test-cover` - runs tests with coverage
- `cover-html` - opens your `cover` file in browser
- `docker-up` - runs docker compose file
- `docker-down` - runs docker compose down command
- `migrate` - runs migration to the latest version or specific version using `-version` flag
- `migrate up` - runs migration to the next step from current
- `migrate down` - runs migration to previous step from current

`migrate` commands accept:
- version - version that you need to migrate to. You can pass `latest` to migrate to the latest version.
- dbname - database name
- host - database host
- user - username for database
- pass - users password for database
- sslmode - sslmode for database, not required. By default it disabled. You can watch all defaults in Makefile.

Feel free to modify environment variables, but beware to not break default configuration rules.

## Folders
### Clients
This folder contains all clients for external services. Implement `healthcheck.Checker` interface if you want to use your service in `/health` endpoint.

`credentials` - contains all credential structures for every service. Feel free to modify and add your own credentials.  
`database`- database client which implements `healthcheck.Checker` interface.  
`rabbitmq` - RabbitMQ client with default logic to push and consume messages. Also implements `healthcheck.Checker` interface.  
`redis` - redis client which implements `healthcheck.Checker` interface.  
`vault` - default Vault client.

### Config
Contains logic to read ENV-variables, validate and store it to default application config. Feel free to modify.

1. Add new constant named `ENV_{your_name}`
2. Add this constant to array
3. Modify default `Config` structure
4. Using your new constant as key, read from `result` map and store you variable to `Config` structure

### Endpoints
Store all endpoints into `MakeEndpoints` function.

Modify `MakeHealthEndpoint` to add new client for healthcheck.

### Healthcheck
Logic to make Healthcheck handle function for route.

### Logger
Contains logger using [logrus](https://github.com/sirupsen/logrus). Added function `WithRequestInfo` to add **requestId** from context to logs. Feel free to modify.

### Middleware
Contains all middlewares for your application. It has default middleware to add `X-Request-ID` header and log every incoming request. Feel free to modify.

### Migrations
Contains all `sql` files to run migrations using [golang-migrate](https://github.com/golang-migrate/migrate).

### Models
Store all structures for request and response in `repository` folder.

### Repository
Core logic of your application. The main rule to implement `func(context.Context, *Request) (*Response, error)` interface. There are some examples in this folder.

### Service
HTTP wrapper for repository. It contains unique logic with [handler](https://pkg.go.dev/github.com/Moranilt/http_template/utils/handler) pakcage using generics.

### Tracer
Default tracer implementation. Feel free to modify.

### Transport
Default settings to create http-transport using [gorilla mux](https://github.com/gorilla/mux). Feel free to modify or add more transports.

### Utils
Helper functions to make your life easier.

## TODO
- stream response
- add README template into `docs` folder
