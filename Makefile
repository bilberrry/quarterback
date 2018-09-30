
GO_CMD=go
GO_BUILD=$(GO_CMD) build
GO_CLEAN=$(GO_CMD) clean
GO_DEPS=$(GO_CMD) get -d -v
GO_VET=$(GO_CMD) vet
GO_FMT=$(GO_CMD)fmt -w
BIN_NAME := quarterback

all: deps fmt vet build

deps:
	@echo "==> Install dependencies for $(BIN_NAME) ..."; \
	$(GO_DEPS)

build:
	@echo "==> Build $(BIN_NAME) ..."; \
	$(GO_BUILD) -o $(BIN_NAME)

clean:
	@echo "==> Clean $(BIN_NAME) ..."; \
	$(GO_CLEAN)
	rm -f $(BIN_NAME)

fmt:
	@echo "==> Formatting $(BIN_NAME) ..."; \
	$(GO_FMT) -w *.go */*.go

vet:
	@echo "==> Vet $(BIN_NAME) ..."; \
	$(GO_VET)

release:
	@goreleaser --skip-validate