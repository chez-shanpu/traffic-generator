SHELL:=/bin/bash

GO=go

GO_VET_OPTS=-v
GO_TEST_OPTS=-v -race


.PHONY: build
build: $(CMDS)
	$(GO) build $(GO_BUILD_OPT) -o ./bin/tg

.PHONY: vet
vet:
	$(GO) vet $(GO_VET_OPTS) ./...

.PHONY: test
test: vet
	$(GO) test $(GO_TEST_OPTS) ./...

.PHONY: clean
clean:
	-$(GO) clean
	-rm $(RM_OPTS) bin/*

.PHONY: all
all: test build

.DEFAULT_GOAL=all
