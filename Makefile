.PHONY: test clean all
.PHONY: preprocess brute lookmeup announce grind curiosity scrolls montecarlo
all:    preprocess brute lookmeup announce grind curiosity scrolls

GO ?= go

montecarlo:
	GOPATH="$(CURDIR)" $(GO) build -o bin/montecarlo src/cmd/montecarlo.go

preprocess:
	GOPATH="$(CURDIR)" $(GO) build -o bin/preprocess src/cmd/preprocess.go

brute:
	GOPATH="$(CURDIR)" $(GO) build -o bin/brute src/cmd/brute.go

lookmeup:
	GOPATH="$(CURDIR)" $(GO) build -o bin/lookmeup src/cmd/lookmeup.go

announce:
	GOPATH="$(CURDIR)" $(GO) build -o bin/announce src/cmd/announce.go

grind:
	GOPATH="$(CURDIR)" $(GO) build -o bin/grind src/cmd/grind.go

curiosity:
	GOPATH="$(CURDIR)" $(GO) build -o bin/curiosity src/cmd/curiosity.go

scrolls:
	GOPATH="$(CURDIR)" $(GO) build -o bin/scrolls src/cmd/scrolls.go

test:
	GOPATH="$(CURDIR)" $(GO) test hstools -race -v -short

clean:
	- rm -r bin pkg
