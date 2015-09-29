package ActorBinaryTree

import "math"

//max two children
type BinaryTreeSet struct {
	opChan     chan Operation
	childReply chan OperationReply

	root         *BinaryTreeNode
	transferRoot *BinaryTreeNode

	currentId int

	done chan bool
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

func MakeBinaryTreeSet() *BinaryTreeSet {
	x := BinaryTreeSet{make(chan Operation, 1024), make(chan OperationReply, 32), makeBinaryTreeNode(0, true), nil, -1, make(chan bool, 1)}
	x.root.parent = x.childReply
	go x.Run()
	return &x
}

//non-blocking
func (b *BinaryTreeSet) Close() {
	b.done <- true
}

func (b *BinaryTreeSet) Run() {
	for {
		select {
		case <-b.done:
			return
		case op := <-b.opChan:
			_, ok := op.(GC)
			if ok {
				b.transferRoot = makeBinaryTreeNode(0, true)
				b.transferRoot.parent = b.childReply
				b.runGC()
			} else {
				b.rootChan() <- op
			}
		case <-b.childReply:
			panic("received a child reply at root that should only happen while gcing")
		}
	}
}

func (b *BinaryTreeSet) runGC() {
	b.rootChan() <- GetElems{b.currentId, b.childReply}
	for {
		select {
		case opRep := <-b.childReply:
			switch opRep.(type) {
			case OperationFinished:
				if opRep.Id() == b.currentId {
					b.currentId--
					if b.currentId == math.MinInt32 {
						b.currentId = -1
					}
					b.root = b.transferRoot
					b.transferRoot = nil
					return
				} else {
					panic("received bad id for OperationFinished")
				}
			case CopyInsert:
				c := opRep.(CopyInsert)
				b.transferRootChan() <- c
			default:
				panic("should only receive OperationFinished and CopyInsert in childReply at root")
			}
		}
	}
}
