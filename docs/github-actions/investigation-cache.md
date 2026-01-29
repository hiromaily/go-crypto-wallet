# Github ActionsでCacheが効いていないことの確認

## 1. DockerfileにDebugコードを仕込む

```dockerfile
# Debug for build step
# “before” が常にほぼ0なら、cache-dance / actions/cache のどこかで復元できていない
RUN --mount=type=bind,source=.,target=/src,ro \
    --mount=type=bind,source=.git,target=/src/.git,ro \
    --mount=type=cache,id=go-mod,target=/go/pkg/mod,sharing=locked \
    --mount=type=cache,id=go-build,target=/root/.cache/go-build,sharing=locked \
    set -eux; \
    echo "=== before ==="; \
    du -sh /go/pkg/mod /root/.cache/go-build || true; \
    APP_VERSION="$(git describe --tags --abbrev=0 2>/dev/null || echo dev)"; \
    start="$(date +%s)"; \
    CGO_ENABLED=0 GOOS=linux go build -trimpath \
      -ldflags="-s -w -X ${APP_VERSION_VAR}=${APP_VERSION}" \
      -o /out/app ./${CMD_PATH}; \
    end="$(date +%s)"; \
    echo "go build seconds=$((end-start))"; \
    echo "=== after ==="; \
    du -sh /go/pkg/mod /root/.cache/go-build || true
```

## 2. JOBを実行する

そのとき、Debug用に手動実行できるworkflowを作成しておくこと

## 3. 結果を解析

### 3.1. Restore cache mounts の出力（cache-hit の行）

`actions/cache/restore@v5`の挙動

```
Cache hit for restore-key: Linux-cache-mount-v2-63c8490f30b0d4f9ee8ff92e6456f4b990a71f4afde280c363445895e89dda3f
Received 198 of 198 (100.0%), 0.0 MBs/sec
Cache Size: ~0 MB (198 B)
/usr/bin/tar -xf /home/runner/work/_temp/2168580f-810d-4bb0-bf5a-269da087c620/cache.tzst -P -C /home/runner/work/go-crypto-wallet/go-crypto-wallet --use-compress-program unzstd
Cache restored successfully
Cache restored from key: Linux-cache-mount-v2-63c8490f30b0d4f9ee8ff92e6456f4b990a71f4afde280c363445895e89dda3f
```

これは、`actions/cache は hit してるが、サイズが ~0MB`の状態で、`保存されている cache-mount がほぼ空`

### 3.2 BuildKit Cache Dance の inject/extractのセクション

```log
Post job cleanup.

FROM ghcr.io/containerd/busybox:latest
COPY buildstamp buildstamp
RUN --mount=type=cache,target=/go/pkg/mod     mkdir -p /var/dance-cache/     && cp -p -R /go/pkg/mod/. /var/dance-cache/ || true


FROM ghcr.io/containerd/busybox:latest
COPY buildstamp buildstamp
RUN --mount=type=cache,target=/root/.cache/go-build     mkdir -p /var/dance-cache/     && cp -p -R /root/.cache/go-build/. /var/dance-cache/ || true
```

=> 1sで終わってるので何もしていないようだった。

これは、cache-dance の post-job extract が “id を指定せず” mount している。
extract の内部Dockerfile相当が以下の通り

```dockerfile
RUN --mount=type=cache,target=/go/pkg/mod ...
RUN --mount=type=cache,target=/root/.cache/go-build ...
```

`ここには id= がない。`

BuildKit は id を省略すると、基本的に id=targetパス 相当（少なくとも Dockerfile で id を明示したものとは一致しません）になる。

一方で、実際の Dockerfile は id を明示してる。

これは、以下のミスマッチが起きている

- 本物のキャッシュ：id=go-mod, id=go-build（ここに 953M/354M 溜まった）
- cache-dance が extract してるキャッシュ：id 指定なし（別ID扱い） → 空

**修正ポイント**

cache-map に id を明示して Dockerfile と揃える。
buildkit-cache-dance は cache-map の value をオブジェクトにして target と id を指定できる。

```yml
- name: BuildKit Cache Dance (inject; extract in post)
  uses: reproducible-containers/buildkit-cache-dance@v3
  with:
    builder: ${{ steps.buildx.outputs.name }}
    cache-dir: ${{ env.CACHE_DIR }}
    dockerfile: ${{ env.DOCKERFILE }}
    cache-map: |
      {
        "go-mod":  { "target": "/go/pkg/mod", "id": "go-mod" },
        "go-build":{ "target": "/root/.cache/go-build", "id": "go-build" }
      }
    skip-extraction: false
```
