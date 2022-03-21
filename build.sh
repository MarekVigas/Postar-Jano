docker build ./db -t marekvigas/sbb-leto-db
docker push marekvigas/sbb-leto-db:latest
docker build ./admin -t marekvigas/sbb-leto-admin
docker push marekvigas/sbb-leto-admin:latest
docker build ./src/go-api -t marekvigas/sbb-leto-api
docker push marekvigas/sbb-leto-api:latest
docker build ./fe -t marekvigas/sbb-leto-form
docker push marekvigas/sbb-leto-form:latest
