package main

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
