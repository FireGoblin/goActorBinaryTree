package ActorBinaryTree

type operation interface {
	ID() int
	Elem() int
	RequesterChan() chan operationReply
	Perform(*binaryTreeNode)
}
