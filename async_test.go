package lapi

import (
	"github.com/goline/errors"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"fmt"
)

var _ = Describe("Async", func() {
	It("Add allows to add function(s)", func() {
		var async Async
		async.Add(func() errors.Error {
			return errors.New("some_code", "some_message")
		})

		for i := 0; i < 3; i++ {
			async.Add(func(index int) errors.Error {
				fmt.Println("f:", index)
				return nil
			}, i)
		}

		Expect(len(async.runners)).To(Equal(4))
	})

	It("Reset allows to clear all runners", func() {
		var async Async
		for i := 0; i < 3; i++ {
			async.Add(func(index int) errors.Error {
				fmt.Println("f:", index)
				return nil
			}, i)
		}
		Expect(len(async.runners)).To(Equal(3))

		async.Reset()
		Expect(len(async.runners)).To(Equal(0))
	})

	It("Run should run all runners and stop on error", func() {
		var async Async
		for i := 1; i < 5; i++ {
			async.Add(func(index int) errors.Error {
				if index%2 == 0 {
					return errors.New("STOP", fmt.Sprintf("%d", index))
				}

				return nil
			}, i)
		}

		err := async.Run()
		Expect(err).NotTo(BeNil())
		Expect(err.Message()).To(Or(Equal("2"), Equal("4")))
	})

	It("Run should run all runners and automatically quit if there is no errors", func() {
		var async Async
		for i := 1; i < 5; i += 2 {
			async.Add(func(index int) errors.Error {
				if index%2 == 0 {
					return errors.New("STOP", fmt.Sprintf("%d", index))
				}

				return nil
			}, i)
		}

		err := async.Run()
		Expect(err).To(BeNil())
	})

	It("Run should run all runners unless there is an error", func() {
		var async Async
		for i := 1; i < 10; i++ {
			async.Add(func(index int) errors.Error {
				if index == 3 {
					return errors.New("STOP", fmt.Sprintf("%d", index))
				}
				return nil
			}, i)
		}

		err := async.Run()
		Expect(err).NotTo(BeNil())
	})
})
