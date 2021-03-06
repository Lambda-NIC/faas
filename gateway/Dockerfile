FROM golang:1.10.4 as build
ARG GIT_COMMIT_SHA
ARG GIT_COMMIT_MESSAGE
ARG VERSION='dev'

WORKDIR /go/src/github.com/Lambda-NIC/faas/gateway

COPY vendor         vendor

COPY handlers       handlers
COPY metrics        metrics
COPY requests       requests
COPY tests          tests

COPY types          types
COPY queue          queue
COPY plugin         plugin
COPY version        version
COPY scaling        scaling
COPY server.go      .

# Run a gofmt and exclude all vendored code.
RUN go test $(go list ./... | grep -v integration | grep -v /vendor/ | grep -v /template/) -cover \
    && CGO_ENABLED=0 GOOS=linux go build --ldflags "-s -w \
    -X github.com/Lambda-NIC/faas/gateway/version.GitCommitSHA=${GIT_COMMIT_SHA}\
    -X \"github.com/Lambda-NIC/faas/gateway/version.GitCommitMessage=${GIT_COMMIT_MESSAGE}\"\
    -X github.com/Lambda-NIC/faas/gateway/version.Version=${VERSION}" \
    -a -installsuffix cgo -o gateway .

FROM alpine:3.8

LABEL org.label-schema.license="MIT" \
    org.label-schema.vcs-url="https://github.com/Lambda-NIC/faas" \
    org.label-schema.vcs-type="Git" \
    org.label-schema.name="openfaas/faas" \
    org.label-schema.vendor="openfaas" \
    org.label-schema.docker.schema-version="1.0"

RUN addgroup -S app \
    && adduser -S -g app app

WORKDIR /home/app

EXPOSE 8080
ENV http_proxy      ""
ENV https_proxy     ""

COPY --from=build /go/src/github.com/Lambda-NIC/faas/gateway/gateway    .
COPY assets     assets

RUN chown -R app:app ./

USER app

CMD ["./gateway"]
