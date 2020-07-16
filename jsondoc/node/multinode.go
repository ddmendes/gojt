package node

import (
	"errors"
	"runtime"
	"sync"
)

// MultiNode represents multiple json components wrapped.
type MultiNode []Node
type getMessage struct {
	node *Node
	i    int
}

// ToMultiNode wraps an []interface{} with a MultiNode component.
func ToMultiNode(iArr []interface{}) Node {
	singleNodes := make([]Node, len(iArr))
	for i, v := range iArr {
		singleNodes[i] = SingleNode{Elem: v}
	}

	return MultiNode(singleNodes)
}

// Get a child Node for each underlying Node.
func (n MultiNode) Get(token string) (Node, error) {
	return n.sequentialGet(token)
}

func (n MultiNode) sequentialGet(token string) (Node, error) {
	var err error
	next := make([]Node, len(n))
	for i, v := range n {
		next[i], err = v.Get(token)
		if err != nil {
			return nil, err
		}
	}
	return MultiNode(next), nil
}

func setupGetWorkers(token string, nItems int, nWorkers int) (chan<- getMessage, <-chan getMessage) {
	var wg sync.WaitGroup
	wg.Add(nItems)

	input := make(chan getMessage, nItems)
	output := make(chan getMessage, nItems)

	worker := func(token string, in <-chan getMessage, out chan<- getMessage, wg *sync.WaitGroup) {
		for m := range in {
			next, _ := (*m.node).Get(token)
			out <- getMessage{&next, m.i}
			wg.Done()
		}
	}

	for i := 0; i < nWorkers; i++ {
		go worker(token, input, output, &wg)
	}

	go func() {
		wg.Wait()
		close(output)
	}()

	return input, output
}

// concurrentGet executes get for each item of underlying array concurrently.
// CAUTION! This is experimental and based on benchmarks is really slower
// than sequential version due to array caching.
func (n MultiNode) concurrentGet(token string) (Node, error) {
	output := make([]Node, len(n))
	send, receive := setupGetWorkers(token, len(n), runtime.NumCPU())
	go func() {
		for i, node := range n {
			send <- getMessage{&node, i}
		}
		close(send)
	}()
	for next := range receive {
		output[next.i] = *next.node
	}
	return MultiNode(output), nil
}

// GetKeys returns all the keys available for this Node.
func (n MultiNode) GetKeys() ([]string, error) {
	if len(n) > 0 {
		return n[0].GetKeys()
	}
	return []string{}, errors.New("Document is empty")
}

// GetInterface gets the interface{} value of this Node component
func (n MultiNode) GetInterface() interface{} {
	interfaces := make([]interface{}, len(n), len(n))
	for i, v := range n {
		interfaces[i] = v.GetInterface()
	}
	return interface{}(interfaces)
}
