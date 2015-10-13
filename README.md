Goal

The aim is to implement an actor based binary tree based on an assignment from Principles of Reactive Programming
on Coursera which was a course on Scala.  In the original assignment used the pattern matching in Scala for making
it easier to handle the operations and operationReplys, and used the akka Actor system.

The structure is a binary tree.  Each node contains an element integer and a boolean representing whether that element
has been removed.  Instead of actually removing a node, we simply flip the removed variable in it to true.  gcing 
involves copying all non-removed nodes to a new tree and then freeing up the old nodes, thus elimination all
removed nodes (besides the root 0 node which may or may not be removed).  

Porting from scala

The interfaces Operation and OperationReply were originally traits in the scala code.  Perform() was the new function
for the golang implementation of Operations.  In Scala used its pattern matching to have a large match statement in
BinaryTreeSet and binaryTreeNode that had what to do with every Operation and Reply.  This could be done in golang, 
but it seemed better to have the function for performing the operation actually in the file with the operation than
forcing binaryTreeNode to contain the code for every operation in it.  It also makes for much smaller and more 
manageable Run commands for nodes, and allows for the addition of new operations without changing Run() commands.

Random finding

When I originally made the run routines I included a default: case in them.  This was a big mistake as it caused an
enormous performance issue that made it take seconds to perform trivial tasks.  Not only is it unnecessary in a loop
where it can't do anything until receives something on a channel, but apparently it causes enormous performance issues.

Dealing with race conditions

There were two race conditions ran into.  One of them lead to adding the mutex to reply_tracker to coordinate map accesses.
The other lead to using injectOperation() in tests instead of sendOperation() so all actions performed by the testProbe
go through the Run() routine instead of concurrently to it.

TODO

For the binary_tree_set to be usable as a data type in another project need to move some of the test_probe code into the
binary_tree_set to give it a convenient public api instead of making them form operations and send them correctly like
test_probe.

//-------------------------------

PACKAGE DOCUMENTATION

package ActorBinaryTree
    import "."


TYPES

type binaryTreeNode struct {
    // contains filtered or unexported fields
}
    max two children

func (b *binaryTreeNode) Run()

func (b *binaryTreeNode) String() string

type BinaryTreeSet struct {
    // contains filtered or unexported fields
}
    max two children

func (b *BinaryTreeSet) Run()

type contains struct {
    // contains filtered or unexported fields
}

func (c contains) Elem() int

func (c contains) ID() int

func (c contains) Perform(node *binaryTreeNode)

func (c contains) RequesterChan() chan OperationReply

func (i contains) String() string

type containsResult struct {
    // contains filtered or unexported fields
}

func (c containsResult) ID() int

func (c containsResult) Result() bool

func (c containsResult) String() string

type Copyinsert struct {
    // contains filtered or unexported fields
}

func (i Copyinsert) Elem() int

func (i Copyinsert) ID() int

func (i Copyinsert) Perform(node *binaryTreeNode)

func (i Copyinsert) RequesterChan() chan OperationReply

func (i Copyinsert) String() string

type gc struct{}
    dummy operation

func (g gc) Elem() int

func (g gc) ID() int

func (g gc) Perform(node *binaryTreeNode)

func (g gc) RequesterChan() chan OperationReply

func (g gc) String() string

type getElems struct {
    // contains filtered or unexported fields
}

func (i getElems) Elem() int

func (i getElems) ID() int

func (i getElems) Perform(node *binaryTreeNode)

func (i getElems) RequesterChan() chan OperationReply

func (i getElems) String() string

type insert struct {
    // contains filtered or unexported fields
}

func (i insert) Elem() int

func (i insert) ID() int

func (i insert) Perform(node *binaryTreeNode)

func (i insert) RequesterChan() chan OperationReply

func (i insert) String() string

type Operation interface {
    ID() int
    Elem() int
    RequesterChan() chan OperationReply
    Perform(*binaryTreeNode)
}

type operationFinished struct {
    // contains filtered or unexported fields
}

func (o operationFinished) ID() int

func (o operationFinished) String() string

type OperationReply interface {
    ID() int
}
    note that any operation also satisfies OperationReply

type remove struct {
    // contains filtered or unexported fields
}

func (r remove) Elem() int

func (r remove) ID() int

func (r remove) Perform(node *binaryTreeNode)

func (r remove) RequesterChan() chan OperationReply

func (i remove) String() string

type ReplyTracker struct {
    // contains filtered or unexported fields
}

type TestProbe struct {
    // contains filtered or unexported fields
}

func (t *TestProbe) Run(succeed chan int, fail chan int)