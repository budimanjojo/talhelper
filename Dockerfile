FROM alpine:3.20.3@sha256:beefdbd8a1da6d2915566fde36db9db0b524eb737fc57cd1367effd16dc0d06d
COPY talhelper /bin/talhelper
ENTRYPOINT ["/bin/talhelper"]
