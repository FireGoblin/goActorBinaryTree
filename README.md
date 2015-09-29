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