package lapi

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("App", func() {
	It("NewApp should return an instance of App", func() {
		Expect(NewApp()).NotTo(BeNil())
	})
})
