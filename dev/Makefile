
build-images:
	docker compose build

dc-up:
	docker compose up -d

dc-down:
	docker compose down

rebuild-api: dc-up
	docker compose exec api sh -c "make build && cp registrations_api /app/registrations_api"
	docker compose restart api

run-api: dc-up
	docker compose exec -i -t api make run

test: dc-up
	docker compose exec -i -t api make test

migrate: dc-up
	docker compose exec -i -t api make migrate
