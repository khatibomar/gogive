# ==============================================================================
# Modules support

deps-reset:
	git checkout -- go.mod
	go mod tidy
	go mod vendor

tidy:
	go mod tidy
	go mod vendor

deps-upgrade:
	# go get $(go list -f '{{if not (or .Main .Indirect)}}{{.Path}}{{end}}' -m all)
	go get -u -t -d -v ./...
	go mod tidy
	go mod vendor

deps-cleancache:
	go clean -modcache

list:
	go list -mod=mod all

# ==============================================================================
# Docker support

docker-dev-build-up:
	docker-compose -f docker-compose-dev.yml build && docker-compose -f docker-compose-dev.yml up -d
docker-dev-up:
	docker-compose -f docker-compose-dev.yml up -d
docker-dev-down:
	docker-compose -f docker-compose-dev.yml down
docker-clean:
	docker system prune -f
