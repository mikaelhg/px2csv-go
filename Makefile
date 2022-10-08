.PHONY: clean build test

BINARY=pcaxis2parquet
PLATFORMS=linux darwin windows
ARCHITECTURES=amd64 arm64
LDFLAGS=-ldflags="-s -w"

build:
	@mkdir -p ./bin
	$(foreach GOOS, $(PLATFORMS),\
		$(foreach GOARCH, $(ARCHITECTURES), $(shell export GOOS=$(GOOS); export GOARCH=$(GOARCH); \
		    go build -buildmode=exe -o ./bin/$(BINARY)-$(GOOS)-$(GOARCH) ./cmd/pcaxis2parquet/ \
	)))

clean:
	go clean
	@rm -f ./bin/pcaxis2parquet*

test:
	zcat ./data/statfin_vtp_pxt_124l.px.gz | time -v \
	    ./bin/pcaxis2parquet --px /dev/stdin --csv /dev/null

test-interpret:
	zcat ./data/statfin_vtp_pxt_124l.px.gz | time -v \
	    go run ./cmd/pcaxis2parquet/main.go --px /dev/stdin --csv /dev/null

test-debug:
	GODEBUG=gctrace=1 zcat ./data/statfin_vtp_pxt_124l.px.gz | time -v \
	    ./bin/pcaxis2parquet --px /dev/stdin --csv /dev/null
