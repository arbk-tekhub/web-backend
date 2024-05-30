## Colors
COLOR_RESET   = \033[0m
COLOR_INFO    = \033[32m
COLOR_COMMENT = \033[33m

## Variables
MAIN_PACKAGE_PATH := ./cmd/api
BINARY_NAME := api
SERVER_PORT := 8080

.PHONY: help
## Help
help:
	@printf "${COLOR_COMMENT}Usage:${COLOR_RESET}\n"
	@printf " make [target] [args...]\n\n"
	@printf "${COLOR_COMMENT}Available targets:${COLOR_RESET}\n"
	@awk '/^[a-zA-Z\-\0-9\.@]+:/ { \
		helpMessage = match(lastLine, /^## (.*)/); \
		if (helpMessage) { \
			helpCommand = substr($$1, 0, index($$1, ":")); \
			helpMessage = substr(lastLine, RSTART + 3, RLENGTH); \
			printf " ${COLOR_INFO}%-16s${COLOR_RESET} %s\n", helpCommand, helpMessage; \
		} \
	} \
	{ lastLine = $$0 }' $(MAKEFILE_LIST)

# ==================================================================================== #
# QUALITY CONTROL
# ==================================================================================== #

.PHONY: tidy
## format code and tidy modfile
tidy:
	go mod tidy -v
	go fmt ./...


.PHONY: audit
## run quality control checks
audit:
	go mod verify
	go vet ./...
	go run honnef.co/go/tools/cmd/staticcheck@latest -checks=all,-ST1000,-U1000 ./...
	go run golang.org/x/vuln/cmd/govulncheck@latest ./...


.PHONY: run_unit_tests
## run unit tests
test:
	go test -race -buildvcs -vet=off ./...

# ==================================================================================== #
# DEVELOPMENT
# ==================================================================================== #

.PHONY: build
## build the application
build:
	go build -ldflags='-s -w' -o=./build/${BINARY_NAME} ${MAIN_PACKAGE_PATH}


.PHONY: run
## run the application
run: build
	./build/${BINARY_NAME} -port ${SERVER_PORT}


.PHONY: run/live
## run the application with reloading on file changes
run/live:
	go run github.com/cosmtrek/air@v1.43.0 \
		--build.cmd "make build" --build.bin "./build/${BINARY_NAME}" --build.delay "100" \
		--build.exclude_dir "" \
		--build.include_ext "go,html,js" \
		--misc.clean_on_exit "true"


.PHONY: swag
## generate swagger docs
swag:
	swag init --dir ${MAIN_PACKAGE_PATH}

# ==================================================================================== #
# Docker
# ==================================================================================== #

.PHONY: compose/build
## Build all/specific resource(s) defined in the compose file
compose/build:
	docker compose build ${svc}

.PHONY: compose/up
## Create and start all/specifc resource(s) defined in the compose file
compose/up:
	docker compose up -d ${svc}

.PHONY: compose/down
## Stops and removes all resources defined in the compose file
compose/down:
	docker compose down

.PHONY: compose/logs
## View output from all/specific running service(s)
compose/logs:
	docker compose logs -f ${svc}