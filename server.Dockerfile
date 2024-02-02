FROM golang:1.21.0-alpine AS builder

RUN apk --no-cache --update --upgrade add git make

WORKDIR /build

COPY . .

RUN --mount=type=cache,target=/go make build_server

FROM scratch
LABEL key="WoW Server"

COPY --from=builder /build/config/config_docker.yaml /config/config.yaml
COPY --from=builder /build/cmd/server/main /

EXPOSE 5555

ENTRYPOINT ["./main"]