
build:
	CGO_ENABLED=0 go build \
		-o "./registrations_api" \
		"./cmd"

run:
	./registrations_api

migrate:
	./registrations_api migrate

export-env:
	export $(grep -v '^#' api.env | xargs)

start-test-docker:
	docker run -p 5432:5432 -e POSTGRES_PASSWORD=password -d postgres

test:
	go test ./...

.PHONY: build test start-test-docker