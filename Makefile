GO=go
GOBIN=${GOPATH}/bin

.PHONY: lint
lint: go.lint.verify
	@echo "===========> Run golangci-lint to lint source codes"
	@$(GOBIN)/golangci-lint run --fix

.PHONY: test
test:
	@sh ${PWD}/script/gotest.sh

.PHONY: go.lint.verify
go.lint.verify:
ifeq (,$(wildcard $(GOBIN)/golangci-lint))
	@echo "===========> Installing golangci-lint"
	@GO111MODULE=on $(GO) get github.com/golangci/golangci-lint/cmd/golangci-lint
	@$(GO) install github.com/golangci/golangci-lint/cmd/golangci-lint
endif

