# Github ActionsでCacheが効いていないことの確認

---

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

---

## 2. JOBを実行する (1回目)

そのとき、Debug用に手動実行できるworkflowを作成しておくこと

---

## 3. 結果を解析 (1回目)

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
  - これは、Dockerfile内の、ビルド時に走らせる `du -sh /go/pkg/mod /root/.cache/go-build || true` からわかる。
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

---

## 4. JOBを実行する (2回目)

---

## 5. 結果を解析 (2回目)

### 5.1. Restore cache mounts の出力（cache-hit の行）

この時点で改善していない。

### 5.2 BuildKit Cache Dance の inject/extractのセクション

```
Post job cleanup.

FROM ghcr.io/containerd/busybox:latest
COPY buildstamp buildstamp
RUN --mount=type=cache,target=/go/pkg/mod,id=go-mod     mkdir -p /var/dance-cache/     && cp -p -R /go/pkg/mod/. /var/dance-cache/ || true


FROM ghcr.io/containerd/busybox:latest
COPY buildstamp buildstamp
RUN --mount=type=cache,target=/root/.cache/go-build,id=go-build     mkdir -p /var/dance-cache/     && cp -p -R /root/.cache/go-build/. /var/dance-cache/ || true
```

[改善ポイント] 68sかかっているため、動作していると思われる。

---

## 6. キャッシュが効いていない原因 (ここまでにおいて)

`actions/cache/save` が `cache-dance` の `extraction（post処理）` より先に走ってしまっている のが原因

ここでいう`extraction`とは、`RUN --mount=type=cache` で使われたキャッシュの中身を、BuildKit の外に“取り出して（extract して）”、`actions/cache` で保存できる形に変換する処理。

**根拠**

- restore は毎回 “~0MB (195B)” を復元している（＝保存されたものが空）
- なのに cache-dance の inject/extract が 68秒かかる（＝extract で大容量コピーしているっぽい）
- Dockerfile の before が 8K（＝inject で中身が入っていない）

**なぜそうなるのか**

- `buildkit-cache-dance@v3` の `extraction` は `post step` で走る
  - ログも Post job cleanup. に出ている

それに対して利用したworkflowは

1. actions/cache/restore（普通のステップ）
2. cache-dance（普通ステップ＋postでextract）
3. ビルド
4. actions/cache/save（普通のステップ）

普通のステップとしての `actions/cache/save` は、post cleanup より先に実行される。
つまり cache-dance が cache-dir に中身を書き出す前に、空の cache-dir を保存してしまい、結果として ~0MB (195B) のキャッシュだけが永続化され続ける。

**対策**

- `restore/save 分割をやめて actions/cache@v5 を使う。`
  - `actions/cache@v5`（restore+save一体型）は save が post で走るので、post の順序を利用できる
  - 重要なのは **postの実行順が“逆順”**ということ
  - cache action を cache-dance より前に置く
    - job 終了時、cache-dance の post（extract）が先に走り、その後に cache action の post（save）が走る

**README の公式例は actions/cache の“一体型”を使っている**

---

## 7. コード変更後、JOBを実行した結果を分析 (3回目)

`Post Cache mounts (restore + save in post)`のセクションにて以下のlogが出力された。

```
Post Cache mounts (restore + save in post)
Post job cleanup.
Cache hit occurred on the primary key Linux-cache-mount-v2-26f17f8f1f45eb19e54c7ec2dce8d9c019e1c8032853d7fe53d6fe19e51dfd43, not saving cache.
```

`actions/cache@v5` が “primary key で cache hit したから保存しない” という意味。
つまり 既に存在している（でも中身が空の）キャッシュを primary key で拾ってしまっていて、上書き保存できない状態である。

`actions/cache` は仕様として、primary key で当たった場合は save をスキップします（同じ key に上書きできない）。

決定的なのが、Dockerfile内の`go build`時、Dockerfile の before が 8.0K のまま。

**解決策：空キャッシュを捨てて「新しい key で保存」させる**

1. 最短（確実）：CACHE_VERSION を bump して新規キーにする
2. restore と save の key を分ける
   actions/cache@v5 は「primary key で hit すると保存しない」ので、restore は固定key、save は run_id を入れたユニークkeyにする、という回避ができるが、actions/cache@v5 は restore+save一体なので、この方式をやるならまた restore/save 分離に戻す必要があり、かつ post順序の問題も再燃するので、このやり方は不可能。`buildkit-cache-dance@v3`を使う場合は実現できないということ。

