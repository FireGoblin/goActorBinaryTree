package ActorBinaryTree

import "fmt"

type remove struct {
	id            int
	elem          int
	requesterChan chan operationReply
}

func (r remove) ID() int {
	return r.id
}

func (r remove) Elem() int {
	return r.elem
}

func (r remove) RequesterChan() chan operationReply {
	return r.requesterChan
}

func (r remove) Perform(node *binaryTreeNode) {
	if r.elem == node.elem {
		node.removed = true
		r.requesterChan <- operationFinished{r.id}
	} else if r.elem < node.elem {
		if node.left != nil {
			node.leftChan() <- r
		} else {
			r.requesterChan <- operationFinished{r.id}
		}
	} else if r.elem > node.elem {
		if node.right != nil {
			node.rightChan() <- r
		} else {
			r.requesterChan <- operationFinished{r.id}
		}
	}
}

func (r remove) String() string {
	return fmt.Sprintf("remove(id: %d, elem: %d)", r.id, r.elem)
}
