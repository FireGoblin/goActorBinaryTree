package main

import (
	"fmt"
	"math/rand"
	"testing"
)

type TestProbe struct {
	childReply chan OperationReply
	done       chan bool

	tree *BinaryTreeSet

	currentTree       map[int]bool //track what has been inserted
	expectedResponses map[int]bool //track expected responses to contains requests
	finishedResponses map[int]bool

	currentId int

	rng *rand.Rand
}

func (t *TestProbe) Run(test *testing.T) {
	for {
		select {
		case msg := <-t.childReply:
			fmt.Println(msg)
			if !t.checkReply(msg) {
				test.FailNow()
			}
		case <-t.done:
			fmt.Println("checking all replies received")
			if !t.checkReceviedAllResponses() {
				test.FailNow()
			}
		}
	}
}

func makeTestProbe() *TestProbe {
	x := TestProbe{make(chan OperationReply, 256), make(chan bool), makeBinaryTreeSet(), make(map[int]bool), make(map[int]bool), make(map[int]bool), 1, rand.New(rand.NewSource(777))}
	return &x
}

func (t *TestProbe) childChan() chan Operation {
	return t.tree.parentChan
}

func (t *TestProbe) sendOperation(o Operation) error {
	switch o.(type) {
	case Insert:
		t.currentTree[o.Elem()] = true
		t.finishedResponses[o.Id()] = true
	case Remove:
		t.currentTree[o.Elem()] = false
		t.finishedResponses[o.Id()] = true
	case Contains:
		t.expectedResponses[o.Id()] = t.currentTree[o.Elem()]
		t.finishedResponses[o.Id()] = true
	default:
		return fmt.Errorf("unknown operation found in test probe")
	}

	t.childChan() <- o

	return nil
}

func (t *TestProbe) checkReply(o OperationReply) bool {
	switch o.(type) {
	case OperationFinished:
		if t.finishedResponses[o.Id()] {
			t.finishedResponses[o.Id()] = false
			return true
		} else {
			return false
		}
	case ContainsResult:
		c := o.(ContainsResult)
		if t.finishedResponses[c.Id()] && t.expectedResponses[c.Id()] == c.Result() {
			t.finishedResponses[c.Id()] = false
			return true
		} else {
			return false
		}
	default:
		return false
	}
}

func (t *TestProbe) checkReceviedAllResponses() bool {
	for _, v := range t.finishedResponses {
		if v {
			return false
		}
	}

	return true
}

func (t *TestProbe) makeInsert(e int) Insert {
	i := t.currentId
	t.currentId++
	return Insert{i, e, t.childReply}
}

func (t *TestProbe) makeContains(e int) Contains {
	i := t.currentId
	t.currentId++
	return Contains{i, e, t.childReply}
}

func (t *TestProbe) makeRemove(e int) Remove {
	i := t.currentId
	t.currentId++
	return Remove{i, e, t.childReply}
}

// func (t *TestProbe) coinFlip() bool {
// 	return t.rng.Int()%2 == 0
// }

// func (t *TestProbe) negativeRand() int {
// 	x := t.rng.Int31()

// 	if t.coinFlip() {
// 		x = -x
// 	}

// 	return int(x)
// }

// func (t *testProbe) makeInsert() Operation {
// 	x := Insert{t.currentId, t.negativeRand(), childReply}
// 	return x
// }

// func (t *testProbe) makeContains() Operation {
// 	// if greater than 0 then pick a number that exists
// 	treeSize := len(t.currentTree)
// 	if t.coinFlip()  && treeSize > 0 {
// 		x := t.rng.Int31n(treeSize)

// 	}
// }
