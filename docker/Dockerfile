# syntax=docker/dockerfile:1

############################
# 1) deps: go mod download via bind mounts
############################
FROM golang:1.25.6 AS deps
WORKDIR /src

RUN --mount=type=bind,source=go.mod,target=/src/go.mod,ro \
    --mount=type=bind,source=go.sum,target=/src/go.sum,ro \
    --mount=type=cache,id=go-mod,target=/go/pkg/mod,sharing=locked \
    go mod download


############################
# 2) build: resolve git tag + build selected cmd
############################
FROM golang:1.25.6 AS build

# Which cmd to build (e.g. cmd/watch)
ARG CMD_PATH=cmd/watch

# ldflags variable path (adjust if needed)
ARG APP_VERSION_VAR=main.appVersion

WORKDIR /src

RUN --mount=type=bind,source=.,target=/src,ro \
    --mount=type=bind,source=.git,target=/src/.git,ro \
    --mount=type=cache,id=go-mod,target=/go/pkg/mod,sharing=locked \
    --mount=type=cache,id=go-build,target=/root/.cache/go-build,sharing=locked \
    APP_VERSION="$(git describe --tags --abbrev=0 2>/dev/null || echo dev)" && \
    echo "Building ${CMD_PATH} with version=${APP_VERSION}" && \
    CGO_ENABLED=0 GOOS=linux \
    go build -trimpath \
      -ldflags="-s -w -X ${APP_VERSION_VAR}=${APP_VERSION}" \
      -o /out/app ./${CMD_PATH}


############################
# 3) runtime: minimal image
############################
FROM gcr.io/distroless/static-debian12:nonroot AS runtime
WORKDIR /
COPY --from=build /out/app /app
USER nonroot:nonroot
ENTRYPOINT ["/app"]
