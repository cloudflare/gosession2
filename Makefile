GOCMD := go
REVISION := $(shell git log -n1 --pretty=format:%h)
FULL_REVISION := $(shell git log -n1 --pretty=format:%H)
TAGS := dev
LDFLAGS := -X main.revision $(FULL_REVISION)

export GOPATH := $(PWD):$(PWD)/gopath
export PATH := $(PWD)/bin:$(PWD)/gopath/bin:$(PATH)

.PHONY: all
all:
	$(GOCMD) install -tags '$(TAGS)' -ldflags '$(LDFLAGS)' -v ./...

.PHONY: race
race:
	$(GOCMD) install -race -v `$(GOCMD) list -f '{{if eq .Name "main"}}{{.ImportPath}}{{end}}' ./...`

.PHONY: test-compile
test-compile: all
	@$(MAKE) --no-print-directory -f Make.tests $@

.PHONY: clean
clean:
	$(RM) -rf bin pkg

.PHONY: gofmt
gofmt:
	$(GOCMD) fmt ./...
