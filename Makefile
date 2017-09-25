.PHONY: deps
deps:
	@go get github.com/onsi/ginkgo/ginkgo
	@go get github.com/onsi/gomega
	@glide i