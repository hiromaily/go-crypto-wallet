# syntax=docker/dockerfile:1

############################
# 1) deps: go mod download via bind mounts
############################
FROM golang:1.25.6 AS deps
WORKDIR /src

# 1. `--mount=type=cache`
# ここで作られるキャッシュは：
# ❌ Docker image のレイヤーには含まれない
# ❌ buildx の type=gha キャッシュにも自動では入らない
# ✅ BuildKit の内部（builder 内）にだけ存在
# つまり、そのままだと GitHub Actions のジョブが終わった瞬間に消える一時データ
#
# 2. `--mount=type=bind,source=.,target=/src,ro`
# Dockerの通常のレイヤーキャッシュ（RUN結果の再利用）はほぼ期待できない
# その代わりに Go側のキャッシュ（cache mount）で速くなる設計
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

# Debug for build step
# “before” が常にほぼ0なら、cache-dance / actions/cache のどこかで復元できていない
# RUN --mount=type=bind,source=.,target=/src,ro \
#     --mount=type=bind,source=.git,target=/src/.git,ro \
#     --mount=type=cache,id=go-mod,target=/go/pkg/mod,sharing=locked \
#     --mount=type=cache,id=go-build,target=/root/.cache/go-build,sharing=locked \
#     set -eux; \
#     echo "=== before ==="; \
#     du -sh /go/pkg/mod /root/.cache/go-build || true; \
#     APP_VERSION="$(git describe --tags --abbrev=0 2>/dev/null || echo dev)"; \
#     start="$(date +%s)"; \
#     CGO_ENABLED=0 GOOS=linux go build -trimpath \
#       -ldflags="-s -w -X ${APP_VERSION_VAR}=${APP_VERSION}" \
#       -o /out/app ./${CMD_PATH}; \
#     end="$(date +%s)"; \
#     echo "go build seconds=$((end-start))"; \
#     echo "=== after ==="; \
#     du -sh /go/pkg/mod /root/.cache/go-build || true

############################
# 3) runtime: minimal image
############################
FROM gcr.io/distroless/static-debian12:nonroot AS runtime
WORKDIR /
COPY --from=build /out/app /app
USER nonroot:nonroot
ENTRYPOINT ["/app"]
