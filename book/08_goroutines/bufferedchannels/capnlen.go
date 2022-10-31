package main

const nbElement = 100

func main() {
	ch := make(chan int, nbElement*2)
	go func() {
		for i := 0; i < nbElement; i++ {
			// (1) Here is the problem, let say that
			// this goroutine reads the len of ch,
			// and it indicates 0, so this condition
			// pass.
			if len(ch)%2 == 0 {
				// (2) But just after this goroutine sleep.

				// (4) Now this goroutine wakes up, but
				// the condition is no more valid.
				ch <- i
			} else {
				i--
			}
		}
	}()

	go func() {
		for i := 0; i < nbElement; i++ {
			// (3) Now this goroutine wakes up. It will add
			// one element to the channel and return to sleep
			// while trying to add the next one.
			ch <- i
		}
	}()

	// (5) The channel is read by the main goroutine
	for i := 0; i < nbElement*2; i++ {
		<-ch
	}
}
