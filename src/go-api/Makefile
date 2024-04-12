
build:
	CGO_ENABLED=0 go build \
		-o "./registrations_api" \
		"./cmd"

start-test-docker:
	docker run -p 5432:5432 -e POSTGRES_PASSWORD=postgres -d postgres

test:
	go test ./...

.PHONY: build test start-test-docker
