package pipeline

import (
	"errors"
	"sync"
)

//ContextFn wraps the users function and provides the necessary controls to alert the pipeline of completion and errors
type ContextFn func(interface{}, func(interface{}) bool, func(error))

//Process listens for work, processes the work via the ContextFn and finally sends the work to a subsequent process
type Process struct {
	Run ContextFn

	parent  *Pipeline
	send    chan interface{}
	receive chan interface{}
	closer  chan bool

	once sync.Once
}

func newProcess(p *Pipeline, fn ContextFn) *Process {
	var proc = &Process{
		parent: p,
		Run:    fn,
	}

	return proc
}

//Send moves data to the next process in the pipeline
func (p *Process) Send(v interface{}) (closed bool) {
	defer func() {
		if r := recover(); r != nil {
			closed = true
		}
	}()

	p.send <- v
	return false
}

//Receive accepts data and processes it
func (p *Process) Receive() {
	for v := range p.receive {
		if v != nil {
			p.Run(v, p.Send, p.Abort)
		}
	}

	p.Close()
}

//Abort calls the parent abort func
func (p *Process) Abort(err error) {
	p.parent.Abort(err)
}

//Closer adds sets closer channel
func (p *Process) Closer(c chan bool) {
	p.closer = c
}

//Close closes the underlying send channel
func (p *Process) Close() {
	p.once.Do(func() {
		if p.send != nil {
			close(p.send)
		}

		if p.closer != nil {
			p.closer <- true
		}
	})
}

//In sets the receiver channel
func (p *Process) In(in chan interface{}) error {
	if in == nil {
		return errors.New("cannot use nil sender")
	}

	p.receive = in
	return nil
}

//Out returns sender channel
func (p *Process) Out() chan interface{} {
	if p.send == nil {
		p.send = make(chan interface{})
	}

	return p.send
}
