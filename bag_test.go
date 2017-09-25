package lapi

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Bag", func() {
	It("NewBag should return an instance of Bag", func() {
		Expect(NewBag()).NotTo(BeNil())
	})
})

var _ = Describe("FactoryBag", func() {
	It("Get should return a value", func() {
		b := &FactoryBag{make(map[string]interface{})}
		b.items["my_key"] = "my_value"
		v, ok := b.Get("my_key")
		Expect(v).To(Equal("my_value"))
		Expect(ok).To(BeTrue())
	})

	It("Has should return a boolean", func() {
		b := &FactoryBag{make(map[string]interface{})}
		b.items["my_key"] = "my_value"
		Expect(b.Has("my_key")).To(BeTrue())
		Expect(b.Has("my_another_key")).To(BeFalse())
	})

	It("Set should allow to set value", func() {
		b := &FactoryBag{make(map[string]interface{})}
		Expect(b.Has("my_key")).To(BeFalse())
		b.Set("my_key", "my_value")
		Expect(b.Has("my_key")).To(BeTrue())
	})

	It("Remove should allow to remove a key", func() {
		b := &FactoryBag{make(map[string]interface{})}
		b.items["my_key"] = "my_value"
		Expect(b.Has("my_key")).To(BeTrue())
		b.Remove("my_key")
		Expect(b.Has("my_key")).To(BeFalse())
	})

	It("All should return all items", func() {
		b := &FactoryBag{make(map[string]interface{})}
		b.items["my_key"] = "my_value"
		b.items["my_another_key"] = 1
		Expect(len(b.All())).To(Equal(2))
	})
})
