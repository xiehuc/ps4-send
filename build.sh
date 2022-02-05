CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -o ps4-send-macos-arm64 main.go
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ps4-send-linux-amd64 main.go
