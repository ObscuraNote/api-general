TEST_DIRS := $(shell go list ./... | grep -E -v '/mocks|/cmd')

test:
	go fmt ./...
	@go test -count 1 -vet=all $(TEST_DIRS)

test-integration:
	go test -count=1 -p=1 -vet=all -tags=integration $(TEST_DIRS)

verbose-test-integration:
	go test -count=1 -p=1 -vet=all -tags=integration -v $(TEST_DIRS)