GOCMD := go
REVISION := $(shell git log -n1 --pretty=format:%h)
FULL_REVISION := $(shell git log -n1 --pretty=format:%H)
TAGS := dev
LDFLAGS := -X main.revision $(FULL_REVISION)

export GOPATH := $(PWD):$(PWD)/gopath
export PATH := $(PWD)/bin:$(PWD)/gopath/bin:$(PATH)

.PHONY: all
all:
	$(GOCMD) install -tags '$(TAGS)' -ldflags '$(LDFLAGS)' -v gophq.io...

.PHONY: race
race:
	$(GOCMD) install -race -v `$(GOCMD) list -f '{{if eq .Name "main"}}{{.ImportPath}}{{end}}' gophq.io...`

.PHONY: test
test:
	$(GOCMD) test -v gophq.io/...

.PHONY: test-race
test-race:
	$(GOCMD) test -race gophq.io/...

.PHONY: run-tls-server
run-tls-server: all
	./bin/gophqd -tls.ca=etc/ca.crt -tls.cert=etc/server.crt -tls.key=etc/server.key

.PHONY: run-tls-producer
run-tls-producer: all
	./bin/gophq -mode=produce -tls.ca=etc/ca.crt -tls.cert=etc/client.crt -tls.key=etc/client.key

TEST=$(subst $(space),$(newline),$(shell cd src && $(GOCMD) list -f '{{if or .TestGoFiles .XTestGoFiles}}{{.Dir}}{{end}}' ./gophq.io...))

.PHONY: test-compile
test-compile: $(addsuffix .test-compile, $(TEST))

%.test-compile: all
	cd $* && $(GOCMD) test -compiler=$(COMPILER) -p 1 -v -c .

.PHONY: clean
clean:
	$(RM) -rf bin pkg

.PHONY: gofmt
gofmt:
	$(GOCMD) fmt ./...
