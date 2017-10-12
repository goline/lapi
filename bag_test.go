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

	It("GetInt should return int64 value", func() {
		b := &FactoryBag{make(map[string]interface{})}
		b.items["my_int64"] = 10
		i, ok := b.GetInt("my_int64")
		Expect(ok).To(BeTrue())
		Expect(i).To(Equal(int64(10)))

		ii, ok := b.GetInt("my_another_int64")
		Expect(ok).To(BeFalse())
		Expect(ii).To(Equal(int64(0)))

		b.items["my_another_int64"] = "12"
		ii, ok = b.GetInt("my_another_int64")
		Expect(ok).To(BeTrue())
		Expect(ii).To(Equal(int64(12)))

		b.items["my_another_int64"] = new(Bag)
		ii, ok = b.GetInt("my_another_int64")
		Expect(ok).To(BeFalse())
		Expect(ii).To(Equal(int64(0)))
	})

	It("GetFloat should return float64 value", func() {
		b := &FactoryBag{make(map[string]interface{})}
		b.items["my_float64"] = 10.01
		f, ok := b.GetFloat("my_float64")
		Expect(ok).To(BeTrue())
		Expect(f).To(Equal(float64(10.01)))

		ff, ok := b.GetFloat("my_another_float64")
		Expect(ok).To(BeFalse())
		Expect(ff).To(Equal(float64(0.0)))

		b.items["my_another_float64"] = "12.2"
		ff, ok = b.GetFloat("my_another_float64")
		Expect(ok).To(BeTrue())
		Expect(ff).To(Equal(float64(12.2)))

		b.items["my_another_float64"] = float32(12.2)
		ff, ok = b.GetFloat("my_another_float64")
		Expect(ok).To(BeTrue())
		Expect(ff).To(Equal(float64(float32(12.2))))

		b.items["my_another_float64"] = new(Bag)
		ff, ok = b.GetFloat("my_another_float64")
		Expect(ok).To(BeFalse())
		Expect(ff).To(Equal(float64(0.0)))
	})

	It("GetString should return string value", func() {
		b := &FactoryBag{make(map[string]interface{})}
		b.items["my_string"] = "10.01"
		s, ok := b.GetString("my_string")
		Expect(ok).To(BeTrue())
		Expect(s).To(Equal("10.01"))

		ss, ok := b.GetString("my_another_string")
		Expect(ok).To(BeFalse())
		Expect(ss).To(BeEmpty())

		b.items["my_another_string"] = 12
		ss, ok = b.GetString("my_another_string")
		Expect(ok).To(BeFalse())
		Expect(ss).To(BeEmpty())
	})

	It("GetBool should return boolean value", func() {
		b := &FactoryBag{make(map[string]interface{})}
		b.items["my_bool"] = true
		v, ok := b.GetBool("my_bool")
		Expect(ok).To(BeTrue())
		Expect(v).To(BeTrue())

		vv, ok := b.GetBool("my_another_bool")
		Expect(ok).To(BeFalse())
		Expect(vv).To(BeFalse())

		b.items["my_another_bool"] = "true"
		vv, ok = b.GetBool("my_another_bool")
		Expect(ok).To(BeTrue())
		Expect(vv).To(BeTrue())

		b.items["my_another_bool"] = "1"
		vv, ok = b.GetBool("my_another_bool")
		Expect(ok).To(BeTrue())
		Expect(vv).To(BeTrue())

		b.items["my_bool"] = "false"
		vv, ok = b.GetBool("my_bool")
		Expect(ok).To(BeTrue())
		Expect(vv).To(BeFalse())

		b.items["my_bool"] = "0"
		vv, ok = b.GetBool("my_bool")
		Expect(ok).To(BeTrue())
		Expect(vv).To(BeFalse())
	})
})
