BINARY = beefs-es-gen

SOURCE_DIR = .

COMMIT ?= $(shell git rev-parse --short HEAD)
LDFLAGS ?= -s -w

.PHONY : clean deps build linux release windows_build darwin_build linux_build bsd_build clean

clean:
	go clean -i $(GO_FLAGS) $(SOURCE_DIR)
	rm -f $(BINARY)
	rm -rf build/
	rm -rf linux/

build:
	go build $(GO_FLAGS) -ldflags "$(LDFLAGS)" -o $(BINARY) $(SOURCE_DIR)
	upx $(BINARY)
	cp $(BINARY) ~/go/bin

install:
	go install $(GO_FLAGS) -ldflags "$(LDFLAGS)" $(SOURCE_DIR)

linux:
	mkdir -p linux
	GOOS=linux GOARCH=amd64 go build $(GO_FLAGS) -ldflags "$(LDFLAGS)" -o linux/$(BINARY) $(SOURCE_DIR)
	upx linux/$(BINARY)