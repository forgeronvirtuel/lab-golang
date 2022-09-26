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

func doStuffAndIncrement(cnt Counter) {
	// do very important stuff
	if cnt != nil {
		cnt.Increment()
	}
}
