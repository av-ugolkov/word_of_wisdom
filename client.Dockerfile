FROM golang:1.21.0-alpine AS builder

RUN apk --no-cache --update --upgrade add git make

WORKDIR /build

COPY . .

RUN --mount=type=cache,target=/go make build_client

FROM scratch
LABEL key="WoW Client"

COPY --from=builder /build/config/config.yaml /config/config.yaml
COPY --from=builder /build/cmd/client/main /

ENTRYPOINT ["./main"]