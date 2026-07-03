# syntax=docker/dockerfile:1

# ---- Build stage ----
# Compiles the im-server binary from source (mirrors launcher/build.sh).
FROM golang:1.25-alpine AS builder

# git is needed to fetch the jugglechat-server / imserver-console module deps.
RUN apk add --no-cache git

WORKDIR /src
# Warm the module cache first for faster incremental builds.
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN cd launcher \
    && CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/imserver main.go \
    && cp -rf scripts/* /out/

# ---- Runtime stage ----
FROM alpine:3.20

WORKDIR /opt
RUN apk add --no-cache bash tzdata ca-certificates \
    && mkdir -p /opt/conf /opt/logs

# Binary + run.sh + config_template.yaml produced by the build stage.
COPY --from=builder /out/ /opt/
RUN chmod +x /opt/imserver /opt/run.sh

# WebSocket + API + Navigator + Admin console
EXPOSE 9003 9001 9002 8090
# pprof
EXPOSE 6060

# run.sh renders config_template.yaml from env vars, then launches /opt/imserver.
ENTRYPOINT ["/opt/run.sh"]
