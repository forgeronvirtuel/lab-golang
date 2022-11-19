package main

import (
	"fmt"
	b2 "github.com/forgeronvirtuel/labgolang/book/09_race_condition/semabank/bank2"
	b3 "github.com/forgeronvirtuel/labgolang/book/09_race_condition/semabank/bank3"
)

func main() {
	b2.Deposit(100)
	fmt.Println(b2.Balance())

	b3.Deposit(100)
	fmt.Println(b3.Balance())
}
