ARG VERSION

## ================================================================================================
## Builder Stage -> creating the binary
## ================================================================================================
FROM golang:1.22.3-alpine3.18 as builder
ARG VERSION

WORKDIR /build
COPY . .
RUN go build -ldflags="-s -w -X github.com/budimanjojo/talhelper/v3/cmd.version=${VERSION}" -o /usr/local/bin/talhelper


## ================================================================================================
## Serving/Production Stage
## ================================================================================================
FROM scratch
COPY --from=builder /usr/local/bin/talhelper /usr/local/bin/talhelper
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
WORKDIR /config
ENTRYPOINT [ "/usr/local/bin/talhelper" ]
