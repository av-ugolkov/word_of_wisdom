.PHONY: build_server
build_server:
	go build --gcflags="all=-N -l" -o ./cmd/server/main ./cmd/server/main.go

.PHONY: build_client
build_client:
	go build --gcflags="all=-N -l" -o ./cmd/client/main ./cmd/client/main.go

.PHONY: run.docker
run.docker:
	docker compose -p word_of_wisdom -f deploy/docker-compose.yml up --build --force-recreate

.PHONY: run.tests
run.tests:
	go test ./... -count=1