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

	It("Parallel should run functions by order", func() {
		m := make(map[int]*Slice)
		m[2] = new(Slice)
		m[0] = new(Slice)
		m[5] = new(Slice)
		m[2].Append("b").Append("a")
		m[0].Append("d").Append("e")
		m[5].Append("h").Append("c").Append("k")
		s := new(Slice)
		Parallel(m, func(item interface{}) {
			s.Append(item.(string))
		})
		Expect(len(s.All())).To(Equal(7))
	})
})
