## help: print this help message
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'
.PHONY: confirm 
confirm:
	@echo -n 'Are you sure? [y/N] ' && read ans && [ $${ans:-N} = y ]

# ==============================================================================
# Modules support

## go/deps-reset: reset all the dependencies in the project
.PHONY: go/deps-reset
go/deps-reset: confirm
	git checkout -- go.mod
	go mod tidy
	go mod vendor
## go/tidy: get all new dependencies and vendor them 
.PHONY: go/tidy
go/tidy: confirm
	go mod tidy
	go mod vendor
## go/deps-upgrade: upgrade current dependencies
.PHONY: go/deps-upgrade
go/deps-upgrade: confirm
	# go get $(go list -f '{{if not (or .Main .Indirect)}}{{.Path}}{{end}}' -m all)
	go get -u -t -d -v ./...
	go mod tidy
	go mod vendor

## go/deps-cleancache: clean mod cache
.PHONY: go/deps-cleancache
go/deps-cleancache: confirm
	go clean -modcache

## go/list: list all dependencies in the project
.PHONY: go/list
go/list:
	go list -mod=mod all

# ==============================================================================
# Docker support

## backend/fresh: build all the containers and run them
.PHONY: backend/fresh
backend/fresh:
	docker-compose -f docker-compose-dev.yml build && docker-compose -f docker-compose-dev.yml up -d
## backend/up: run the containers 
.PHONY: backend/up
backend/up:
	docker-compose -f docker-compose-dev.yml up -d
## backend/down: stop containers
.PHONY: backend/down
backend/down: confirm
	docker-compose -f docker-compose-dev.yml down
## docker/prune: remove all containers on the system (not only current project)
.PHONY: docker/prune
docker/prune: confirm
	docker system prune -f
## backend/build: build only the backend container
.PHONY: backend/build
backend/build:
	docker compose -f docker-compose-dev.yml up backend -d --build

# ===============================================================================
# Migrate

## db/migrate name=$1: create a new database migration
.PHONY: db/migrate
db/migrate:
	migrate create -seq -ext=.sql -dir=./migrations ${name}
## db/migrate/up: apply all up database migrations
.PHONY: db/migrate/up
db/migrate/up: confirm
	migrate -path=./migrations -database="${GOGIVE_DB_DSN}" up
## db/migrate/up: apply all down database migrations
.PHONY: db/migrate/down
db/migrate/down: confirm
	migrate -path=./migrations -database="${GOGIVE_DB_DSN}" down 
