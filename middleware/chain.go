package middleware

import "net/http"

//Adapter represents middleware abstraction
type Adapter func(http.Handler) http.Handler

//Chain holds a slice of adapters
type Chain struct {
	adapters []Adapter
}

//New returns new adapter chain
func New(adapters ...Adapter) Chain {
	return Chain{append(([]Adapter)(nil), adapters...)}
}

//Append adapters to the current chain
func (c Chain) Append(adapters ...Adapter) Chain {
	newCons := make([]Adapter, 0, len(c.adapters)+len(adapters))
	newCons = append(newCons, c.adapters...)
	newCons = append(newCons, adapters...)

	return Chain{newCons}
}

//Extend the adapter chain
func (c Chain) Extend(chain Chain) Chain {
	return c.Append(chain.adapters...)
}

//Then chains middleware
func (c Chain) Then(h http.Handler) http.Handler {
	if h == nil {
		h = http.DefaultServeMux
	}

	for i := range c.adapters {
		h = c.adapters[len(c.adapters)-1-i](h)
	}

	return h
}
