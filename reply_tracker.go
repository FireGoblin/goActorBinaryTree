package ActorBinaryTree

import "fmt"
import "sync"

type ReplyTracker struct {
	m   map[int]bool
	key *sync.Mutex
}

func (r ReplyTracker) sentOp(o Operation) {
	r.key.Lock()
	r.m[o.ID()] = true
	r.key.Unlock()
}

func (r ReplyTracker) receivedReply(o OperationReply) error {
	r.key.Lock()
	if !r.m[o.ID()] {
		r.key.Unlock()
		return fmt.Errorf("received reply %t that had not been sent", r.m[o.ID()])
	}

	r.m[o.ID()] = false

	r.key.Unlock()
	return nil
}

func (r ReplyTracker) get(i int) bool {
	r.key.Lock()
	x := r.m[i]
	r.key.Unlock()

	return x
}

func (r ReplyTracker) checkAllReceived() bool {
	r.key.Lock()
	for _, v := range r.m {
		if v {
			r.key.Unlock()
			return false
		}
	}
	r.key.Unlock()

	return true
}

func (r ReplyTracker) displayUnreceived() {
	r.key.Lock()
	for k, v := range r.m {
		if v {
			fmt.Println(k)
		}
	}
	r.key.Unlock()
}
