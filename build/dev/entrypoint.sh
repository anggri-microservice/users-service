#!/bin/bash
#export DOCKER_HOST=$(ip -4 addr show docker0 | grep -Po 'inet \K[\d.]+')

cp .netrc ~/
chmod 600 ~/.netrc

export DOCKER_HOST="host.docker.internal"
export SRV_DOT_ENV="true"
ping -q -c1 $DOCKER_HOST > /dev/null 2>&1
if [ $? -ne 0 ]; then
  HOST_IP=$(ip route | awk 'NR==1 {print $3}')
  echo "${HOST_IP}	${DOCKER_HOST}" | sudo tee -a /etc/hosts > /dev/null
fi

sudo chown -R anggri:anggri /go/pkg
sudo chown -R anggri:anggri ./vendor

export GOPRIVATE="github.com/anggri-microservice/*"
export CGO_ENABLED="0"

if [ "${COMPOSE_PROJECT_NAME}" = "emporia" ]; then
    make go-mod-replace-clear
    make go-mod-replace-insert
fi

echo "SRV_DC_ONLINE = $SRV_DC_ONLINE"
if [ "${SRV_DC_ONLINE}" = "true" ]; then
    echo "Vendoring, please wait..."
    go mod vendor
fi

CompileDaemon \
    -color=true \
    -graceful-kill=true \
    -pattern="^(\.env.+|\.env)|(.+\.go|.+\.c)$" \
    -build="go build -mod=vendor -o $SERVICE_NAME ./cmd/$SERVICE_NAME/..." \
    -command="./${SERVICE_NAME}"

if [ "${COMPOSE_PROJECT_NAME}" = "emporia" ]; then
    make go-mod-replace-clear
fi
