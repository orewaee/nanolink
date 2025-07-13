.PHONY: build
build: clean
	go build -o build/nanolink cmd/cli/main.go
clean:
	rm -rf build/nanolink
