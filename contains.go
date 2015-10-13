package ActorBinaryTree

import "fmt"

type contains struct {
	id            int
	elem          int
	requesterChan chan operationReply
}

func (c contains) ID() int {
	return c.id
}

func (c contains) Elem() int {
	return c.elem
}

func (c contains) RequesterChan() chan operationReply {
	return c.requesterChan
}

func (c contains) Perform(node *binaryTreeNode) {
	if c.elem == node.elem {
		c.requesterChan <- containsResult{c.id, !node.removed}
	} else if c.elem < node.elem {
		if node.left != nil {
			node.leftChan() <- c
		} else {
			c.requesterChan <- containsResult{c.id, false}
		}
	} else {
		if node.right != nil {
			node.rightChan() <- c
		} else {
			c.requesterChan <- containsResult{c.id, false}
		}
	}
}

func (c contains) String() string {
	return fmt.Sprintf("contains(id: %d, elem: %d)", c.id, c.elem)
}
