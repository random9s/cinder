package pipeline

import (
	"errors"
	"fmt"
	"runtime"
	"sync"
	"time"
)

/*Pipeline contains a group of chained processes and waits for completion or an error to occur
 *
 * Error Handling:
 *  There is an error channel wrapped in a func called Abort that all processes can access.
 *  Calling the Abort fn causes the Pipeline to Close all running processes, decrement the wg counter
 *  and then return the error
 *
 * Concurrency Patterns Implemented:
 *  Fan In: Using the fan in method allows the user to connect one receiver channel to N sender channels
 *  Fan Out: Using the fan out method allows the user to connect one send channel to N receiver channels
 *  Append: Simply connects one sender to one receiver.
 *  Connect N to M: Create method to bridge N senders to N receivers
 *
 */
type Pipeline struct {
	Procs  []Processor
	tail   []Processor //the reason this is a slice is because of the fan out func, out tail could be multiple procs
	closer chan bool

	once sync.Once

	//Err Handling//
	abort       chan error
	abortedWith chan error
	globalAlert chan int
	started     bool
}

//New returns a newly initialized pipeline
func New() *Pipeline {
	var pipe = &Pipeline{
		Procs:       make([]Processor, 0),
		tail:        make([]Processor, 0),
		closer:      make(chan bool, 1),
		abort:       make(chan error, 1),
		abortedWith: make(chan error, 1),
		globalAlert: make(chan int),
		started:     false,
	}
	pipe.Procs = append(pipe.Procs, newProcess(pipe, nil))

	go pipe.watchForErrors()
	return pipe
}

//Start sends data to the first process in the pipeline
func (p *Pipeline) Start(v interface{}) {
	if len(p.Procs) == 0 {
		panic(errors.New("cannot send on nil channel"))
	}

	//Before the first send, attach the closer channel to all final processes
	p.once.Do(func() {
		p.started = true
		//upgrade the buffer size, if needed
		p.closer = make(chan bool, len(p.tail))
		for _, tail := range p.tail {
			tail.Closer(p.closer)
		}
	})

	p.Procs[0].Send(v)
}

//Wait starts signaling to close channels after all items have been sent, waits for response from closer
func (p *Pipeline) Wait() error {
	if p.closer == nil {
		panic(errors.New("cannot receive from nil channel"))
	}

	//wait for all procs in tail to close
	p.Procs[0].Close()
	if !p.started {
		return nil
	}

	for i := 0; i < len(p.tail); i++ {
		<-p.closer
	}

	return p.err()
}

//WaitWithTimeout starts signaling to close channels after all items have been sent, waits for response from closer or timesout
func (p *Pipeline) WaitWithTimeout(dur time.Duration) error {
	if p.closer == nil {
		panic(errors.New("cannot receive from nil channel"))
	}

	p.Procs[0].Close()
	if !p.started {
		return nil
	}

	//wait for all procs in tail to close or timeout to occur
	for i := 0; i < len(p.tail); i++ {
		select {
		case <-p.closer:
		case <-time.After(dur):
			p.Abort(errors.New("timeout: process could not be completed"))
		}
	}

	return p.err()
}

//Abort sends an error over the abort channel, which will cause a shutdown
func (p *Pipeline) Abort(err error) {
	defer func() {
		recover()
	}()

	if p.abort != nil {
		fmt.Printf("ABORTED WITH: %s", err)
		p.abort <- err
	}
}

//The pipeline reads from the abort channel until an error occurs
//It's up to the user to send errors via the abort func
func (p *Pipeline) watchForErrors() {
	//This will block until an error occurs
	p.abortedWith <- <-p.abort
	close(p.globalAlert)
	fmt.Println("SHUTDOWN")
	p.shutdown()
}

//shutdown closes all running processes in the pipeline
func (p *Pipeline) shutdown() {
	for _, proc := range p.Procs {
		proc.Close()
	}
}

//Error returns an error, if one occured
func (p *Pipeline) err() error {
	var err error

	select {
	case err = <-p.abortedWith:
	default:
	}

	return err
}

