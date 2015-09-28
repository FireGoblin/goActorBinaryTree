package ActorBinaryTree

import (
	"fmt"
	"testing"
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

// 	one := testProbe.makeInsert(1)
// 	two := testProbe.makeContains(2)
// 	three := testProbe.makeRemove(1)
// 	four := testProbe.makeInsert(2)
// 	five := testProbe.makeContains(1)
// 	six := testProbe.makeContains(2)

// 	testProbe.sendOperation(one)
// 	testProbe.sendOperation(two)
// 	testProbe.sendOperation(three)
// 	testProbe.sendOperation(four)
// 	testProbe.sendOperation(five)
// 	testProbe.sendOperation(six)

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

	count := 10000

	for i := 0; i < count; i++ {
		err := testProbe.sendOperation(testProbe.randomOperation())
		if err != nil {
			t.FailNow()
		}

		if (i+1)%100 == 0 {
			fmt.Println("*%d*", i+1)
		}
		// if testProbe.rng.Float32() < 0.05 {
		// 	fmt.Println("Sending gc")
		// 	testProbe.sendGC()
		// }
	}

	fmt.Println("waiting")

	select {
	case <-succeed:
	case <-fail:
		t.FailNow()
	}
}
