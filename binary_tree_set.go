package ActorBinaryTree

import (
	"fmt"
	"math"
)

//max two children
type BinaryTreeSet struct {
	opChan     chan Operation
	childReply chan OperationReply

	root         *BinaryTreeNode
	transferRoot *BinaryTreeNode

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
	x := BinaryTreeSet{make(chan Operation, 1024), make(chan OperationReply, 32), makeBinaryTreeNode(0, true), nil, make(chan Operation, 1024), -1}
	x.root.parent = x.childReply
	go x.Run()
	return &x
}

func (b *BinaryTreeSet) Run() {
	for {
		select {
		case op := <-b.pendingQueue:
			_, ok := op.(GC)
			if ok {
				b.transferRoot = makeBinaryTreeNode(0, true)
				b.transferRoot.parent = b.childReply
				fmt.Println("moving to GC")
				b.runGC()
			} else {
				b.rootChan() <- op
			}
		case op := <-b.opChan:
			_, ok := op.(GC)
			if ok {
				b.transferRoot = makeBinaryTreeNode(0, true)
				b.transferRoot.parent = b.childReply
				fmt.Println("moving to GC")
				b.runGC()
			} else {
				b.rootChan() <- op
			}
		case <-b.childReply:
			panic("received a child reply at root that should only happen while gcing")
		default:
		}
	}
}

func (b *BinaryTreeSet) runGC() {
	b.rootChan() <- GetElems{b.currentId, b.childReply}
	for {
		select {
		case opRep := <-b.childReply:
			//fmt.Println("opRep received:", opRep)
			switch opRep.(type) {
			case OperationFinished:
				if opRep.Id() == b.currentId {
					b.currentId--
					if b.currentId == math.MinInt32 {
						b.currentId = -1
					}
					fmt.Println("GC Over")
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
				panic("should only receive OperationFinished in childReply at root")
			}
		case op := <-b.opChan:
			//fmt.Println("move to pending queue")
			b.pendingQueue <- op
		default:
		}
	}
}