//Chain appends a process to our current chain of processes
func (p *Pipeline) Chain(fn ContextFn) error {
	if len(p.Procs) == 0 {
		return errors.New("cannot chain from nil process")
	}

	proc := newProcess(p, fn)

	if len(p.Procs) > 0 {
		//Connect the last process to the current process
		if err := Connect(p.Procs[len(p.Procs)-1], proc); err != nil {
			return err
		}
	}

	p.tail = []Processor{proc}
	p.Procs = append(p.Procs, proc)
	return nil
}

//FanOut connects procs to the last process
func (p *Pipeline) FanOut(fn ContextFn, n int) error {
	if len(p.Procs) == 0 {
		return errors.New("cannot fan out from nil process")
	}

	var procs = make([]Processor, 0)
	for i := 0; i < n; i++ {
		procs = append(procs, newProcess(p, fn))
	}

	var sendProc = p.Procs[len(p.Procs)-1]
	for _, recProc := range procs {
		if err := Connect(sendProc, recProc); err != nil {
			return err
		}
	}

	p.tail = procs
	p.Procs = append(p.Procs, procs...)
	return nil
}

//FanIn connects process to the last n procs
func (p *Pipeline) FanIn(n int, fn ContextFn) error {
	if len(p.Procs) == 0 {
		return errors.New("cannot fan out from nil process")
	}

	if len(p.Procs) < n {
		return fmt.Errorf("cannot fan-in from previous %d procs: only %d procs exist", n, len(p.Procs))
	}

	proc := newProcess(p, fn)

	//intermediary channel that replaces a processes receiver
	var midCh = make(chan interface{})
	err := proc.In(midCh)
	if err != nil {
		return err
	}

	go proc.Receive()

	var sendChs = make([]chan interface{}, 0)
	var sendProcs = p.Procs[len(p.Procs)-n:]
	for _, sendProc := range sendProcs {
		sendChs = append(sendChs, sendProc.Out())
	}

	//Fan in data
	go func() {
		var wg = sync.WaitGroup{}
		for _, sendCh := range sendChs {
			wg.Add(1)
			go func(ch chan interface{}, g *sync.WaitGroup) {
				for {
					select {
					case v, ok := <-ch:
						if !ok {
							wg.Done()
							return
						}

						midCh <- v
					case <-p.globalAlert:
						wg.Done()
						return
					}

					runtime.Gosched()
				}
			}(sendCh, &wg)
		}

		wg.Wait()
		close(midCh)
	}()

	p.tail = []Processor{proc}
	p.Procs = append(p.Procs, proc)
	return nil
}

//ConnectNtoM ...
func (p *Pipeline) ConnectNtoM(n, m int, fn ContextFn) error {
	if len(p.Procs) == 0 {
		return errors.New("cannot fan out from nil process")
	}

	if len(p.Procs) < n {
		return fmt.Errorf("cannot fan-in from previous %d procs: only %d procs exist", n, len(p.Procs))
	}

	var procs = make([]Processor, 0)
	for i := 0; i < m; i++ {
		procs = append(procs, newProcess(p, fn))
	}

	//intermediary channel that replaces a processes receiver
	var midCh = make(chan interface{})

	for _, proc := range procs {
		err := proc.In(midCh)
		if err != nil {
			return err
		}

		go proc.Receive()
	}

	var sendChs = make([]chan interface{}, 0)
	var sendProcs = p.Procs[len(p.Procs)-n:]
	for _, sendProc := range sendProcs {
		sendChs = append(sendChs, sendProc.Out())
	}

	//Fan in data
	go func() {
		var wg = sync.WaitGroup{}
		for _, sendCh := range sendChs {
			wg.Add(1)
			go func(ch chan interface{}, g *sync.WaitGroup) {
				for {
					select {
					case v, ok := <-ch:
						if !ok {
							wg.Done()
							return
						}

						midCh <- v
					case <-p.globalAlert:
						wg.Done()
						return
					}

					runtime.Gosched()
				}
			}(sendCh, &wg)
		}

		wg.Wait()
		close(midCh)
	}()

	p.tail = procs
	p.Procs = append(p.Procs, procs...)
	return nil
}
