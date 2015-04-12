GO     ?= go
GO     := env GOPATH="$(CURDIR):$(GOPATH)" $(GO)

.PHONY: montecarlo test clean

montecarlo:
	$(GO) install montecarlo

test:
	$(GO) test hstools -race -v -short

clean:
	- rm -r bin pkg
