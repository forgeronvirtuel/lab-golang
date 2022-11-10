package main

import (
	"fmt"
	"time"
)

var deposits = make(chan int)
var balances = make(chan int)

func Deposit(amount int) {
	deposits <- amount
}

func Balance() int {
	return <-balances
}

func teller() {
	var balance int
	for {
		select {
		case amount := <-deposits:
			balance += amount
		case balances <- balance:
		}
	}
}

func init() {
	go teller()
}

func main() {
	go func() {
		fmt.Println("Alice", Balance())
		Deposit(200)
		fmt.Println("Alice", Balance())
	}()

	go func() {
		fmt.Println("Bob", Balance())
		Deposit(-200)
		fmt.Println("Bob", Balance())
	}()

	time.Sleep(1 * time.Second)
}
