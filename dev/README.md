# Local development

In order to start local development setup run `make dc-up`. To run backend component
test call `make test`.

The BE sources are mounted to the container i.e.: backend-api executable 
can be rebuilt by running `make rebuild-api` and started by `make run-api`. The API
is served on port 5000 by default and forwarded to host on port 48080.

## Database

Database is running on default port 5432

### Backup & Restore
```bash
docker exec -t your-db-container pg_dumpall -c -U postgres > dump_`date +%d-%m-%Y"_"%H_%M_%S`.sql
```

```bash
cat your_dump.sql | docker exec -i your-db-container psql -U postgres
```