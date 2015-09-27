package main

import (
	//"fmt"
	"testing"
	"time"
)

func TestInsertsAndContains(t *testing.T) {
	testProbe := makeTestProbe()

	one := testProbe.makeContains(1)
	testProbe.sendOperation(one)
	oneResult := <-testProbe.childReply

	x, ok := oneResult.(ContainsResult)
	if ok {
		if x.Id() != 1 || x.Result() || !testProbe.checkReply(x) {
			t.FailNow()
		}
	} else {
		t.FailNow()
	}

	two := testProbe.makeInsert(1)
	three := testProbe.makeContains(1)
	testProbe.sendOperation(two)
	testProbe.sendOperation(three)

	twoResult := <-testProbe.childReply
	threeResult := <-testProbe.childReply

	y, ok := twoResult.(OperationFinished)
	if ok {
		if y.Id() != 2 || !testProbe.checkReply(y) {
			t.FailNow()
		}
	} else {
		t.FailNow()
	}

	x, ok = threeResult.(ContainsResult)
	if ok {
		if x.Id() != 3 || !x.Result() || !testProbe.checkReply(x) {
			t.FailNow()
		}
	} else {
		t.FailNow()
	}
}

func TestInstructionExample(t *testing.T) {
	testProbe := makeTestProbe()

	go testProbe.Run(t)

	one := testProbe.makeInsert(1)
	two := testProbe.makeContains(2)
	three := testProbe.makeRemove(1)
	four := testProbe.makeInsert(2)
	five := testProbe.makeContains(1)
	six := testProbe.makeContains(2)

	testProbe.sendOperation(one)
	testProbe.sendOperation(two)
	testProbe.sendOperation(three)
	testProbe.sendOperation(four)
	testProbe.sendOperation(five)
	testProbe.sendOperation(six)

	time.Sleep(100 * time.Millisecond)

	testProbe.done <- true

	time.Sleep(10 * time.Millisecond)
}