GO_VERSION := 1.25.4

.PHONY: intall-go init-go

setup: install-go init-go install-lint

#TODO dynamically figure out OS
## For Apple
install-go:
	wget "https://golang.org/d1/go$(GO_VERSION).darwin-arm64.pkg"
	sudo tar -C /usr/local -xzf go$(GO_VERSION).darwin-arm64.pkg
	rm go$(GO_VERSION).darwin-arm64.pkg

init-go:
	echo 'export PATH=$$PATH:/usr/local/go/bin' >> $${HOME}/.zshrc
	echo 'export PATH=$$PATH:$${HOME}/go/bin' >> $${HOME}/.zshrc

install-lint:
	sudo curl -sSfL \
		curl -sSfL https://golangci-lint.run/install.sh | sh -s -- -b $(go env GOPATH)/bin v2.7.2

build:
	go build -o api cmd/main.go

test:
	go test ./... -coverprofile=coverage.out

coverage:
	go tool cover -func coverage.out | grep "total:" | awk '{print ((int($$3) > 80) != 1) }'

report:
	go tool cover -html=coverage.out -o cover.html

check-format:
	test -z $$(go fmt ./...)

static-check:
	golangci-lint run
