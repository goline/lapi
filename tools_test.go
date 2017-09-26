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
		m := make(map[int][]interface{})
		m[2] = []interface{}{"a", "b"}
		m[0] = []interface{}{"d", "e"}
		m[5] = []interface{}{"g"}
		s := make([]string, 0)
		Parallel(m, func(item interface{}) {
			i := item.(string)
			s = append(s, i)
		})
		Expect(len(s)).To(Equal(5))
	})
})
