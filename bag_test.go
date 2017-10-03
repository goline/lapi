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

	It("GetInt64 should return int64 value", func() {
		b := &FactoryBag{make(map[string]interface{})}
		b.items["my_int64"] = 10
		i, ok := b.GetInt64("my_int64")
		Expect(ok).To(BeTrue())
		Expect(i).To(Equal(int64(10)))
	})

	It("GetFloat64 should return float64 value", func() {
		b := &FactoryBag{make(map[string]interface{})}
		b.items["my_float64"] = 10.01
		f, ok := b.GetFloat64("my_float64")
		Expect(ok).To(BeTrue())
		Expect(f).To(Equal(float64(10.01)))
	})

	It("GetString should return string value", func() {
		b := &FactoryBag{make(map[string]interface{})}
		b.items["my_string"] = "10.01"
		s, ok := b.GetString("my_string")
		Expect(ok).To(BeTrue())
		Expect(s).To(Equal("10.01"))
	})
})
