FROM golang:1.15-alpine
ARG SRV_UID=1000
ARG SRV_GID=1000
ENV SERVICE_NAME users-service

COPY build/assets/ghz /usr/local/bin
COPY build/assets/repositories /etc/apk/repositories

RUN chmod 644 /etc/apk/repositories && \
    apk update && apk add --no-cache sudo git openssh-client
RUN apk add --no-cache make postgresql-client curl tzdata && \
    chmod a+x /usr/local/bin/ghz

COPY build/assets/sudoers /etc/sudoers

RUN chmod 0400 /etc/sudoers && \
    addgroup anggri -g ${SRV_GID} && \
    addgroup sudo && \
    adduser -h /home/anggri -s /bin/sh -u ${SRV_UID} -G users -D anggri && \
    adduser anggri anggri && \
    adduser anggri sudo

USER anggri
WORKDIR "/go/src/github.com/anggri-microservice/${SERVICE_NAME}"

COPY .netrc /tmp/.netrc

COPY build/dev/evans_linux_amd64.tar.gz /tmp

RUN sudo chown -R anggri:anggri /go && \
    echo "Installing additional tools, please wait..." && \
    go get -u golang.org/x/lint/golint && echo "- Success installing golint" && \
    go get github.com/githubnemo/CompileDaemon && echo "- Success installing CompileDaemon" && \
    cp /tmp/.netrc ~/.netrc && \
    chmod 600 ~/.netrc && \
    export GOPRIVATE=github.com/anggri-microservice/* && \
    sudo tar -C /usr/bin -xzf /tmp/evans_linux_amd64.tar.gz && \
    sudo rm /tmp/evans_linux_amd64.tar.gz && \
    export GO111MODULE=off && \
    go get -tags 'postgres' -u github.com/golang-migrate/migrate/cmd/migrate && echo "- Success installing golang-migrate" && \
    go get github.com/bayucandra/go-seed-pg && echo "- Success installing go-seed-pg" && \
    unset GO111MODULE

ENTRYPOINT ["sh", "./build/dev/entrypoint.sh"]
