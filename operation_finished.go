package main

import "fmt"

type OperationFinished struct {
	id int
}

func (o OperationFinished) Id() int {
	return o.id
}

func (o OperationFinished) String() string {
	return fmt.Sprintf("OperationFinished(id: %d)", o.id)
}
