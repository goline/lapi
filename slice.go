package lapi

import (
	"sync"
)

type Slice struct {
	sync.RWMutex
	items []interface{}
}

type SliceItem struct {
	Index int
	Value interface{}
}

type SliceFunc func(item interface{})

func (s *Slice) Iter() <-chan SliceItem {
	c := make(chan SliceItem)

	f := func() {
		s.Lock()
		defer s.Lock()
		for index, value := range s.items {
			c <- SliceItem{index, value}
		}
		close(c)
	}
	go f()

	return c
}

func (s *Slice) All() []interface{} {
	return s.items
}

func (s *Slice) Append(item interface{}) *Slice {
	s.Lock()
	defer s.Unlock()

	s.items = append(s.items, item)
	return s
}

func (s *Slice) Remove(index int) *Slice {
	s.Lock()
	defer s.Unlock()

	s.items = append(s.items[:index], s.items[index+1:]...)
	return s
}

func (s *Slice) Get(index int) interface{} {
	s.Lock()
	defer s.Unlock()

	return s.items[index]
}

func (s *Slice) Run(f SliceFunc) {
	var wg sync.WaitGroup
	for i := range s.items {
		wg.Add(1)
		item := s.items[i]
		go func() {
			defer wg.Done()
			f(item)
		}()
	}
	wg.Wait()
}
