GO_VERSION := 1.25.4


.PHONY: intall-go init-go

setup: install-go init-go install-lint copy-hooks

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
	sudo curl -sSfL https://golangci-lint.run/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.64.5

copy-hooks:
	chmod +x ./scripts/hooks/pre-commit.sh
	cp -r ./scripts/hooks/ .git/pre-commit

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

