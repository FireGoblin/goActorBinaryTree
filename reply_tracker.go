package ActorBinaryTree

import "fmt"
import "sync"

type replyTracker struct {
	m   map[int]bool
	key *sync.Mutex
}

func (r replyTracker) sentOp(o operation) {
	r.key.Lock()
	r.m[o.ID()] = true
	r.key.Unlock()
}

func (r replyTracker) receivedReply(o operationReply) error {
	r.key.Lock()
	if !r.m[o.ID()] {
		r.key.Unlock()
		return fmt.Errorf("received reply %t that had not been sent", r.m[o.ID()])
	}

	r.m[o.ID()] = false

	r.key.Unlock()
	return nil
}

func (r replyTracker) get(i int) bool {
	r.key.Lock()
	x := r.m[i]
	r.key.Unlock()

	return x
}

func (r replyTracker) checkAllReceived() bool {
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

func (r replyTracker) displayUnreceived() {
	r.key.Lock()
	for k, v := range r.m {
		if v {
			fmt.Println(k)
		}
	}
	r.key.Unlock()
}