## 8. バージョン変更後、JOBを実行する (4回目)

もう一度`V3`で実行したが、ビルド時のbeforeのサイズが依然小さい。

**v3に保存されたキャッシュの中身が本当に大きいか？」を確定すること**

疑う原因

- v3キーに保存されたキャッシュが “実は空”
- cache-dance が “正しい cache id” を extract/inject できていない

```
- name: Inspect cache-dir before post steps
  run: |
    set -eux
    ls -lah "${CACHE_DIR}" || true
    du -sh "${CACHE_DIR}" || true
    find "${CACHE_DIR}" -maxdepth 2 -type f | head -50 || true
```

これのlogが、以下だった。

```
+ ls -lah cache-mount total 8.0K
drwxr-xr-x 2 runner runner 4.0K Jan 29 08:24 .
drwxr-xr-x 26 runner runner 4.0K Jan 29 11:30 .. + du -sh cache-mount 4.0K cache-mount
```

**結論**

cache-dance の inject/extract が “cache-dir に何も書いていない” ため、actions/cache@v5 が保存/復元している cache-mount/ が常に空（4KB）のままである。

その結果、

- restore は primary key で hit（でも中身は空）
- build の before は 8KB（注入されない）
- Go は毎回 download して 40秒台

という状態になります。

**「なぜ cache-dance が cache-dir に書かないのか？」**

1. cache-map が解釈されていない（＝結局 idなし mount を触ってる / 何もコピーしてない）
   以前あなたが貼ってくれた cache-dance の post は、id= が無い mount でした。もし今も内部では id 無しの別キャッシュ（空）を触っていると、cache-dir には当然何も出ません。
2. そもそも cache-dance が参照している Dockerfile が違う（or パス解決の問題）
   `dockerfile: ${{ env.DOCKERFILE }}` を渡していますが、cache-dance は Dockerfile を解析して cache-map と合わせるので、違うDockerfileを見ていると何もしないことがあります。
3. cache-dance の “extract 先” が cache-dir 配下ではない（入力ミス/相対パスのズレ）
   cache-dir: cache-mount は相対パスで、workspace基準で合ってるはずですが、念のため絶対パスにして潰せます。

**いったん id 指定はやめる（READMEの基本形に戻す）**

いまの症状は「cache-map が解釈されてない/動いてない」可能性が高い。まず README の最小構成で “動くこと” を確認するのが早い。

=> cache-dance が id 指定無しで mount してくるなら、Dockerfile 側も合わせる。

この問題は[Issue#33](https://github.com/reproducible-containers/buildkit-cache-dance/issues/33)で報告されている。

BuildKit Cache Dance を使ってもキャッシュが取り出せず、保存されたキャッシュが空になってしまう
→ 結果として復元しても意味がない（前回の成果物が注入されない）

また、[複数のキャッシュを使うと上書きされる](https://github.com/reproducible-containers/buildkit-cache-dance/issues/39)という問題もある。

---

## [原点] 問題点の整理**

`RUN --mount=type=cache` の中身は、`buildx cache-to/from (type=gha 等)` だけでは永続化されない（＝昔から課題）というのが広く知られていて、Docker側/BuildKit側にも関連Issueがあります。
その上で、代替案は「mount cache 自体をやめる」か「mount cache を外へ持ち出す別手段を使う」の2系統となる。

### RUN --mount=type=cache をやめて、レイヤーキャッシュ + buildx cache (type=gha) に寄せる

- `COPY go.mod go.sum` → `RUN go mod download`
- `COPY .` → `RUN go build`

**問題**

ソースが少しでも変わった時に `RUN go build` は“実行し直し”になりやすい

**キーポイント**

Go が速くなる鍵は GOCACHE（コンパイルキャッシュ） と GOMODCACHE（モジュールキャッシュ）

- GOMODCACHE は「COPY go.mod go.sum → RUN go mod download」の層に分けると、ソースが変わっても温存できる（＝依存再DLを避けられる）
- 一方 GOCACHE は、Dockerビルドを跨いで永続化されないと “前回のコンパイル成果” を使えない。
