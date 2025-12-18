FROM alpine:3.23.2@sha256:c93cec902b6a0c6ef3b5ab7c65ea36beada05ec1205664a4131d9e8ea13e405d
ARG TARGETPLATFORM
COPY --chmod=755 $TARGETPLATFORM/talhelper /bin
ENTRYPOINT ["/bin/talhelper"]
