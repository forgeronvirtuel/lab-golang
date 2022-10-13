package main

func main() {
	ch := make(chan int, 10)
	go func() {
		for i := 0; i < 30; i++ {
			// (1) Here is the problem, let say that
			// this goroutine reads the len of ch
			// and it indicate 1, so this condition
			// pass.
			if len(ch) > 0 {
				// (2) But just after the call to len(ch)
				// this goroutine sleep.

				// (4) Now this goroutine wakes up, but
				// the condition is no more valid.
				ch <- i
			}
		}
	}()

	go func() {
		for i := 0; i < 200; i++ {
			// (3) Now this groutine wakes up. It will add
			// one element to the channel and return to sleep
			// while trying to add the next one.
			ch <- i
		}
	}()
}
