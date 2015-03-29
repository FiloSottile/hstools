GO     ?= go
GO     := env GOPATH="$(CURDIR):$(GOPATH)" $(GO)

montecarlo:
	$(GO) install montecarlo

clean:
	- rm -r bin pkg
