package ActorBinaryTree

import "fmt"

//dummy operation
type GC struct{}

func (g GC) Id() int {
	return 0
}

func (g GC) Elem() int {
	return 0
}

func (g GC) RequesterChan() chan OperationReply {
	return nil
}

func (g GC) Perform(node *BinaryTreeNode) {
	return
}

func (g GC) String() string {
	return fmt.Sprintf("GC")
}
