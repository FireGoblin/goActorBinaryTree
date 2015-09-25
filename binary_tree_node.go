package main

//max two children
type BinaryTreeNode struct {
	parentChan chan Operation
	childChan  chan OperationReply

	left  *BinaryTreeNode //use left.parentChan for sending operations to it
	right *BinaryTreeNode //use right.parentChan for sending operations to it

	elem    int
	removed bool
}

func (b *BinaryTreeNode) leftChan() chan Operation {
	if b.left == nil {
		return nil
	}

	return b.left.parentChan
}

func (b *BinaryTreeNode) rightChan() chan Operation {
	if b.right == nil {
		return nil
	}

	return b.right.parentChan
}

func makeBinaryTreeNode(element int) *BinaryTreeNode {
	//TODO: Tweak buffer sizes
	x := BinaryTreeNode{make(chan Operation, 64), make(chan OperationReply, 8), nil, nil, element, false}
	return &x
}
