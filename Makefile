.PHONY: clean test build

build:
	@mkdir -p ./bin
	go build -o ./bin/pcaxis2parquet ./cmd/pcaxis2parquet/

cross: build
	GOOS=linux GOARCH=amd64 go build -buildmode=exe -ldflags="-s -w" \
		-o ./bin/pcaxis2parquet-linux-amd64 ./cmd/pcaxis2parquet/
	GOOS=linux GOARCH=arm64 go build -ldflags="-s -w" \
		-o ./bin/pcaxis2parquet-linux-arm64 ./cmd/pcaxis2parquet/
	GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" \
		-o ./bin/pcaxis2parquet-darwin-amd64 ./cmd/pcaxis2parquet/
	GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w" \
		-o ./bin/pcaxis2parquet-darwin-arm64 ./cmd/pcaxis2parquet/
	GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" \
		-o ./bin/pcaxis2parquet-windows-amd64.exe ./cmd/pcaxis2parquet/
	GOOS=windows GOARCH=arm64 go build -ldflags="-s -w" \
		-o ./bin/pcaxis2parquet-windows-arm64.exe ./cmd/pcaxis2parquet/
	# GOOS=js GOARCH=wasm go build -ldflags="-s -w" \
	#  	-o ./bin/pcaxis2parquet.wasm ./cmd/pcaxis2parquet/

clean:
	go clean
	@rm -f ./bin/pcaxis2parquet
	@rm -f ./bin/pcaxis2parquet-linux-amd64
	@rm -f ./bin/pcaxis2parquet-linux-arm64
	@rm -f ./bin/pcaxis2parquet-darwin-amd64
	@rm -f ./bin/pcaxis2parquet-darwin-arm64
	@rm -f ./bin/pcaxis2parquet-windows-amd64.exe
	@rm -f ./bin/pcaxis2parquet-windows-arm64.exe
	@rm -f ./bin/pcaxis2parquet.wasm

test:
	zcat ./data/statfin_vtp_pxt_124l.px.gz | time -v \
	    ./bin/pcaxis2parquet --px /dev/stdin --csv /dev/null

test-interpret:
	zcat ./data/statfin_vtp_pxt_124l.px.gz | time -v \
	    go run ./cmd/pcaxis2parquet/main.go --px /dev/stdin --csv /dev/null

test-debug:
	GODEBUG=gctrace=1 zcat ./data/statfin_vtp_pxt_124l.px.gz | time -v \
	    ./bin/pcaxis2parquet --px /dev/stdin --csv /dev/null
