FROM alpine:3.23.3@sha256:eb37f58646a901dc7727cf448cae36daaefaba79de33b5058dab79aa4c04aefb
ARG TARGETPLATFORM
COPY --chmod=755 $TARGETPLATFORM/talhelper /bin
ENTRYPOINT ["/bin/talhelper"]
