FROM --platform=${BUILDPLATFORM:-linux/amd64} ghcr.io/openfaas/license-check:0.4.1 as license-check

FROM --platform=${BUILDPLATFORM:-linux/amd64} golang:1.20 as build

ENV GO111MODULE=on
ENV CGO_ENABLED=1

ARG TARGETPLATFORM
ARG BUILDPLATFORM
ARG TARGETOS
ARG TARGETARCH

ARG GIT_COMMIT
ARG VERSION

COPY --from=license-check /license-check /usr/bin/

WORKDIR /go/src/github.com/forge4flow/flow-events-connector

COPY .        .

RUN license-check -path ./ --verbose=false "BoiseITGuru" "Forge4Flow Authors" "Forge4Flow Author(s)"

# Run a gofmt and exclude all vendored code.
# TODO: Fix Testing
# RUN test -z "$(gofmt -l $(find . -type f -name '*.go' -not -path "./vendor/*"))"
# RUN go test $(go list ./... | grep -v integration | grep -v /vendor/ | grep -v /template/) -cover

RUN CGO_ENABLED=${CGO_ENABLED} GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -a -installsuffix cgo -o flow-events-connector main.go

FROM --platform=${TARGETPLATFORM:-linux/amd64} alpine:3.18.3 as ship

LABEL org.label-schema.license="MIT" \
    org.label-schema.vcs-url="https://github.com/forge4flow/flow-events-connector" \
    org.label-schema.vcs-type="Git" \
    org.label-schema.name="forge4flow/flow-events-connector" \
    org.label-schema.vendor="forge4flow" \
    org.label-schema.docker.schema-version="1.0"

LABEL org.opencontainers.image.source=https://github.com/forge4flow/flow-events-connector

RUN addgroup -S app \
    && adduser -S -g app app \
    && apk add --no-cache ca-certificates

WORKDIR /home/app

COPY --from=build /go/src/github.com/forge4flow/flow-events-connector    .

RUN chown -R app:app ./

USER app

CMD ["./flow-events-connector"]