package ActorBinaryTree

import "fmt"

type ReplyTracker map[int]bool

func (r ReplyTracker) sentOp(o Operation) {
	r[o.Id()] = true
}

func (r ReplyTracker) receivedReply(o OperationReply) error {
	if !r[o.Id()] {
		return fmt.Errorf("received reply %d that had not been sent")
	}

	r[o.Id()] = false

	return nil
}

func (r ReplyTracker) checkAllReceived() bool {
	for _, v := range r {
		if v {
			return false
		}
	}

	return true
}
