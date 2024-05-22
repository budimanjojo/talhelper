# Additional build flags passed to the `go build` command
ARG BUILD_FLAGS


## ================================================================================================
## Builder Stage -> creating the binary
## ================================================================================================
FROM golang:1.22.3-alpine3.18 as builder
ARG BUILD_FLAGS
WORKDIR /build
COPY . .
RUN go build ${BUILD_FLAGS} -o /usr/local/bin/talhelper



## ================================================================================================
## Serving/Production Stage
## ================================================================================================
FROM scratch
COPY --from=builder /usr/local/bin/talhelper /usr/local/bin/talhelper
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
WORKDIR /config
ENTRYPOINT [ "/usr/local/bin/talhelper" ]
