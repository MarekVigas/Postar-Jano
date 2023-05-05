# Local setup

```bash
export $(grep -v '^#' ../../api.env | xargs)
make build
./registration_api
```

Unset env variables.
```bash
unset $(grep -v '^#' ../../api.env | sed -E 's/(.*)=.*/\1/' | xargs)
```