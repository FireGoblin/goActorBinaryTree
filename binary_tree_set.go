package ActorBinaryTree

import "math"

// BinaryTreeSet
type BinaryTreeSet struct {
	opChan     chan operation
	childReply chan operationReply

	root         *binaryTreeNode
	transferRoot *binaryTreeNode

	currentID int

	done chan bool
}

func (b *BinaryTreeSet) rootChan() chan operation {
	if b.root == nil {
		return nil
	}

	return b.root.opChan
}

func (b *BinaryTreeSet) transferRootChan() chan operation {
	if b.root == nil {
		return nil
	}

	return b.transferRoot.opChan
}

func MakeBinaryTreeSet() *BinaryTreeSet {
	x := BinaryTreeSet{make(chan operation, 1024), make(chan operationReply, 32), makeBinaryTreeNode(0, true), nil, -1, make(chan bool, 1)}
	x.root.parent = x.childReply
	go x.run()
	return &x
}

// Close in a non-blocking manner
func (b *BinaryTreeSet) Close() {
	b.done <- true
}

func (b *BinaryTreeSet) run() {
	for {
		select {
		case <-b.done:
			return
		case op := <-b.opChan:
			_, ok := op.(gc)
			if ok {
				b.transferRoot = makeBinaryTreeNode(0, true)
				b.transferRoot.parent = b.childReply
				b.rungc()
			} else {
				b.rootChan() <- op
			}
		case <-b.childReply:
			panic("received a child reply at root that should only happen while gcing")
		}
	}
}

func (b *BinaryTreeSet) rungc() {
	b.rootChan() <- getElems{b.currentID, b.childReply}
	for {
		select {
		case opRep := <-b.childReply:
			switch opRep.(type) {
			case operationFinished:
				if opRep.ID() == b.currentID {
					b.currentID--
					if b.currentID == math.MinInt32 {
						b.currentID = -1
					}
					b.root = b.transferRoot
					b.transferRoot = nil
					return
				}

				panic("received bad id for operationFinished")
			case copyInsert:
				c := opRep.(copyInsert)
				b.transferRootChan() <- c
			default:
				panic("should only receive operationFinished and copyInsert in childReply at root")
			}
		}
	}
}
