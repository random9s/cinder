package pipeline

//Processor accepts data, does whatever, and sends data out
type Processor interface {
	Sender
	Receiver
	Closer(chan bool)
	Close()
}

//Sender implements a send channel
type Sender interface {
	Send(interface{}) bool
	Out() chan interface{}
}

//Receiver implements a receive channel
type Receiver interface {
	Receive()
	In(chan interface{}) error
}

//Connect bridges two processes via channels
func Connect(s Sender, r Receiver) error {
	err := r.In(s.Out())
	if err != nil {
		return err
	}

	go r.Receive()
	return nil
}
