package ActorBinaryTree

import "fmt"

type Insert struct {
	id            int
	elem          int
	requesterChan chan OperationReply
}

func (i Insert) Id() int {
	return i.id
}

func (i Insert) Elem() int {
	return i.elem
}

func (i Insert) RequesterChan() chan OperationReply {
	return i.requesterChan
}

func (i Insert) Perform(node *BinaryTreeNode) {
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

func (i Insert) String() string {
	return fmt.Sprintf("Insert(id: %d, elem: %d)", i.id, i.elem)
}
