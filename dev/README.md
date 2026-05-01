# Local development

In order to start local development setup run `make dc-up`. To run backend component
test call `make test`.

The BE sources are mounted to the container i.e.: backend-api executable 
can be rebuilt by running `make rebuild-api` and started by `make run-api`. The API
is served on port 5000 by default and forwarded to host on port 48080.

To create an admin user run

```bash
docker exec -i -t dev-api-1 /app/registrations_api add-user --username admin@sbb.sk --password Pass123
```

## Hot-reload (possible improvement)

For faster development turnaround, [Air](https://github.com/air-verse/air) can be used to watch for Go source changes and automatically rebuild and restart the API inside the container — eliminating the need to run `make rebuild-api` manually.

**Setup:**

1. Install Air in `Dockerfile` (`devImage` stage):
```dockerfile
RUN go install github.com/air-verse/air@latest
```

1. Add `.air.toml` in `src/go-api`:
```toml
[build]
  cmd = "make build && cp registrations_api /app/registrations_api"
  bin = "/app/registrations_api"

[watch]
  include_ext = ["go"]
  exclude_dir = ["vendor"]
```

1. Change the `api` service command in `docker-compose.yml`:
```yaml
command: "air"
```

After that, any `.go` file change triggers an automatic rebuild and restart (~1-2s) with no manual step.

## Database

Migrations are applied on start and by default database port is not forwarded to localhost.

### Backup & Restore
```bash
docker exec -t dev-db-1 pg_dump -U postgres > dump_`date +%d-%m-%Y"_"%H_%M_%S`.sql
```

```bash
cat your_dump.sql | docker exec -i dev-db-1 psql -U postgres
```