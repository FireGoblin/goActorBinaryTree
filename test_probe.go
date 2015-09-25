package main

type TestProbe struct {
	childReply chan OperationReply

	tree *BinaryTreeSet

	currentTree       map[int]bool //track what has been inserted
	expectedResponses map[int]bool //track expected responses to contains requests
}

func makeTestProbe() *TestProbe {
	x := TestProbe{make(chan OperationReply, 256), makeBinaryTreeSet(), make(map[int]bool), make(map[int]bool)}
	return &x
}
