#!/bin/bash
if [ "$(docker network ls | grep mai-test | wc -l)" -eq 0 ]; then
    docker network create mai-test
fi

export TEST_PG_HOST=db-dev
export TEST_PG_PORT=5432
export TEST_PG_USER=postgres
export TEST_PG_PASSWORD=1qaz@WSX
export TEST_PG_DBNAME=mai_webmin_user_service_test
export TEST_PG_SSLMODE=disable

docker run --rm \
  --net mai-test \
  --name ${TEST_PG_HOST} \
  --hostname ${TEST_PG_HOST} \
  -e POSTGRES_USER=${TEST_PG_USER} \
  -e POSTGRES_PASSWORD=${TEST_PG_PASSWORD} \
  -e POSTGRES_DB=${TEST_PG_DBNAME} \
  -d \
  postgres:alpine

docker run --rm \
  --net mai-test \
  postgres:alpine \
  ping -w 10000 -c 1 ${TEST_PG_HOST}


docker run --rm \
    -v "$(pwd)/deployments/migrator/migrations:/migrations" \
    --net mai-test \
    migrate/migrate \
    -path=/migrations/ \
    -database postgres://${TEST_PG_USER}:${TEST_PG_PASSWORD}@${TEST_PG_HOST}:${TEST_PG_PORT}/${TEST_PG_DBNAME}?sslmode=${TEST_PG_SSLMODE} up

if [ $? -ne 0 ]; then
    docker stop ${TEST_PG_HOST}
    exit
fi

docker run --rm \
  --net mai-test \
  -e TEST_PG_HOST=${TEST_PG_HOST} \
  -e TEST_PG_PORT=${TEST_PG_PORT} \
  -e TEST_PG_USER=${TEST_PG_USER} \
  -e TEST_PG_PASSWORD=${TEST_PG_PASSWORD} \
  -e TEST_PG_DBNAME=${TEST_PG_DBNAME} \
  -e TEST_PG_SSLMODE=${TEST_PG_SSLMODE} \
  -e GO_ENV=test \
  -e GO_SEED_PG_HOST=${TEST_PG_HOST} \
  -e GO_SEED_PG_PORT=${TEST_PG_PORT} \
  -e GO_SEED_PG_USER=${TEST_PG_USER} \
  -e GO_SEED_PG_PASSWORD=${TEST_PG_PASSWORD} \
  -e GO_SEED_PG_DBNAME=${TEST_PG_DBNAME} \
  -e GO_SEED_PG_SCHEMA=back_office \
  -e GO_SEED_PG_SSLMODE=${TEST_PG_SSLMODE} \
  -e GO_SEED_SOURCE_PATH=deployments/migrator/seed-data \
  -e CGO_ENABLED=0 \
  test-image \
  sh build/ci/ci-integration-test-docker-script.sh

docker stop ${TEST_PG_HOST}
