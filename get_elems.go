package ActorBinaryTree

import (
	"fmt"
	"math"
)

type getElems struct {
	id int
	//elem          int
	requesterChan chan OperationReply
}

func (i getElems) ID() int {
	return i.id
}

func (i getElems) Elem() int {
	return 0
}

func (i getElems) RequesterChan() chan OperationReply {
	return i.requesterChan
}

func (i getElems) Perform(node *binaryTreeNode) {
	node.getElemResponse = operationFinished{i.ID()}
	if node.left != nil {
		op := getElems{math.MinInt32, i.requesterChan}
		node.gcOperationResponses.sentOp(op)
		node.leftChan() <- op
	}
	if node.right != nil {
		op := getElems{math.MaxInt32, i.requesterChan}
		node.gcOperationResponses.sentOp(op)
		node.rightChan() <- op
	}
	if !node.removed {
		op := Copyinsert{node.elem, node.elem, node.childReply}
		node.gcOperationResponses.sentOp(op)
		i.requesterChan <- op
	}

	if node.gcOperationResponses.checkAllReceived() {
		node.parent <- node.getElemResponse
	}
}

func (i getElems) String() string {
	return fmt.Sprintf("getElems(id: %d)", i.id)
}
