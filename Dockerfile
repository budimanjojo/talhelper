FROM alpine:3.22.2@sha256:265b17e252b9ba4c7b7cf5d5d1042ed537edf6bf16b66130d93864509ca5277f
COPY talhelper /bin/talhelper
ENTRYPOINT ["/bin/talhelper"]
