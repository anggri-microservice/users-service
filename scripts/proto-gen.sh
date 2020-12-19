#!/bin/bash

PROTO_PATHS="
    test/test
" 

for PROTO_PATH in ${PROTO_PATHS};
do
    PROTO_PATH_OUTPUT=pb/${PROTO_PATH}.pb.go && \
        protoc -I pb/ \
        pb/${PROTO_PATH}.proto \
        --go_out=plugins=grpc:pb && \
        sed s/,omitempty// ${PROTO_PATH_OUTPUT} > ${PROTO_PATH_OUTPUT}.tmp && mv ${PROTO_PATH_OUTPUT}.tmp ${PROTO_PATH_OUTPUT} && \
        protoc-go-inject-tag -input=${PROTO_PATH_OUTPUT}
done


#ls proto | awk '{print "protoc --proto_path=proto proto/"$1"/*.proto --go_out=plugins=grpc:."}' | sh

#mockgen -source pb/auth/auth.pb.go -destination pb/mock_authpb/mock_authpb.go