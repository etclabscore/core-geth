FROM golang:1.25-bookworm

RUN apt-get update && apt-get install -y --no-install-recommends \
      ca-certificates git build-essential python3 \
    && rm -rf /var/lib/apt/lists/*

COPY . /src
WORKDIR /src

RUN set -e; \
    for sub in tests/testdata:https://github.com/ethereum/tests \
               tests/testdata-etc:https://github.com/etclabscore/tests; do \
      dir="${sub%%:*}"; url="${sub#*:}"; \
      if [ -z "$(ls -A "$dir" 2>/dev/null)" ]; then \
        rev="$(git ls-tree HEAD "$dir" | awk '{print $3}')"; \
        rm -rf "$dir" && git init -q "$dir" && git -C "$dir" fetch -q --depth 1 "$url" "$rev" \
          && git -C "$dir" checkout -q FETCH_HEAD; \
      fi; \
    done

RUN go mod download

RUN go build -gcflags=all='-N -l' -o /src/build/bin/ ./cmd/...

RUN go test -count=1 -run '^$' ./... > /dev/null

RUN (go test -short -count=1 -timeout 20m ./consensus/... ./params/... ./core/types/... ./eth/protocols/... \
      || echo "WARNING: go test reported failures")
