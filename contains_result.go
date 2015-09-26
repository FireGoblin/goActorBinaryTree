package main

import "fmt"

type ContainsResult struct {
	id     int
	result bool
}

func (c ContainsResult) Id() int {
	return c.id
}

func (c ContainsResult) Result() bool {
	return c.result
}

func (c ContainsResult) String() string {
	return fmt.Sprintf("ContainsResult(id: %d, result: %t)", c.id, c.result)
}
