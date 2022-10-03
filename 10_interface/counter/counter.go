package counter

import "fmt"

type IncrementalCounter interface {
	Increment()
}

type DecrementalCounter interface {
	Decrement()
}

type Counter struct {
	counter int
}

func (c *Counter) Increment() {
	c.counter++
}

func (c *Counter) Decrement() {
	c.counter--
}

func (c *Counter) String() string {
	return fmt.Sprintf("%d", c.counter)
}

func doStuffAndIncrement(cnt IncrementalCounter) {
	// do very important stuff
	if cnt != nil {
		cnt.Increment()
	}
}
