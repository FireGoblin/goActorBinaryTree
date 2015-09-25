package main

type Insert struct {
	id            int
	elem          int
	requesterChan chan OperationReply
}

func (i *Insert) Id() int {
	return i.id
}

func (i *Insert) Elem() int {
	return i.elem
}

func (i *Insert) RequesterChan() chan OperationReply {
	return i.requesterChan
}

func (i *Insert) Perform(node *BinaryTreeNode) {
	if i.elem < node.elem {
		if node.left != nil {
			node.leftChan() <- i
		} else {
			node.left = makeBinaryTreeNode(i.elem, false)
		}
	} else if i.elem > node.elem {
		if node.right != nil {
			node.rightChan() <- i
		} else {
			node.right = makeBinaryTreeNode(i.elem, false)
		}
	}
}
