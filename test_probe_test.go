package ActorBinaryTree

import (
	"fmt"
	"testing"
	"time"
)

func TestinsertsAndcontains(t *testing.T) {
	testProbe := makeTestProbe()

	one := testProbe.makecontains(1)
	testProbe.sendoperation(one)
	oneResult := <-testProbe.childReply

	x, ok := oneResult.(containsResult)
	if ok {
		if x.ID() != 1 || x.Result() || !testProbe.checkReply(x) {
			t.FailNow()
		}
	} else {
		t.FailNow()
	}

	two := testProbe.makeinsert(1)
	three := testProbe.makecontains(1)
	testProbe.sendoperation(two)
	testProbe.sendoperation(three)

	twoResult := <-testProbe.childReply
	threeResult := <-testProbe.childReply

	y, ok := twoResult.(operationFinished)
	if ok {
		if y.ID() != 2 || !testProbe.checkReply(y) {
			t.FailNow()
		}
	} else {
		t.FailNow()
	}

	x, ok = threeResult.(containsResult)
	if ok {
		if x.ID() != 3 || !x.Result() || !testProbe.checkReply(x) {
			t.FailNow()
		}
	} else {
		t.FailNow()
	}

	fmt.Println("inserts and contains succeeded")
}

func TestInstructionExample(t *testing.T) {
	fmt.Println("--------------------------")
	testProbe := makeTestProbe()

	succeed := make(chan int)
	fail := make(chan int)

	go testProbe.Run(succeed, fail)

	var ops []operation
	ops = append(ops, testProbe.makeinsert(1))
	ops = append(ops, testProbe.makecontains(2))
	ops = append(ops, testProbe.makeremove(1))
	ops = append(ops, testProbe.makeinsert(2))
	ops = append(ops, testProbe.makecontains(1))
	ops = append(ops, testProbe.makecontains(2))

	for _, op := range ops {
		testProbe.sendoperation(op)
		time.Sleep(1 * time.Millisecond)
	}

	select {
	case <-succeed:
	case <-fail:
		t.FailNow()
	}
}

func TestSmallgc(t *testing.T) {
	fmt.Println("--------------------------")
	testProbe := makeTestProbe()

	succeed := make(chan int)
	fail := make(chan int)

	go testProbe.Run(succeed, fail)

	var opsBefore []operation
	var opsAfter []operation

	opsBefore = append(opsBefore, testProbe.makeinsert(-122))
	opsBefore = append(opsBefore, testProbe.makeinsert(99))
	opsBefore = append(opsBefore, testProbe.makeinsert(-13))
	opsBefore = append(opsBefore, testProbe.makeinsert(104))
	opsBefore = append(opsBefore, testProbe.makeremove(-122))

	opsAfter = append(opsAfter, testProbe.makecontains(-122))
	opsAfter = append(opsAfter, testProbe.makecontains(99))
	opsAfter = append(opsAfter, testProbe.makecontains(-13))
	opsAfter = append(opsAfter, testProbe.makecontains(104))
	opsAfter = append(opsAfter, testProbe.makecontains(777))

	for _, op := range opsBefore {
		testProbe.sendoperation(op)
		time.Sleep(1 * time.Millisecond)
	}

	testProbe.injectgc()
	time.Sleep(1 * time.Millisecond)

	for _, op := range opsAfter {
		testProbe.sendoperation(op)
		time.Sleep(1 * time.Millisecond)
	}

	select {
	case <-succeed:
	case <-fail:
		t.FailNow()
	}
}

func TestWorkWithgc(t *testing.T) {
	fmt.Println("--------------------------")
	testProbe := makeTestProbe()

	succeed := make(chan int)
	fail := make(chan int)

	go testProbe.Run(succeed, fail)

	count := 100000

	start := time.Now()

	for i := 0; i < count; i++ {
		op := testProbe.randomoperation()

		testProbe.injectoperation(op)

		if testProbe.rng.Float32() < 0.01 {
			testProbe.injectgc()
		}
	}

	elapsed := time.Since(start)
	fmt.Println("sending messages took", elapsed)

	fmt.Println("waiting")

	start = time.Now()

	select {
	case c := <-succeed:
		fmt.Println(c)
		if c != count {
			t.FailNow()
		}
	case c := <-fail:
		fmt.Println(c)
		t.FailNow()
	}

	elapsed = time.Since(start)
	fmt.Println("waited for response", elapsed)
}
