package lapi

// Container acts as a dependency-injection manager
type Container interface {
	Binder
	Resolver
	Injector
}

// Binder uses to bind a concrete to an abstract
type Binder interface {
	// Bind stores a concrete of an abstract, as default sharing is enable
	Bind(abstract interface{}, concrete interface{}) error
}

// Resolver helps to resolve dependencies
type Resolver interface {
	// Resolve processes and returns a concrete of proposed abstract
	Resolve(abstract interface{}) (concrete interface{}, err error)
}

// Injector works as a tool to inject dependencies
type Injector interface {
	// Inject resolves target's dependencies
	Inject(target interface{}) error
}
