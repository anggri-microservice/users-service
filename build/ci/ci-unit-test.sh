#!/bin/bash

docker run --rm \
  -e TEST_PG_HOST=${TEST_PG_HOST} \
  -e TEST_PG_PORT=${TEST_PG_PORT} \
  -e TEST_PG_USER=${TEST_PG_USER} \
  -e TEST_PG_PASSWORD=${TEST_PG_PASSWORD} \
  -e TEST_PG_DBNAME=${TEST_PG_DBNAME} \
  -e TEST_PG_SSLMODE=${TEST_PG_SSLMODE} \
  -e CGO_ENABLED=0 \
  test-image \
  go test ./... -v -coverprofile .coverage-unit.txt -tags=unit && go tool cover -func .coverage-unit.txt
