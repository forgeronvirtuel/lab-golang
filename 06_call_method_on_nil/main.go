package main

import (
	"fmt"
	"net/url"
)

type MyArray []string

func (m MyArray) Get(idx int) string {
	if m != nil {
		return m[idx]
	}
	return ""
}

func (m MyArray) Add(idx int, value string) {
	m[idx] = value
}

type MyStruct struct {
	internalValue string
}

func (m *MyStruct) Get() string {
	if m == nil {
		return ""
	}
	return m.internalValue
}

func (m *MyStruct) Add(v string) {
	m.internalValue = v
}

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

	arr := make(MyArray, 10)
	arr.Add(0, "item1")
	arr.Add(1, "item2")
	arr.Add(2, "item3")
	fmt.Printf("item1: `%s`\n", arr.Get(0))
	fmt.Printf("item2: `%s`\n", arr.Get(1))
	fmt.Printf("item3: `%s`\n", arr.Get(2))

	arr = nil
	fmt.Printf("item1: `%s`\n", arr.Get(0))
	fmt.Printf("item2: `%s`\n", arr.Get(1))
	fmt.Printf("item3: `%s`\n", arr.Get(2))

	mystruct := &MyStruct{}
	mystruct.Add("value")
	fmt.Printf("item: `%s`\n", mystruct.Get())

	mystruct = nil
	fmt.Printf("item: `%s`\n", mystruct.Get())
}
