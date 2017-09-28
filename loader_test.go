package lapi

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ServerLoader", func() {
	It("ss", func() {
		l := &ServerLoader{}
		l.WithPriority(5)
		Expect(l.Priority()).To(Equal(5))
	})
})
