package ActorBinaryTree

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

type TestProbe struct {
	childReply chan OperationReply

	tree *BinaryTreeSet

	currentTree       map[int]bool //track what has been inserted
	expectedResponses map[int]bool //track expected responses to contains requests
	finishedResponses ReplyTracker

	currentId int

	rng *rand.Rand

	responseCount int
}

func (t *TestProbe) displayUnreceived() {
	t.finishedResponses.displayUnreceived()
}

func (t *TestProbe) Run(succeed chan bool, fail chan bool) {
	for {
		select {
		case msg := <-t.childReply:
			//fmt.Println(msg)
			t.responseCount++
			// if t.responseCount%100 == 0 {
			// 	fmt.Println(t.responseCount)
			// }
			if !t.checkReply(msg) {
				fail <- true
			}
		case <-time.After(10 * time.Second):
			fmt.Println("checking all replies received")
			if !t.checkReceviedAllResponses() {
				fmt.Println("not all responses found")
				t.displayUnreceived()
				fail <- true
			} else {
				fmt.Println("all responses received")
				succeed <- true
			}
		}
	}
}

func makeTestProbe() *TestProbe {
	x := TestProbe{make(chan OperationReply, 1024), makeBinaryTreeSet(), make(map[int]bool), make(map[int]bool), make(map[int]bool), 1, rand.New(rand.NewSource(777)), 0}
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
	default:
		return fmt.Errorf("unknown operation found in test probe")
	}

	t.childChan() <- o

	return nil
}

func (t *TestProbe) sendGC() {
	t.childChan() <- GC{}
}

func (t *TestProbe) checkReply(o OperationReply) bool {
	//fmt.Println(o)
	switch o.(type) {
	case OperationFinished:
		if t.finishedResponses[o.Id()] {
			t.finishedResponses.receivedReply(o)
			return true
		} else {
			return false
		}
	case ContainsResult:
		c := o.(ContainsResult)
		if t.finishedResponses[c.Id()] && t.expectedResponses[c.Id()] == c.Result() {
			t.finishedResponses.receivedReply(c)
			return true
		} else {
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
