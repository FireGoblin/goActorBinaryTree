Goal

The aim is to implement an actor based binary tree based on an assignment from Principles of Reactive Programming
on Coursera which was a course on Scala.  In the original assignment used the pattern matching in Scala for making
it easier to handle the operations and operationReplys, and used the akka Actor system.

The structure is a binary tree.  Each node contains an element integer and a boolean representing whether that element
has been removed.  Instead of actually removing a node, we simply flip the removed variable in it to true.  GCing 
involves copying all non-removed nodes to a new tree and then freeing up the old nodes, thus elimination all
removed nodes (besides the root 0 node which may or may not be removed).  

Porting from scala

The interfaces Operation and OperationReply were originally traits in the scala code.  Perform() was the new function
for the golang implementation of Operations.  In Scala used its pattern matching to have a large match statement in
BinaryTreeSet and BinaryTreeNode that had what to do with every Operation and Reply.  This could be done in golang, 
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

type BinaryTreeNode struct {
    // contains filtered or unexported fields
}
    max two children

func (b *BinaryTreeNode) Run()

func (b *BinaryTreeNode) String() string

type BinaryTreeSet struct {
    // contains filtered or unexported fields
}
    max two children

func (b *BinaryTreeSet) Run()

type Contains struct {
    // contains filtered or unexported fields
}

func (c Contains) Elem() int

func (c Contains) Id() int

func (c Contains) Perform(node *BinaryTreeNode)

func (c Contains) RequesterChan() chan OperationReply

func (i Contains) String() string

type ContainsResult struct {
    // contains filtered or unexported fields
}

func (c ContainsResult) Id() int

func (c ContainsResult) Result() bool

func (c ContainsResult) String() string

type CopyInsert struct {
    // contains filtered or unexported fields
}

func (i CopyInsert) Elem() int

func (i CopyInsert) Id() int

func (i CopyInsert) Perform(node *BinaryTreeNode)

func (i CopyInsert) RequesterChan() chan OperationReply

func (i CopyInsert) String() string

type GC struct{}
    dummy operation

func (g GC) Elem() int

func (g GC) Id() int

func (g GC) Perform(node *BinaryTreeNode)

func (g GC) RequesterChan() chan OperationReply

func (g GC) String() string

type GetElems struct {
    // contains filtered or unexported fields
}

func (i GetElems) Elem() int

func (i GetElems) Id() int

func (i GetElems) Perform(node *BinaryTreeNode)

func (i GetElems) RequesterChan() chan OperationReply

func (i GetElems) String() string

type Insert struct {
    // contains filtered or unexported fields
}

func (i Insert) Elem() int

func (i Insert) Id() int

func (i Insert) Perform(node *BinaryTreeNode)

func (i Insert) RequesterChan() chan OperationReply

func (i Insert) String() string

type Operation interface {
    Id() int
    Elem() int
    RequesterChan() chan OperationReply
    Perform(*BinaryTreeNode)
}

type OperationFinished struct {
    // contains filtered or unexported fields
}

func (o OperationFinished) Id() int

func (o OperationFinished) String() string

type OperationReply interface {
    Id() int
}
    note that any operation also satisfies OperationReply

type Remove struct {
    // contains filtered or unexported fields
}

func (r Remove) Elem() int

func (r Remove) Id() int

func (r Remove) Perform(node *BinaryTreeNode)

func (r Remove) RequesterChan() chan OperationReply

func (i Remove) String() string

type ReplyTracker struct {
    // contains filtered or unexported fields
}

type TestProbe struct {
    // contains filtered or unexported fields
}

func (t *TestProbe) Run(succeed chan int, fail chan int)