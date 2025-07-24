# Local development

In order to start local development setup run `make dc-up`. To run backend component
test call `make test`.

The BE sources are mounted to the container i.e.: backend-api executable 
can be rebuilt by running `make rebuild-api` and started by `make run-api`. The API
is served on port 5000 by default and forwarded to host on port 48080.

To create an admin user run

```bash
docker exec -i -t dev-api-1 /src/registrations_api add-user --username admin@sbb.sk --password Pass123
```

## Database

Migrations are applied on start and by default database port is not forwarded to localhost.

### Backup & Restore
```bash
docker exec -t dev-db-1 pg_dump -U postgres > dump_`date +%d-%m-%Y"_"%H_%M_%S`.sql
```

```bash
cat your_dump.sql | docker exec -i dev-db-1 psql -U postgres
```