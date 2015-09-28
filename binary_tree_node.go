package ActorBinaryTree

import "fmt"
import "sync"

//max two children
type BinaryTreeNode struct {
	parent chan OperationReply

	opChan     chan Operation
	childReply chan OperationReply

	left  *BinaryTreeNode //use left.opChan for sending operations to it
	right *BinaryTreeNode //use right.opChan for sending operations to it

	elem    int
	removed bool

	//elements for tracking gc
	gcOperationResponses ReplyTracker
	getElemResponse      OperationFinished
}

func (b *BinaryTreeNode) String() string {
	return fmt.Sprintf("Node(elem: %d, removed: %t)", b.elem, b.removed)
}

func (b *BinaryTreeNode) leftChan() chan Operation {
	if b.left == nil {
		return nil
	}

	return b.left.opChan
}

func (b *BinaryTreeNode) rightChan() chan Operation {
	if b.right == nil {
		return nil
	}

	return b.right.opChan
}

func makeBinaryTreeNode(element int, initiallyRemoved bool) *BinaryTreeNode {
	//TODO: Tweak buffer sizes
	x := BinaryTreeNode{nil, make(chan Operation, 1024), make(chan OperationReply, 32), nil, nil, element, initiallyRemoved, ReplyTracker{make(map[int]bool), &sync.Mutex{}}, OperationFinished{0}}
	go x.Run()
	return &x
}

func (b *BinaryTreeNode) Run() {
	for {
		select {
		case op := <-b.opChan:
			op.Perform(b)
		case opRep := <-b.childReply:
			//fmt.Println(b)
			//fmt.Println(opRep)
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
