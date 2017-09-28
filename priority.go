package lapi

type Prioritizer interface {
	Priority() int
}

type PriorityAware struct {
	priority int
}

func (p *PriorityAware) WithPriority(priority int) *PriorityAware {
	p.priority = priority
	return p
}

func (p *PriorityAware) Priority() int {
	return p.priority
}
