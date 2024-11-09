# Makefile

.PHONY: test

test:
	@echo "Running tests with coverage..."
	@go test -coverprofile cover.out.tmp ./...
	@cat cover.out.tmp | grep -v '^github.com/moonicy/gometrics/cmd' | grep -v '^github.com/moonicy/gometrics/proto' > cover.out
	@go tool cover -func=cover.out
