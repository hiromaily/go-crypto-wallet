# syntax=docker/dockerfile:1

# without BuildKit cache mounts and bind mounts

############################
# 1) build
# Note: set sha checked by `docker pull golang:1.25.6 && docker inspect`
############################
#FROM golang:1.25.6 AS build
FROM golang:1.25.6@sha256:ce63a16e0f7063787ebb4eb28e72d477b00b4726f79874b3205a965ffd797ab2 AS build
WORKDIR /src

# go mod download (cacheable but extracting cost may be high)
COPY go.mod go.sum ./
RUN go mod download

ARG CMD_PATH=cmd/watch
ARG APP_VERSION=dev

COPY cmd pkg internal ./

# go build
RUN CGO_ENABLED=0 GOOS=linux \
    go build -trimpath \
      -ldflags="-s -w -X main.appVersion=${APP_VERSION}" \
      -o /out/app ./${CMD_PATH}

############################
# 2) runtime: minimal image
############################
FROM gcr.io/distroless/static-debian12:nonroot AS runtime
WORKDIR /
COPY --from=build /out/app /app
USER nonroot:nonroot
ENTRYPOINT ["/app"]
