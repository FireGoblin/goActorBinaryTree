package ActorBinaryTree

import "fmt"

type containsResult struct {
	id     int
	result bool
}

func (c containsResult) ID() int {
	return c.id
}

func (c containsResult) Result() bool {
	return c.result
}

func (c containsResult) String() string {
	return fmt.Sprintf("containsResult(id: %d, result: %t)", c.id, c.result)
}
