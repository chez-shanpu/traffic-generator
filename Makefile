SHELL:=/bin/bash

GO=go

GO_VET_OPTS=-v
GO_TEST_OPTS=-v -race

CMD_NAMES:=$(basename $(wildcard cmd/*.go))
CMDS:=$(subst cmd,bin,$(CMD_NAMES))


.SECONDEXPANSION:
#bin/%: $(wildcard cmd/*/*.go) $(wildcard cmd/*/*/*.go) $(wildcard pkg/*/*.go) go.mod bin
bin/%:
	$(GO) build $(GO_BUILD_OPT) -o $@


.PHONY: build
build: $(CMDS)

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