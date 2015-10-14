package ActorBinaryTree

import "fmt"
import "sync"

type binaryTreeNode struct {
	parent chan operationReply

	opChan     chan operation
	childReply chan operationReply

	left  *binaryTreeNode //use left.opChan for sending operations to it
	right *binaryTreeNode //use right.opChan for sending operations to it

	elem    int
	removed bool

	//elements for tracking gc
	gcoperationResponses replyTracker
	getElemResponse      operationFinished
}

func (b *binaryTreeNode) String() string {
	return fmt.Sprintf("Node(elem: %d, removed: %t)", b.elem, b.removed)
}

func (b *binaryTreeNode) leftChan() chan operation {
	if b.left == nil {
		return nil
	}

	return b.left.opChan
}

func (b *binaryTreeNode) rightChan() chan operation {
	if b.right == nil {
		return nil
	}

	return b.right.opChan
}

func newBinaryTreeNode(element int, initiallyremoved bool) *binaryTreeNode {
	//TODO: Tweak buffer sizes
	x := binaryTreeNode{nil, make(chan operation, 32), make(chan operationReply, 32), nil, nil, element, initiallyremoved, replyTracker{make(map[int]bool), &sync.Mutex{}}, operationFinished{0}}
	go x.Run()
	return &x
}

func (b *binaryTreeNode) Run() {
	for {
		select {
		case op := <-b.opChan:
			op.Perform(b)
		case opRep := <-b.childReply:
			b.gcoperationResponses.receivedReply(opRep)
			if b.gcoperationResponses.checkAllReceived() {
				b.parent <- b.getElemResponse
				b.left = nil
				b.right = nil
				return
			}
		}
	}
}
