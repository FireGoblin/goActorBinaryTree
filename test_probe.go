package ActorBinaryTree

import (
	"fmt"
	"math"
	"math/rand"
	"sync"
	"time"
)

type TestProbe struct {
	opChan     chan Operation
	childReply chan OperationReply

	tree *BinaryTreeSet

	currentTree       map[int]bool //track what has been inserted
	expectedResponses map[int]bool //track expected responses to contains requests
	finishedResponses ReplyTracker

	currentId int

	rng *rand.Rand

	replyCount int
}

func (t *TestProbe) displayUnreceived() {
	t.finishedResponses.displayUnreceived()
}

func (t *TestProbe) Run(succeed chan int, fail chan int) {
	for {
		select {
		case msg := <-t.childReply:
			//fmt.Println(msg)
			t.replyCount++
			if !t.checkReply(msg) {
				fail <- t.replyCount
				fmt.Println("failed in checkReply")
				return
			}
		case op := <-t.opChan:
			t.sendOperation(op)
		case <-time.After(1 * time.Second):
			fmt.Println("checking all replies received")
			if !t.checkReceviedAllResponses() {
				fmt.Println("not all responses found")
				t.displayUnreceived()
				fail <- t.replyCount
				return
			} else {
				fmt.Println("all responses received")
				succeed <- t.replyCount
				return
			}
		}
	}
}

func makeTestProbe() *TestProbe {
	x := TestProbe{make(chan Operation, 1024), make(chan OperationReply, 1024), makeBinaryTreeSet(), make(map[int]bool), make(map[int]bool), ReplyTracker{make(map[int]bool), &sync.Mutex{}}, 1, rand.New(rand.NewSource(777)), 0}
	return &x
}

func (t *TestProbe) childChan() chan Operation {
	return t.tree.opChan
}

func (t *TestProbe) sendOperation(o Operation) error {
	//fmt.Println(o)
	switch o.(type) {
	case Insert:
		t.currentTree[o.Elem()] = true
		t.finishedResponses.sentOp(o)
	case Remove:
		t.currentTree[o.Elem()] = false
		t.finishedResponses.sentOp(o)
	case Contains:
		t.expectedResponses[o.Id()] = t.currentTree[o.Elem()]
		t.finishedResponses.sentOp(o)
	case GC:
	default:
		return fmt.Errorf("unknown operation found in test probe")
	}

	t.childChan() <- o

	return nil
}

func (t *TestProbe) injectOperation(o Operation) {
	t.opChan <- o
}

func (t *TestProbe) injectGC() {
	t.opChan <- GC{}
}

func (t *TestProbe) checkReply(o OperationReply) bool {
	//fmt.Println(o)
	switch o.(type) {
	case OperationFinished:
		if t.finishedResponses.get(o.Id()) {
			t.finishedResponses.receivedReply(o)
			return true
		} else {
			fmt.Println("failing reply", o)
			return false
		}
	case ContainsResult:
		c := o.(ContainsResult)
		if t.finishedResponses.get(c.Id()) && t.expectedResponses[c.Id()] == c.Result() {
			t.finishedResponses.receivedReply(c)
			return true
		} else {
			fmt.Println("failing reply", o)
			return false
		}
	default:
		return false
	}
}

func (t *TestProbe) checkReceviedAllResponses() bool {
	return t.finishedResponses.checkAllReceived()
}

func (t *TestProbe) incrementId() {
	t.currentId++
	if t.currentId == math.MaxInt32 {
		t.currentId = 1
	}
}

func (t *TestProbe) makeInsert(e int) Insert {
	i := t.currentId
	t.incrementId()
	return Insert{i, e, t.childReply}
}

func (t *TestProbe) makeContains(e int) Contains {
	i := t.currentId
	t.incrementId()
	return Contains{i, e, t.childReply}
}

func (t *TestProbe) makeRemove(e int) Remove {
	i := t.currentId
	t.incrementId()
	return Remove{i, e, t.childReply}
}

func (t *TestProbe) coinFlip() bool {
	return t.rng.Int()%2 == 0
}

// random element between (MinInt8, MaxInt8) (exclusive)
func (t *TestProbe) randomElement8() int {
	flip := t.coinFlip()
	val := t.rng.Int31n(math.MaxInt8)

	if flip {
		if val == 0 {
			val = math.MinInt8 + 1
		} else {
			val = -val
		}
	}

	return int(val)
}

// random element between (MinInt16, MaxInt16) (exclusive)
func (t *TestProbe) randomElement16() int {
	flip := t.coinFlip()
	val := t.rng.Int31n(math.MaxInt16)

	if flip {
		if val == 0 {
			val = math.MinInt16 + 1
		} else {
			val = -val
		}
	}

	return int(val)
}

// random element between (MinInt32, MaxInt32) (exclusive)
func (t *TestProbe) randomElement32() int {
	flip := t.coinFlip()
	val := t.rng.Int31n(math.MaxInt32)

	if flip {
		if val == 0 {
			val = math.MinInt32 + 1
		} else {
			val = -val
		}
	}

	return int(val)
}

func (t *TestProbe) randomOperation() Operation {
	val := t.rng.Int31n(4)

	switch val {
	case 0:
		return t.makeInsert(t.randomElement8())
	case 1:
		return t.makeInsert(t.randomElement8())
	case 2:
		return t.makeContains(t.randomElement8())
	case 3:
		return t.makeRemove(t.randomElement8())
	}

	return nil
}
