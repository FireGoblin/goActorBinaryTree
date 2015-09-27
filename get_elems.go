package main

import (
	"fmt"
	"math"
)

type GetElems struct {
	id int
	//elem          int
	requesterChan chan OperationReply
}

func (i GetElems) Id() int {
	return i.id
}

func (i GetElems) Elem() int {
	return 0
}

func (i GetElems) RequesterChan() chan OperationReply {
	return i.requesterChan
}

func (i GetElems) Perform(node *BinaryTreeNode) {
	node.getElemResponse = &OperationFinished{i.Id()}
	if node.left == nil {
		op := GetElems{math.MinInt32, i.requesterChan}
		node.gcOperationResponses.sentOp(op)
		node.leftChan() <- op
	}
	if node.right == nil {
		op := GetElems{math.MaxInt32, i.requesterChan}
		node.gcOperationResponses.sentOp(op)
		node.rightChan() <- op
	}
	if !node.removed {
		op := CopyInsert{node.elem, node.elem, node.childReply}
		node.gcOperationResponses.sentOp(op)
		i.requesterChan <- op
	}
}

func (i GetElems) String() string {
	return fmt.Sprintf("GetElems(id: %d)", i.id)
}