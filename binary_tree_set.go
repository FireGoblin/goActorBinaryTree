package ActorBinaryTree

import "math"

//max two children
type BinaryTreeSet struct {
	opChan     chan Operation
	childReply chan OperationReply

	root         *BinaryTreeNode
	transferRoot *BinaryTreeNode

	gcChan       chan bool
	pendingQueue chan Operation

	currentId int
}

func (b *BinaryTreeSet) rootChan() chan Operation {
	if b.root == nil {
		return nil
	}

	return b.root.opChan
}

func (b *BinaryTreeSet) transferRootChan() chan Operation {
	if b.root == nil {
		return nil
	}

	return b.transferRoot.opChan
}

func makeBinaryTreeSet() *BinaryTreeSet {
	x := BinaryTreeSet{make(chan Operation, 1024), make(chan OperationReply, 32), makeBinaryTreeNode(0, true), nil, make(chan bool, 32), make(chan Operation, 1024), -1}
	x.root.parent = x.childReply
	go x.Run()
	return &x
}

func (b *BinaryTreeSet) Run() {
	for {
		select {
		case op := <-b.opChan:
			b.rootChan() <- op
		case <-b.childReply:
			panic("received a child reply at root that should only happen while gcing")
		case gc := <-b.gcChan:
			if gc {
				b.transferRoot = makeBinaryTreeNode(0, true)
				b.transferRoot.parent = b.childReply
				b.runGC()
			}
		default:
		}
	}
}

func (b *BinaryTreeSet) runGC() {
	for {
		select {
		case op := <-b.opChan:
			b.pendingQueue <- op
		case opRep := <-b.childReply:
			switch opRep.(type) {
			case OperationFinished:
				if opRep.Id() == b.currentId {
					b.currentId--
					if b.currentId == math.MinInt32 {
						b.currentId = -1
					}
					b.gcChan <- false
				} else {
					panic("received bad id for OperationFinished")
				}
			case CopyInsert:
				c := opRep.(CopyInsert)
				b.transferRootChan() <- c
			default:
				panic("should only receive OperationFinished in childReply at root")
			}
		case gc := <-b.gcChan:
			if !gc {
				b.root = b.transferRoot
				b.transferRoot = nil
				for op := range b.pendingQueue {
					b.rootChan() <- op
				}
				return
			}
		default:
		}
	}
}
