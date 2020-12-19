#!/bin/bash

pwd
go-seed-pg && \
    go test ./cmd/... ./internal/... -v -coverprofile .coverage-integration.txt -tags=integration && \
    go tool cover -func .coverage-integration.txt
