FROM golang:1.13-alpine AS builder
ENV SERVICE_NAME users-service

WORKDIR "/go/src/gitlab.com/pt-mai/${SERVICE_NAME}"
COPY ./go.mod .

COPY ./id_rsa /tmp/id_rsa
# RUN go mod download
# RUN go mod verify

RUN apk update && apk add --no-cache git openssh-client && \
    chmod 600 /tmp/id_rsa && \
    mkdir ~/.ssh && \
    mv /tmp/id_rsa ~/.ssh/id_rsa && \
    eval $(ssh-agent) && \
    echo -e "StrictHostKeyChecking no" >> /etc/ssh/ssh_config && \
    ssh-add ~/.ssh/id_rsa && \
    git config --global url."git@gitlab.com:".insteadOf "https://gitlab.com/" && \
    export GOPRIVATE=gitlab.com/pt-mai && \
    go mod vendor
COPY . .
RUN GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /usr/bin/${SERVICE_NAME} ./cmd/$SERVICE_NAME/...



FROM alpine
#COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /usr/bin/${SERVICE_NAME} /usr/bin/${SERVICE_NAME}
RUN adduser -D -g '' mai
USER mai
ENTRYPOINT ["/usr/bin/users-service"]