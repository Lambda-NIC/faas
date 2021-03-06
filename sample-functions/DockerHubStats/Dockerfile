FROM golang:1.9.7-alpine as builder

MAINTAINER alex@openfaas.com
ENTRYPOINT []

RUN apk --no-cache add make curl \
    && curl -sL https://github.com/Lambda-NIC/faas/releases/download/0.9.6/fwatchdog > /usr/bin/fwatchdog \
    && chmod +x /usr/bin/fwatchdog

WORKDIR /go/src/github.com/Lambda-NIC/faas/sample-functions/DockerHubStats

COPY . /go/src/github.com/Lambda-NIC/faas/sample-functions/DockerHubStats

RUN make install

FROM alpine:3.8

# Needed to reach the hub
RUN apk --no-cache add ca-certificates 

COPY --from=builder /usr/bin/fwatchdog  /usr/bin/fwatchdog
COPY --from=builder /go/bin/DockerHubStats  /usr/bin/DockerHubStats
ENV fprocess "/usr/bin/DockerHubStats"

HEALTHCHECK --interval=5s CMD [ -e /tmp/.lock ] || exit 1
CMD ["/usr/bin/fwatchdog"]
