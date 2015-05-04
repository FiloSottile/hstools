GO     ?= go
GO     := env GOPATH="$(CURDIR):$(GOPATH)" $(GO)

.PHONY: montecarlo preprocess test clean

montecarlo:
	$(GO) install montecarlo

preprocess:
	$(GO) install preprocess

test:
	$(GO) test hstools -race -v -short

clean:
	- rm -r bin pkg
