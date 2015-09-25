package main

type Remove struct {
	id            int
	elem          int
	requesterChan chan OperationReply
}

func (r *Remove) Id() int {
	return r.id
}

func (r *Remove) Elem() int {
	return r.elem
}

func (r *Remove) RequesterChan() chan OperationReply {
	return r.requesterChan
}

func (r *Remove) Perform(node *BinaryTreeNode) {
	if r.elem == node.elem {
		node.removed = true
	} else if r.elem < node.elem {
		if node.left != nil {
			node.leftChan() <- r
		}
	} else if r.elem > node.elem {
		if node.right != nil {
			node.rightChan() <- r
		}
	}
}
