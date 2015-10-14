package ActorBinaryTree

import (
	"fmt"
	"math"
	"math/rand"
	"sync"
	"time"
)

type testProbe struct {
	opChan     chan operation
	childReply chan operationReply

	tree *BinaryTreeSet

	currentTree       map[int]bool //track what has been inserted
	expectedResponses map[int]bool //track expected responses to contains requests
	finishedResponses replyTracker

	currentID int

	rng *rand.Rand

	replyCount int
}

func (t *testProbe) displayUnreceived() {
	t.finishedResponses.displayUnreceived()
}

func (t *testProbe) run(succeed chan int, fail chan int) {
	for {
		select {
		case msg := <-t.childReply:
			t.replyCount++
			if !t.checkReply(msg) {
				fail <- t.replyCount
				fmt.Println("failed in checkReply")
				return
			}
		case op := <-t.opChan:
			t.sendoperation(op)
		case <-time.After(1 * time.Second):
			fmt.Println("checking all replies received")
			if !t.checkReceviedAllResponses() {
				fmt.Println("not all responses found")
				t.displayUnreceived()
				fail <- t.replyCount
				return
			}

			fmt.Println("all responses received")
			succeed <- t.replyCount
			return
		}
	}
}

func maketestProbe() *testProbe {
	x := testProbe{make(chan operation, 1024), make(chan operationReply, 1024), NewBinaryTreeSet(), make(map[int]bool), make(map[int]bool), replyTracker{make(map[int]bool), &sync.Mutex{}}, 1, rand.New(rand.NewSource(777)), 0}
	return &x
}

func (t *testProbe) childChan() chan operation {
	return t.tree.opChan
}

func (t *testProbe) sendoperation(o operation) error {
	switch o.(type) {
	case insert:
		t.currentTree[o.Elem()] = true
		t.finishedResponses.sentOp(o)
	case remove:
		t.currentTree[o.Elem()] = false
		t.finishedResponses.sentOp(o)
	case contains:
		t.expectedResponses[o.ID()] = t.currentTree[o.Elem()]
		t.finishedResponses.sentOp(o)
	case gc:
	default:
		return fmt.Errorf("unknown operation found in test probe")
	}

	t.childChan() <- o

	return nil
}

func (t *testProbe) injectoperation(o operation) {
	t.opChan <- o
}

func (t *testProbe) injectgc() {
	t.opChan <- gc{}
}

func (t *testProbe) checkReply(o operationReply) bool {
	switch o.(type) {
	case operationFinished:
		if t.finishedResponses.get(o.ID()) {
			t.finishedResponses.receivedReply(o)
			return true
		}

		fmt.Println("failing reply", o)
		return false
	case containsResult:
		c := o.(containsResult)
		if t.finishedResponses.get(c.ID()) && t.expectedResponses[c.ID()] == c.Result() {
			t.finishedResponses.receivedReply(c)
			return true
		}

		fmt.Println("failing reply", o)
		return false
	default:
		return false
	}
}

func (t *testProbe) checkReceviedAllResponses() bool {
	return t.finishedResponses.checkAllReceived()
}

func (t *testProbe) incrementID() {
	t.currentID++
	if t.currentID == math.MaxInt32 {
		t.currentID = 1
	}
}

func (t *testProbe) makeinsert(e int) insert {
	i := t.currentID
	t.incrementID()
	return insert{i, e, t.childReply}
}

func (t *testProbe) makecontains(e int) contains {
	i := t.currentID
	t.incrementID()
	return contains{i, e, t.childReply}
}

func (t *testProbe) makeremove(e int) remove {
	i := t.currentID
	t.incrementID()
	return remove{i, e, t.childReply}
}

func (t *testProbe) coinFlip() bool {
	return t.rng.Int()%2 == 0
}

// random element between (MinInt8, MaxInt8) (exclusive)
func (t *testProbe) randomElement8() int {
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
func (t *testProbe) randomElement16() int {
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
func (t *testProbe) randomElement32() int {
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

func (t *testProbe) randomoperation() operation {
	val := t.rng.Int31n(4)

	switch val {
	case 0:
		return t.makeinsert(t.randomElement8())
	case 1:
		return t.makeinsert(t.randomElement8())
	case 2:
		return t.makecontains(t.randomElement8())
	case 3:
		return t.makeremove(t.randomElement8())
	}

	return nil
}
