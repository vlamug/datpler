FROM alpine:latest

COPY bin/metricplower /usr/bin

ENTRYPOINT ["/usr/bin/metricplower"]