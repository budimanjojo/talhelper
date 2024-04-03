FROM golang:1.22.2-alpine3.18 as builder
WORKDIR /build
COPY . .
RUN go build -o /usr/local/bin/talhelper

FROM scratch
COPY --from=builder /usr/local/bin/talhelper /usr/local/bin/talhelper
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
WORKDIR /config
ENTRYPOINT [ "/usr/local/bin/talhelper" ]
