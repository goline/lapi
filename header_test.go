package lapi

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"sync"
)

var _ = Describe("Header", func() {
	It("NewHeader should return an instance of Header", func() {
		Expect(NewHeader()).NotTo(BeNil())
	})
})

var _ = Describe("FactoryHeader", func() {
	It("Get should return a value", func() {
		h := &FactoryHeader{new(sync.Map)}
		h.items.Store("Content-Type", "application/json")
		values, ok := h.Get("content-Type")
		Expect(ok).To(BeTrue())
		Expect(values).To(Equal("application/json"))
	})

	It("Has should return a boolean", func() {
		h := &FactoryHeader{new(sync.Map)}
		h.items.Store("Content-Type", "application/json")
		Expect(h.Has("content-Type")).To(BeTrue())
		Expect(h.Has("content-type")).To(BeTrue())
		Expect(h.Has("ContentType")).To(BeFalse())
	})

	It("Set should set key-value", func() {
		h := &FactoryHeader{new(sync.Map)}
		h.Set("content-type", "application/json")
		v, ok := h.items.Load("Content-Type")
		Expect(ok).To(BeTrue())
		Expect(v).To(Equal("application/json"))
	})

	It("Remove should delete a key", func() {
		h := &FactoryHeader{new(sync.Map)}
		h.items.Store("Content-Type", "application/json")
		h.Remove("content-TYPE")
		_, ok := h.items.Load("Content-Type")
		Expect(ok).To(BeFalse())
	})

	It("All should return all items", func() {
		h := &FactoryHeader{new(sync.Map)}
		h.items.Store("Content-Type", "application/json")
		h.items.Store("Content-Length", "1234")
		Expect(len(h.All())).To(Equal(2))
	})
})
