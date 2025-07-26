## Project structure

- api package - registers all handlers and call services
- config package - defines runtime config option for the application
- model package - defines entities exchanges between services and repository
- resources package - defines entities exchanged between handlers and services
- repository package - implements DB queries
- services package - implements the business logic, calls repository layer

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