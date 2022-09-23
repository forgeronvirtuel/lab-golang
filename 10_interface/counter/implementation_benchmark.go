package counter

type Counter interface {
	Increment()
}

type ConcreteCounter struct {
	counter int
}

func (c *ConcreteCounter) Increment() {
	c.counter++
}
