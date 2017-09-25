package lapi

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Connection", func() {
	It("NewConnection should return an instance of Connection", func() {
		Expect(NewConnection(nil, nil)).NotTo(BeNil())
	})
})

var _ = Describe("FactoryConnection", func() {
	It("Request should return instance of Request", func() {
		r := &FactoryRequest{}
		c := &FactoryConnection{}
		c.request = r
		Expect(c.Request()).NotTo(BeNil())
	})

	It("WithRequest should set instance of Request", func() {
		r := &FactoryRequest{}
		c := &FactoryConnection{}
		Expect(c.WithRequest(r).Request()).NotTo(BeNil())
	})

	It("Response should return instance of Response", func() {
		r := &FactoryResponse{}
		c := &FactoryConnection{}
		c.response = r
		Expect(c.Response()).NotTo(BeNil())
	})

	It("WithResponse should set instance of Response", func() {
		r := &FactoryResponse{}
		c := &FactoryConnection{}
		Expect(c.WithResponse(r).Response()).NotTo(BeNil())
	})
})
