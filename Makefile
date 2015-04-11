GO     ?= go
GO     := env GOPATH="$(CURDIR):$(GOPATH)" $(GO)

montecarlo:
	$(GO) install montecarlo

test:
	$(GO) test hspredict -race -v -short
	$(GO) test montecarlo -race -v -short

clean:
	- rm -r bin pkg
