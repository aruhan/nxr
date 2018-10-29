GOCMD=go

build:
	$(GOCMD) build -v -o nxr.exe ./cmd/nxr

run:
	$(GOCMD) run -v ./cmd/nxr

.PHONY: clean build run