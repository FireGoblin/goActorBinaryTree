package ActorBinaryTree

import "fmt"

type operationFinished struct {
	id int
}

func (o operationFinished) ID() int {
	return o.id
}

func (o operationFinished) String() string {
	return fmt.Sprintf("operationFinished(id: %d)", o.id)
}
