package main

type OperationFinished struct {
	id int
}

func (o OperationFinished) Id() int {
	return o.id
}
