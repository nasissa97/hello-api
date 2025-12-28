GO_VERSION := 1.25.4


.PHONY: intall-go init-go

setup: install-go init-go install-lint copy-hooks install-godog

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

install-godog:
	go install github.com/cucumber/godog/cmd/godog@latest

copy-hooks:
	@echo "Installing git hooks..."
	@mkdir -p .git/hooks
	@for hook in $(shell ls scripts/hooks); do \
		cp scripts/hooks/$$hook .git/hooks/$${hook%.*}; \
		chmod +x .git/hooks/$${hook%.*}; \
		echo "Installed $${hook%.*}"; \
	done

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

install-helm:
	curl -fsSL -o get_helm.sh https://raw.githubusercontent.com/helm/helm/main/scripts/get-helm-4
	chmod 700 get_helm.sh
	./get_helm.sh

install-redis:
	helm repo add bitnami https://charts.bitnami.com/bitnami
	helm repo update
	kubectl create secret generic redis-cluster \
		--from-literal=redis-password=$$(openssl rand -base64 12) \
		--dry-run=client -o yaml | kubectl apply -f -
	
	helm upgrade --install redis-cluster bitnami/redis \
		--set auth.existingSecret=redis-cluster \
		--set auth.existingSecretPasswordKey=redis-password \
		--set architecture=replication \
		--set master.persistence.enabled=false \
		--set replica.persistence.enabled=false

build-image-prod:
	docker build -t hello-api:min .

build-image-dev:
	docker build -t hello-api:dev --target dev .

load-image-in-cluster:
	kind load docker-image hello-api:min
