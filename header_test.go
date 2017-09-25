package lapi

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Header", func() {
	It("NewHeader should return an instance of Header", func() {
		Expect(NewHeader()).NotTo(BeNil())
	})
})

var _ = Describe("FactoryHeader", func() {
	It("Get should return a value", func() {
		h := &FactoryHeader{make(map[string]string)}
		h.items["Content-Type"] = "application/json"
		values, ok := h.Get("content-Type")
		Expect(ok).To(BeTrue())
		Expect(values).To(Equal("application/json"))
	})

	It("Has should return a boolean", func() {
		h := &FactoryHeader{make(map[string]string)}
		h.items["Content-Type"] = "application/json"
		Expect(h.Has("content-Type")).To(BeTrue())
		Expect(h.Has("content-type")).To(BeTrue())
		Expect(h.Has("ContentType")).To(BeFalse())
	})

	It("Set should set key-value", func() {
		h := &FactoryHeader{make(map[string]string)}
		h.Set("content-type", "application/json")
		Expect(h.items["Content-Type"]).To(Equal("application/json"))
	})

	It("Remove should delete a key", func() {
		h := &FactoryHeader{make(map[string]string)}
		h.items["Content-Type"] = "application/json"
		h.Remove("content-TYPE")
		Expect(len(h.items)).To(Equal(0))
	})

	It("All should return all items", func() {
		h := &FactoryHeader{make(map[string]string)}
		h.items["Content-Type"] = "application/json"
		h.items["Content-Length"] = "1234"
		Expect(len(h.items)).To(Equal(2))
	})
})
