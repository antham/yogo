compile:
	mkdir -p build
	go build -ldflags "-s" -o build/yogo main.go
format:
	find -name "*.go" -exec go fmt {} \;
