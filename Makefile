.PHONY: clean build test

BINARY=px2csv
PLATFORMS=linux darwin windows
ARCHITECTURES=amd64 arm64
LDFLAGS=-ldflags="-s -w"

build:
	@mkdir -p ./bin
	$(foreach GOOS, $(PLATFORMS),\
		$(foreach GOARCH, $(ARCHITECTURES), $(shell export GOOS=$(GOOS); export GOARCH=$(GOARCH); \
		    go build -buildmode=exe -o ./bin/$(BINARY)-$(GOOS)-$(GOARCH) ./cmd/px2csv/ \
	)))

clean:
	go clean
	@rm -f ./bin/px2csv*

test:
	zcat ./data/statfin_vtp_pxt_124l.px.gz | time -v \
	    ./bin/px2csv-linux-amd64 --px /dev/stdin --csv /dev/null

test-interpret:
	zcat ./data/statfin_vtp_pxt_124l.px.gz | time -v \
	    go run ./cmd/px2csv/main.go --px /dev/stdin --csv /dev/null

test-gctrace:
	GODEBUG=gctrace=1 zcat ./data/statfin_vtp_pxt_124l.px.gz | time -v \
	    ./bin/px2csv-linux-amd64 --px /dev/stdin --csv /dev/null

test-validate:
	@zcat ./data/010_kats_tau_101.px.gz | \
	    ./bin/px2csv-linux-amd64 --px /dev/stdin --csv /dev/stdout | \
		    diff -q data/testout.csv /dev/stdin || { echo "VALIDATE: FAILED!"; exit 1; }
	@echo "VALIDATE: SUCCESS"
