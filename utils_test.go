package lapi

import (
	"errors"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Utils", func() {
	It("PanicOnError should panic if there is an error", func() {
		defer func() {
			if r := recover(); r == nil {
				Expect(r).NotTo(BeNil())
			}
		}()
		PanicOnError(errors.New("ERROR"))
	})

	It("Must verifies and panics if there is an error", func() {
		defer func() {
			if r := recover(); r == nil {
				Expect(r).NotTo(BeNil())
			}
		}()
		Must(nil, errors.New("ERROR"))
	})
})
