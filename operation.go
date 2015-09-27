package ActorBinaryTree

type Operation interface {
	Id() int
	Elem() int
	RequesterChan() chan OperationReply
	Perform(*BinaryTreeNode)
}
