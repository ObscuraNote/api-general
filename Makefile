ENV_LOCAL_TEST=\
  POSTGRES_HOST=localhost
	
TEST_DIRS := $(shell go list ./... | grep -E -v '/mocks|/cmd')

mock:
	rm -rf mocks
	go install go.uber.org/mock/mockgen@latest
	go generate -tags=tool mockgen.go

test: mock
	go fmt ./...
	@go test -count 1 -vet=all $(TEST_DIRS)

test-integration: mock
	@$(ENV_LOCAL_TEST) go test -count=1 -p=1 -vet=all -tags=integration $(TEST_DIRS)

verbose-test-integration: mock
	@$(ENV_LOCAL_TEST) go test -count=1 -p=1 -vet=all -tags=integration -v $(TEST_DIRS)