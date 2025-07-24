## Running tests

By default, tests start a postgres test container. If you want to run them
using an existing database specify the following env vars:

```
POSTGRES_PASSWORD=password
POSTGRES_USER=postgres
POSTGRES_HOST=db
POSTGRES_PORT=5432
POSTGRES_DB=postgres
```