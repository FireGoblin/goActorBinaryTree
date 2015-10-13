package ActorBinaryTree

import "fmt"

//dummy operation
type gc struct{}

func (g gc) ID() int {
	return 0
}

func (g gc) Elem() int {
	return 0
}

func (g gc) RequesterChan() chan OperationReply {
	return nil
}

func (g gc) Perform(node *binaryTreeNode) {
	return
}

func (g gc) String() string {
	return fmt.Sprintf("gc")
}
