package main

import (
	"fmt"
	"testing"
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
