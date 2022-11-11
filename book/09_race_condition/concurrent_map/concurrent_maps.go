package main

import (
	"fmt"
	"time"
)

type Blob []byte

type Item struct {
	key   string
	value Blob
}

func NewItem(key string, data Blob) *Item {
	return &Item{
		key:   key,
		value: data,
	}
}

type itemChan struct {
	key string
	ch  chan Blob
}

func newItemChan(key string) *itemChan {
	return &itemChan{
		key: key,
		ch:  make(chan Blob),
	}
}

var (
	cache  = make(map[string]Blob)
	setter = make(chan *Item)
	getter = make(chan *itemChan)
)

func Set(item *Item) {
	setter <- item
}

func Get(key string) Blob {
	item := newItemChan(key)
	getter <- item
	return <-item.ch
}

func handleCache() {
	for {
		select {
		case itm := <-setter:
			cache[itm.key] = itm.value
		case itmChan := <-getter:
			itmChan.ch <- cache[itmChan.key]
		}
	}
}

func init() {
	go handleCache()
}

func main() {
	go func() {
		Set(NewItem("value1", []byte("This is sparta 1")))
		Set(NewItem("value2", []byte("This is sparta 2")))
		Set(NewItem("value3", []byte("This is sparta 3")))
	}()

	go func() {
		time.Sleep(1 * time.Second)
		fmt.Println(string(Get("value1")))
		fmt.Println(string(Get("value2")))
		fmt.Println(string(Get("value3")))
		fmt.Println(string(Get("value4")))
	}()

	time.Sleep(2 * time.Second)
}
