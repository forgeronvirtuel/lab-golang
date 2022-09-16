package main

import (
	"fmt"
	"net/url"
)

func main() {
	values := url.Values{}
	values.Add("item1", "value1")
	values.Add("item2", "value2")
	values.Add("item3", "value3")

	fmt.Printf("item1: `%s`\n", values.Get("item1"))
	fmt.Printf("item2: `%s`\n", values.Get("item2"))
	fmt.Printf("item3: `%s`\n", values.Get("item3"))

	values = nil

	fmt.Printf("item1: `%s`\n", values.Get("item1"))
	fmt.Printf("item2: `%s`\n", values.Get("item2"))
	fmt.Printf("item3: `%s`\n", values.Get("item3"))
}
