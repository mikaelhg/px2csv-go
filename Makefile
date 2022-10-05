.PHONY: clean test build

build:
	GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" \
		-o ./bin/pcaxis2parquet-linux-amd64 ./cmd/pcaxis2parquet/
	GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" \
		-o ./bin/pcaxis2parquet-darwin-amd64 ./cmd/pcaxis2parquet/
	GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" \
		-o ./bin/pcaxis2parquet-windows-amd64.exe ./cmd/pcaxis2parquet/

clean:
	go clean
	rm ./bin/pcaxis2parquet-linux-amd64
	rm ./bin/pcaxis2parquet-darwin-amd64
	rm ./bin/pcaxis2parquet-windows-amd64.exe

test:
	time -v go run ./cmd/pcaxis2parquet/main.go \
		../gpcaxis/data/statfin_altp_pxt_12bd.px
