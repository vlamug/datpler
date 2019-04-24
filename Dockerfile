FROM alpine:latest

COPY bin/ratibor /usr/bin

ENTRYPOINT ["/usr/bin/ratibor"]