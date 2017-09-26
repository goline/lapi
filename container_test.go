package lapi

import (
	"reflect"

	"github.com/goline/errors"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"sync"
)

var _ = Describe("Container", func() {
	It("NewContainer should return Container", func() {
		Expect(NewContainer()).NotTo(BeNil())
	})
})

type sampleInvalidError struct{}
type facErr struct {
	code string
	msg  string
}
type InjectErrorResolveErrorNotExistAbstract struct {
	Err errors.Error `inject:"*"`
}
type InjectOk struct {
	Err errors.Error `inject:"*"`
}
type InjectRecursiveOk struct {
	Err                          errors.Error `inject:"*"`
	Foo                          InjectFooer  `inject:"*"`
	NotInjectableNonInterface    string       `inject:"*"`
	notInjectablePrivateProperty errors.Error `inject:"*"`
}
type InjectFooer interface {
	Foo() string
}
type InjectFoo struct {
	Baz InjectBazer `inject:"*"`
}

func (f *InjectFoo) Foo() string { return "Foo.." }

type InjectBazer interface {
	Baz() string
}
type InjectBaz struct{}

func (b *InjectBaz) Baz() string { return "Baz.." }

var _ = Describe("FactoryContainer", func() {
	It("Bind should return error code ERR_BIND_INVALID_ARGUMENTS", func() {
		c := NewContainer()
		err := c.Bind("string", &errors.FactoryError{})
		Expect(err).NotTo(BeNil())
		Expect(err.(errors.Error).Code()).To(Equal(ERR_BIND_INVALID_ARGUMENTS))
	})

	It("Bind should return error code ERR_BIND_NOT_IMPLEMENT_INTERFACE", func() {
		c := NewContainer()
		err := c.Bind((*Bag)(nil), &errors.FactoryError{})
		Expect(err).NotTo(BeNil())
		Expect(err.(errors.Error).Code()).To(Equal(ERR_BIND_NOT_IMPLEMENT_INTERFACE))
	})

	It("Bind should return error code ERR_BIND_INVALID_CONCRETE", func() {
		c := NewContainer()
		a := make([]interface{}, 2)
		a[0] = "a_string"
		a[1] = sampleInvalidError{}
		for _, i := range a {
			err := c.Bind((*Bag)(nil), i)
			Expect(err).NotTo(BeNil())
			Expect(err.(errors.Error).Code()).To(Equal(ERR_BIND_INVALID_CONCRETE))
		}
	})

	It("Bind should return nil when binding struct", func() {
		c := NewContainer()
		err := c.Bind(FactoryBag{}, &FactoryBag{})
		Expect(err).To(BeNil())
	})

	It("Bind should return error code ERR_BIND_INVALID_STRUCT", func() {
		c := NewContainer()
		err := c.Bind(FactoryBag{}, "not_a_struct")
		Expect(err).NotTo(BeNil())
		Expect(err.(errors.Error).Code()).To(Equal(ERR_BIND_INVALID_STRUCT))
	})

	It("Bind should return error code ERR_BIND_INVALID_STRUCT_CONCRETE", func() {
		c := NewContainer()
		err := c.Bind(FactoryBag{}, &errors.FactoryError{})
		Expect(err).NotTo(BeNil())
		Expect(err.(errors.Error).Code()).To(Equal(ERR_BIND_INVALID_STRUCT_CONCRETE))
	})

	It("Bind should return nil when binding function", func() {
		c := NewContainer()
		err := c.Bind((*errors.Error)(nil), func(code string, message string, err error) errors.Error {
			return errors.New(code, message)
		})
		Expect(err).To(BeNil())
	})

	It("Bind should return nil when binding a pointer", func() {
		c := NewContainer()
		err := c.Bind((*errors.Error)(nil), errors.New("my_code", "my_message"))
		Expect(err).To(BeNil())
	})

	It("structOf should return nil", func() {
		c := &FactoryContainer{}
		_, err := c.structOf(reflect.TypeOf(FactoryBag{}))
		Expect(err).To(BeNil())
	})

	It("Resolve should return error code ERR_RESOLVE_NOT_EXIST_ABSTRACT", func() {
		c := NewContainer()
		_, err := c.Resolve(&FactoryBag{})
		Expect(err).NotTo(BeNil())
		Expect(err.(errors.Error).Code()).To(Equal(ERR_RESOLVE_NOT_EXIST_ABSTRACT))
	})

	It("Resolve should return error code ERR_RESOLVE_INSUFFICIENT_ARGUMENTS", func() {
		c := NewContainer()
		c.Bind((*errors.Error)(nil), func(code string, message string, err error) errors.Error {
			return errors.New(code, message)
		})
		_, err := c.Resolve((*errors.Error)(nil), "my_code", "my_message")
		Expect(err).NotTo(BeNil())
		Expect(err.(errors.Error).Code()).To(Equal(ERR_RESOLVE_INSUFFICIENT_ARGUMENTS))
	})

	It("Resolve should return error code ERR_RESOLVE_NON_VALUES_RETURNED", func() {
		c := NewContainer()
		c.Bind((*errors.Error)(nil), func(code string, message string) {})
		_, err := c.Resolve((*errors.Error)(nil), "my_code", "my_message")
		Expect(err).NotTo(BeNil())
		Expect(err.(errors.Error).Code()).To(Equal(ERR_RESOLVE_NON_VALUES_RETURNED))
	})

	It("Resolve should return nil", func() {
		c := NewContainer()
		c.Bind((*errors.Error)(nil), func(code string, message string) errors.Error {
			return errors.New(code, message)
		})
		e, err := c.Resolve((*errors.Error)(nil), "my_code", "my_message")
		Expect(err).To(BeNil())
		Expect(e.(errors.Error).Code()).To(Equal("my_code"))
	})

	It("Resolve should return nil when resolving pointer", func() {
		c := NewContainer()
		c.Bind((*errors.Error)(nil), errors.New("my_code", "my_message"))
		e, err := c.Resolve((*errors.Error)(nil))
		Expect(err).To(BeNil())
		Expect(e.(errors.Error).Code()).To(Equal("my_code"))
	})

	It("Resolve should return nil when resolving struct", func() {
		c := NewContainer()
		c.Bind(FactoryBag{}, NewBag())
		v, err := c.Resolve(FactoryBag{})
		Expect(err).To(BeNil())
		Expect(v).NotTo(BeNil())
	})

	It("Resolve should return error code ERR_RESOLVE_INVALID_ARGUMENTS", func() {
		c := NewContainer()
		_, err := c.Resolve("string")
		Expect(err).NotTo(BeNil())
		Expect(err.(errors.Error).Code()).To(Equal(ERR_RESOLVE_INVALID_ARGUMENTS))
	})

	It("Inject should return error code ERR_INJECT_INVALID_TARGET_TYPE", func() {
		c := NewContainer()
		e := facErr{"my_code", "my_message"}
		err := c.Inject(e)
		Expect(err).NotTo(BeNil())
		Expect(err.(errors.Error).Code()).To(Equal(ERR_INJECT_INVALID_TARGET_TYPE))
	})

	It("Inject should return error code ERR_RESOLVE_NOT_EXIST_ABSTRACT", func() {
		c := NewContainer()
		in := &InjectErrorResolveErrorNotExistAbstract{}
		err := c.Inject(in)
		Expect(err).NotTo(BeNil())
		Expect(err.(errors.Error).Code()).To(Equal(ERR_RESOLVE_NOT_EXIST_ABSTRACT))
	})

	It("Inject should return nil", func() {
		c := NewContainer()
		c.Bind((*errors.Error)(nil), errors.New("my_code", "my_message"))
		in := &InjectOk{}
		err := c.Inject(in)
		Expect(err).To(BeNil())
		Expect(in.Err.Code()).To(Equal("my_code"))
		Expect(in.Err.Message()).To(Equal("my_message"))
	})

	It("Inject should do a recursive injection", func() {
		// InjectRecursiveOk -> Err
		// InjectRecursiveOk -> Foo -> Baz
		c := NewContainer()
		c.Bind((*errors.Error)(nil), errors.New("my_code", "my_message"))
		c.Bind((*InjectFooer)(nil), &InjectFoo{})
		c.Bind((*InjectBazer)(nil), &InjectBaz{})
		in := &InjectRecursiveOk{}
		err := c.Inject(in)
		Expect(err).To(BeNil())

		// Asserting for in.Err
		Expect(in.Err.Code()).To(Equal("my_code"))
		Expect(in.Err.Message()).To(Equal("my_message"))

		// Asserting for in.Foo
		Expect(in.Foo.Foo()).To(Equal("Foo.."))
		Expect(in.Foo.(*InjectFoo).Baz.Baz()).To(Equal("Baz.."))
	})

	It("instanceOf should return false", func() {
		c := &FactoryContainer{new(sync.Map)}
		b := c.instanceOf(reflect.TypeOf("a_string"), reflect.TypeOf((*Bag)(nil)))
		Expect(b).To(BeFalse())
	})

	It("instanceOf should return false (ConcreteTypeNotSupport)", func() {
		c := &FactoryContainer{new(sync.Map)}
		i, _ := c.interfaceOf((*Bag)(nil))
		b := c.instanceOf(i, reflect.TypeOf("a_string"))
		Expect(b).To(BeFalse())
	})
})
