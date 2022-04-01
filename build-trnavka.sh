docker build ./db -t marekvigas/sbb-leto-db:trnavka
docker push marekvigas/sbb-leto-db:trnavka
docker build ./admin -t marekvigas/sbb-leto-admin:trnavka
docker push marekvigas/sbb-leto-admin:trnavka
docker build ./src/go-api -t marekvigas/sbb-leto-api:trnavka
docker push marekvigas/sbb-leto-api:trnavka
docker build ./fe --build-arg REACT_APP_API_HOST=https://tabory.trnavka.sk --build-arg REACT_APP_RESULT_REDIRECT=https://trnavka.sk/tabory -t marekvigas/sbb-leto-form:trnavka
docker push marekvigas/sbb-leto-form:trnavka
