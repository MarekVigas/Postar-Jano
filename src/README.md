# Local setup

Environment variables
```
cat api.env

POSTGRES_PASSWORD=password
POSTGRES_USER=postgres
POSTGRES_HOST=localhost
POSTGRES_PORT=5432
POSTGRES_DB=postgres
MAILGUN_DOMAIN=http://bla.bla
MAILGUN_KEY=secret
MAILGUN_EU_BASE=True
JWT_SECRET=secret
CONFIRMATION_EMAIL_TEMPLATE=/app/bin/templates/confirmation.html
PROMO_EMAIL_TEMPLATE=/app/bin/templates/templates/promo.html
PROMO_SECRET=secret
PROMO_SIMPLE=True
HOST="0.0.0.0"
PORT=5000
# PROMO_ACTIVATION_DATE="2006-01-02T15:04:05+02:00"
# PROMO_EXPIRATION_DATE="2006-01-02T15:04:05+02:00"

```

```bash
export $(grep -v '^#' api.env | xargs)
make build
./registration_api
```

Unset env variables.
```bash
unset $(grep -v '^#' api.env | sed -E 's/(.*)=.*/\1/' | xargs)
```

Create user
```bash
./registrations_api add-user --username test --password nbusr123
2023-05-06T18:04:57.315+0200    INFO    command/create_user.go:72       User created successfully.
```

Login
```bash
curl -X POST localhost:50000/api/sign/in -H "Accept: application/json"  -H "Content-Type: application/json"  -d '{"username":"test", "password":"nbusr123"}'

{"token":"<token>"}

```

Create promo code + send email
```bash
curl -X POST localhost:50000/api/promo_codes -H "Authorization: Bearer <JWT_TOKEN>" -H "Content-Type: application/json" -d '{"email":"test@mailinator.com", "registration_count":1, "send_email":true}'
```
