package main

//max two children
type BinaryTreeSet struct {
	parentChan chan Operation
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

	return b.root.parentChan
}

func (b *BinaryTreeSet) transferRootChan() chan Operation {
	if b.root == nil {
		return nil
	}

	return b.transferRoot.parentChan
}

func makeBinaryTreeSet() *BinaryTreeSet {
	x := BinaryTreeSet{make(chan Operation, 256), make(chan OperationReply, 16), makeBinaryTreeNode(0, true), nil, make(chan bool, 16), make(chan Operation, 1024), -1}
	go x.Run()
	return &x
}

func (b *BinaryTreeSet) Run() {
	for {
		select {
		case op := <-b.parentChan:
			b.rootChan() <- op
		case <-b.childReply:
			panic("received a child reply at root that should only happen while gcing")
		case gc := <-b.gcChan:
			if gc {
				b.transferRoot = makeBinaryTreeNode(0, true)
				b.runGC()
			}
		default:
		}
	}
}

func (b *BinaryTreeSet) runGC() {
	for {
		select {
		case op := <-b.parentChan:
			switch op.(type) {
			case CopyInsert:
				b.transferRootChan() <- op
			default:
				b.pendingQueue <- op
			}
		case opRep := <-b.childReply:
			switch opRep.(type) {
			case OperationFinished:
				if opRep.Id() == b.currentId {
					b.currentId--
					b.gcChan <- false
				} else {
					panic("received bad id for OperationFinished")
				}
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
