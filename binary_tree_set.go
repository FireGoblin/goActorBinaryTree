package main

//max two children
type BinaryTreeSet struct {
	parentChan chan Operation
	childReply chan OperationReply

	child *BinaryTreeNode

	gcChan chan bool
}

func (b *BinaryTreeSet) childChan() chan Operation {
	if b.child == nil {
		return nil
	}

	return b.child.parentChan
}

func makeBinaryTreeSet() *BinaryTreeSet {
	x := BinaryTreeSet{make(chan Operation, 256), make(chan OperationReply, 16), makeBinaryTreeNode(0, true), make(chan bool, 1)}
	return &x
}
