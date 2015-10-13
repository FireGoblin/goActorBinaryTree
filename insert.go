package ActorBinaryTree

import "fmt"

type insert struct {
	id            int
	elem          int
	requesterChan chan OperationReply
}

func (i insert) ID() int {
	return i.id
}

func (i insert) Elem() int {
	return i.elem
}

func (i insert) RequesterChan() chan OperationReply {
	return i.requesterChan
}

func (i insert) Perform(node *binaryTreeNode) {
	if i.elem == node.elem {
		node.removed = false
		i.requesterChan <- operationFinished{i.id}
	} else if i.elem < node.elem {
		if node.left != nil {
			node.leftChan() <- i
		} else {
			node.left = makebinaryTreeNode(i.elem, false)
			node.left.parent = node.childReply
			i.requesterChan <- operationFinished{i.id}
		}
	} else if i.elem > node.elem {
		if node.right != nil {
			node.rightChan() <- i
		} else {
			node.right = makebinaryTreeNode(i.elem, false)
			node.right.parent = node.childReply
			i.requesterChan <- operationFinished{i.id}
		}
	}
}

func (i insert) String() string {
	return fmt.Sprintf("insert(id: %d, elem: %d)", i.id, i.elem)
}
