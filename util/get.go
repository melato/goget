package util

// Get provides a goroutine-safe way to call a function that takes an argument and returns a value.
// It uses a background goroutine and a channel to process the requests.
type Get[A, V any] struct {
	loadFunc func(A) V
	queue    chan *getRequest[A, V]
}

func NewGet[A, V any](loadFunc func(A) V) *Get[A, V] {
	var t Get[A, V]
	t.loadFunc = loadFunc
	t.queue = make(chan *getRequest[A, V])
	go t.reader()
	return &t
}

type getRequest[A, V any] struct {
	arg   A
	Value V
	done  chan struct{}
}

// Get returns the result of calling the function with the given argument.
func (t *Get[A, V]) newReadRequest(arg A) *getRequest[A, V] {
	r := &getRequest[A, V]{arg: arg}
	r.done = make(chan struct{})
	t.queue <- r
	return r
}

func (r *getRequest[A, V]) Close(value V) {
	r.Value = value
	close(r.done)
}

func (r *getRequest[A, V]) Wait() {
	_ = <-r.done
}

func (t *Get[A, V]) reader() {
	for r := range t.queue {
		r.Close(t.loadFunc(r.arg))
	}
}

func (t *Get[A, V]) Get(arg A) V {
	read := t.newReadRequest(arg)
	read.Wait()
	return read.Value
}
