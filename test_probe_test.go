package ActorBinaryTree

import (
	"fmt"
	"testing"
	"time"
)

// func TestInsertsAndContains(t *testing.T) {
// 	testProbe := makeTestProbe()

// 	one := testProbe.makeContains(1)
// 	testProbe.sendOperation(one)
// 	oneResult := <-testProbe.childReply

// 	x, ok := oneResult.(ContainsResult)
// 	if ok {
// 		if x.Id() != 1 || x.Result() || !testProbe.checkReply(x) {
// 			t.FailNow()
// 		}
// 	} else {
// 		t.FailNow()
// 	}

// 	two := testProbe.makeInsert(1)
// 	three := testProbe.makeContains(1)
// 	testProbe.sendOperation(two)
// 	testProbe.sendOperation(three)

// 	twoResult := <-testProbe.childReply
// 	threeResult := <-testProbe.childReply

// 	y, ok := twoResult.(OperationFinished)
// 	if ok {
// 		if y.Id() != 2 || !testProbe.checkReply(y) {
// 			t.FailNow()
// 		}
// 	} else {
// 		t.FailNow()
// 	}

// 	x, ok = threeResult.(ContainsResult)
// 	if ok {
// 		if x.Id() != 3 || !x.Result() || !testProbe.checkReply(x) {
// 			t.FailNow()
// 		}
// 	} else {
// 		t.FailNow()
// 	}

// 	fmt.Println("Inserts and Contains succeeded")
// }

// func TestInstructionExample(t *testing.T) {
// 	fmt.Println("--------------------------")
// 	testProbe := makeTestProbe()

// 	succeed := make(chan bool)
// 	fail := make(chan bool)

// 	go testProbe.Run(succeed, fail)

// 	var ops []Operation
// 	ops = append(ops, testProbe.makeInsert(1))
// 	ops = append(ops, testProbe.makeContains(2))
// 	ops = append(ops, testProbe.makeRemove(1))
// 	ops = append(ops, testProbe.makeInsert(2))
// 	ops = append(ops, testProbe.makeContains(1))
// 	ops = append(ops, testProbe.makeContains(2))

// 	for _, op := range ops {
// 		testProbe.sendOperation(op)
// 		time.Sleep(1 * time.Millisecond)
// 	}

// 	select {
// 	case <-succeed:
// 	case <-fail:
// 		t.FailNow()
// 	}
// }

// func TestSmallGC(t *testing.T) {
// 	fmt.Println("--------------------------")
// 	testProbe := makeTestProbe()

// 	succeed := make(chan bool)
// 	fail := make(chan bool)

// 	go testProbe.Run(succeed, fail)

// 	var opsBefore []Operation
// 	var opsAfter []Operation

// 	opsBefore = append(opsBefore, testProbe.makeInsert(-122))
// 	opsBefore = append(opsBefore, testProbe.makeInsert(99))
// 	opsBefore = append(opsBefore, testProbe.makeInsert(-13))
// 	opsBefore = append(opsBefore, testProbe.makeInsert(104))
// 	opsBefore = append(opsBefore, testProbe.makeRemove(-122))

// 	opsAfter = append(opsAfter, testProbe.makeContains(-122))
// 	opsAfter = append(opsAfter, testProbe.makeContains(99))
// 	opsAfter = append(opsAfter, testProbe.makeContains(-13))
// 	opsAfter = append(opsAfter, testProbe.makeContains(104))
// 	opsAfter = append(opsAfter, testProbe.makeContains(777))

// 	for _, op := range opsBefore {
// 		testProbe.sendOperation(op)
// 		time.Sleep(1 * time.Millisecond)
// 	}

// 	testProbe.sendGC()
// 	time.Sleep(1 * time.Millisecond)

// 	for _, op := range opsAfter {
// 		testProbe.sendOperation(op)
// 		time.Sleep(1 * time.Millisecond)
// 	}

// 	select {
// 	case <-succeed:
// 	case <-fail:
// 		t.FailNow()
// 	}
// }

func TestWorkWithGC(t *testing.T) {
	fmt.Println("--------------------------")
	testProbe := makeTestProbe()

	succeed := make(chan bool)
	fail := make(chan bool)

	go testProbe.Run(succeed, fail)

	count := 1000

	//display := true

	start := time.Now()

	for i := 0; i < count; i++ {
		op := testProbe.randomOperation()

		err := testProbe.sendOperation(op)
		if err != nil {
			t.FailNow()
		}

		if testProbe.rng.Float32() < 0.05 {
			//testProbe.sendGC()
		}
	}

	elapsed := time.Since(start)
	fmt.Println("sending messages took", elapsed)

	fmt.Println("waiting")

	start = time.Now()

	select {
	case <-succeed:
	case <-fail:
		t.FailNow()
	}

	elapsed = time.Since(start)
	fmt.Println("waited for response", elapsed)
}
