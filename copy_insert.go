package main

import "fmt"

type CopyInsert struct {
	id            int
	elem          int
	requesterChan chan OperationReply
}

func (i CopyInsert) Id() int {
	return i.id
}

func (i CopyInsert) Elem() int {
	return i.elem
}

func (i CopyInsert) RequesterChan() chan OperationReply {
	return i.requesterChan
}

func (i CopyInsert) Perform(node *BinaryTreeNode) {
	if i.elem == node.elem {
		node.removed = false
		i.requesterChan <- OperationFinished{i.id}
	} else if i.elem < node.elem {
		if node.left != nil {
			node.leftChan() <- i
		} else {
			node.left = makeBinaryTreeNode(i.elem, false)
			node.left.parent = node.childReply
			i.requesterChan <- OperationFinished{i.id}
		}
	} else if i.elem > node.elem {
		if node.right != nil {
			node.rightChan() <- i
		} else {
			node.right = makeBinaryTreeNode(i.elem, false)
			node.right.parent = node.childReply
			i.requesterChan <- OperationFinished{i.id}
		}
	}
}

func (i CopyInsert) String() string {
	return fmt.Sprintf("Insert(id: %d, elem: %d)", i.id, i.elem)
}
