package ActorBinaryTree

import "fmt"
import "sync"

type binaryTreeNode struct {
	parent chan OperationReply

	opChan     chan Operation
	childReply chan OperationReply

	left  *binaryTreeNode //use left.opChan for sending operations to it
	right *binaryTreeNode //use right.opChan for sending operations to it

	elem    int
	removed bool

	//elements for tracking gc
	gcOperationResponses ReplyTracker
	getElemResponse      operationFinished
}

func (b *binaryTreeNode) String() string {
	return fmt.Sprintf("Node(elem: %d, removed: %t)", b.elem, b.removed)
}

func (b *binaryTreeNode) leftChan() chan Operation {
	if b.left == nil {
		return nil
	}

	return b.left.opChan
}

func (b *binaryTreeNode) rightChan() chan Operation {
	if b.right == nil {
		return nil
	}

	return b.right.opChan
}

func makebinaryTreeNode(element int, initiallyremoved bool) *binaryTreeNode {
	//TODO: Tweak buffer sizes
	x := binaryTreeNode{nil, make(chan Operation, 1024), make(chan OperationReply, 32), nil, nil, element, initiallyremoved, ReplyTracker{make(map[int]bool), &sync.Mutex{}}, operationFinished{0}}
	go x.Run()
	return &x
}

func (b *binaryTreeNode) Run() {
	for {
		select {
		case op := <-b.opChan:
			op.Perform(b)
		case opRep := <-b.childReply:
			b.gcOperationResponses.receivedReply(opRep)
			if b.gcOperationResponses.checkAllReceived() {
				b.parent <- b.getElemResponse
				b.left = nil
				b.right = nil
				return
			}
		}
	}
}
