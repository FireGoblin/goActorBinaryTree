package main

import "fmt"

type Remove struct {
	id            int
	elem          int
	requesterChan chan OperationReply
}

func (r Remove) Id() int {
	return r.id
}

func (r Remove) Elem() int {
	return r.elem
}

func (r Remove) RequesterChan() chan OperationReply {
	return r.requesterChan
}

func (r Remove) Perform(node *BinaryTreeNode) {
	if r.elem == node.elem {
		node.removed = true
		r.requesterChan <- OperationFinished{r.id}
	} else if r.elem < node.elem {
		if node.left != nil {
			node.leftChan() <- r
		} else {
			r.requesterChan <- OperationFinished{r.id}
		}
	} else if r.elem > node.elem {
		if node.right != nil {
			node.rightChan() <- r
		} else {
			r.requesterChan <- OperationFinished{r.id}
		}
	}
}

func (i Remove) String() string {
	return fmt.Sprintf("Remove(id: %d, elem: %d)", i.id, i.elem)
}
