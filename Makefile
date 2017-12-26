.PHONY: build

BUILD_NUMBER ?= $(BUILD_NUMBER:)
BUILD_DATE = $(shell date -u)
BUILD_HASH = $(shell git rev-parse HEAD)
ifeq ($(BUILD_NUMBER),)
	BUILD_NUMBER := dev
endif

DIST_ROOT=dist
DIST_PATH=$(DIST_ROOT)/single-sign-on

GO=go
GO_LINKER_FLAGS ?= -ldflags \
				   "-X github.com/KenmyZhang/single-sign-on/model.BuildNumber=$(BUILD_NUMBER)\
				    -X 'github.com/KenmyZhang/single-sign-on/model.BuildDate=$(BUILD_DATE)'\
				    -X github.com/KenmyZhang/single-sign-on/model.BuildHash=$(BUILD_HASH)"

build-linux:
	@echo Build Linux amd64
	env GOOS=linux GOARCH=amd64 $(GO) install $(GOFLAGS) $(GO_LINKER_FLAGS) ./cmd/single-sign-on

build: build-linux

package: build
	@ echo Packaging single-sign-on

	rm -Rf $(DIST_ROOT)

	mkdir -p $(DIST_PATH)/bin
	mkdir -p $(DIST_PATH)/logs

	cp -RL config $(DIST_PATH)
	cp -RL i18n $(DIST_PATH)
	cp -RL templates $(DIST_PATH)
	cp -RL fonts $(DIST_PATH)	

	cp $(GOPATH)/bin/single-sign-on $(DIST_PATH)/bin

	tar -C dist -czf $(DIST_PATH)-linux-amd64.tar.gz single-sign-on
	rm -f $(DIST_PATH)/bin/single-sign-on

