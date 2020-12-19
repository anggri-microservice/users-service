#!/bin/bash
#going to root project dir

if [[ "$(which golint | wc -l)" -ne 1 ]]; then
    export GO111MODULES=off
    go get golang.org/x/lint/golint
    unset GO111MODULES
fi

if [[ "$(which CompileDaemon | wc -l)" -ne 1 ]]; then
    export GO111MODULES=off
    go get github.com/githubnemo/CompileDaemon
    unset GO111MODULES
fi

echo "Start linting========"
CompileDaemon -build="echo true" \
    -command="golint ./cmd/... ./internal/..."

echo ""
echo "End Linting********"
