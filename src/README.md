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
