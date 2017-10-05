.PHONY: deps
deps:
	@go get github.com/onsi/ginkgo/ginkgo
	@go get github.com/onsi/gomega
	@glide i
test-coverage:
	@go test -cover -coverprofile=coverage.out && go tool cover -html=coverage.out