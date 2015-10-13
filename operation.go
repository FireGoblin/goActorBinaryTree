package ActorBinaryTree

type Operation interface {
	ID() int
	Elem() int
	RequesterChan() chan OperationReply
	Perform(*binaryTreeNode)
}
