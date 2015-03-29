GO     ?= go
GOPATH := $(CURDIR)
GO     := env GOPATH="$(GOPATH)" $(GO)

montecarlo:
	$(GO) install montecarlo

clean:
	- rm -r bin pkg
