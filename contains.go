package ActorBinaryTree

import "fmt"

type Contains struct {
	id            int
	elem          int
	requesterChan chan OperationReply
}

func (c Contains) Id() int {
	return c.id
}

func (c Contains) Elem() int {
	return c.elem
}

func (c Contains) RequesterChan() chan OperationReply {
	return c.requesterChan
}

func (c Contains) Perform(node *BinaryTreeNode) {
	if c.elem == node.elem {
		c.requesterChan <- ContainsResult{c.id, !node.removed}
	} else if c.elem < node.elem {
		if node.left != nil {
			node.leftChan() <- c
		} else {
			c.requesterChan <- ContainsResult{c.id, false}
		}
	} else {
		if node.right != nil {
			node.rightChan() <- c
		} else {
			c.requesterChan <- ContainsResult{c.id, false}
		}
	}
}

func (i Contains) String() string {
	return fmt.Sprintf("Contains(id: %d, elem: %d)", i.id, i.elem)
}
