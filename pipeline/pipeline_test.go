package pipeline

import (
	"errors"
	"math"
	"testing"
)

var head ContextFn = func(v interface{}, send func(interface{}) bool, abort func(error)) {
	var double = v.(int) * 2
	send(double)
}

var middle ContextFn = func(v interface{}, send func(interface{}) bool, abort func(error)) {
	var pow = math.Pow(float64(v.(int)), 2)
	send(int(pow))
}

var middleFail ContextFn = func(v interface{}, send func(interface{}) bool, abort func(error)) {
	if v.(int) == 2 {
		abort(errors.New("test error"))
		return
	}

	send(v)
}

var tail ContextFn = func(v interface{}, send func(interface{}) bool, abort func(error)) {
	//noop, i'm lazy
}

func TestPipeline(t *testing.T) {
	p := New()
	err := p.Chain(head)
	if err != nil {
		t.Error(err)
	}

	err = p.Chain(middle)
	if err != nil {
		t.Error(err)
	}

	err = p.Chain(tail)
	if err != nil {
		t.Error(err)
	}

	for i := 0; i < 5; i++ {
		p.Start(i)
	}

	if err := p.Wait(); err != nil {
		t.Error(err)
	}
}

func TestPipelineWithControlledError(t *testing.T) {
	p := New()

	err := p.Chain(head)
	if err != nil {
		t.Error(err)
	}

	err = p.Chain(middleFail)
	if err != nil {
		t.Error(err)
	}

	err = p.Chain(tail)
	if err != nil {
		t.Error(err)
	}

	for i := 0; i < 5; i++ {
		p.Start(i)
	}

	if err := p.Wait(); err == nil {
		t.Error(errors.New("err should have occurred"))
	}

}

func TestPipelineMultiplex(t *testing.T) {
	p := New()
	err := p.Chain(head)
	if err != nil {
		t.Error(err)
	}

	if err = p.FanOut(tail, 4); err != nil {
		t.Error(err)
	}

	for i := 0; i < 10; i++ {
		p.Start(i)
	}

	if err := p.Wait(); err != nil {
		t.Error(err)
	}
}

func TestPipelineMultiplexDemultiplex(t *testing.T) {
	p := New()
	err := p.Chain(head)
	if err != nil {
		t.Error(err)
	}

	err = p.FanOut(middle, 4)
	if err != nil {
		t.Error(err)
	}

	err = p.FanIn(4, tail)
	if err != nil {
		t.Error(err)
	}

	for i := 0; i < 10; i++ {
		p.Start(i)
	}

	if err := p.Wait(); err != nil {
		t.Error(err)
	}
}

func TestErrorPipelineMultiplexDemultiplex(t *testing.T) {
	p := New()
	err := p.Chain(head)
	if err != nil {
		t.Error(err)
	}

	err = p.FanOut(middleFail, 4)
	if err != nil {
		t.Error(err)
	}

	err = p.FanIn(4, tail)
	if err != nil {
		t.Error(err)
	}

	for i := 0; i < 10; i++ {
		p.Start(i)
	}

	if err := p.Wait(); err == nil {
		t.Error(errors.New("err should have occurred"))
	}
}

func TestPipelineMtoNProcs(t *testing.T) {
	p := New()

	err := p.Chain(head)
	if err != nil {
		t.Error(err)
	}

	err = p.FanOut(middle, 4)
	if err != nil {
		t.Error(err)
	}

	err = p.ConnectNtoM(4, 4, tail)
	if err != nil {
		t.Error(err)
	}

	for i := 0; i < 10; i++ {
		p.Start(i)
	}

	if err := p.Wait(); err != nil {
		t.Error(err)
	}
}

func TestErrorPipelineMtoNProcs(t *testing.T) {
	p := New()

	err := p.Chain(head)
	if err != nil {
		t.Error(err)
	}

	err = p.FanOut(middleFail, 4)
	if err != nil {
		t.Error(err)
	}

	err = p.ConnectNtoM(4, 4, tail)
	if err != nil {
		t.Error(err)
	}

	for i := 0; i < 10; i++ {
		p.Start(i)
	}

	if err := p.Wait(); err == nil {
		t.Error(errors.New("err should have occurred"))
	}
}

func TestPipelineMtoNProcsAndFanIn(t *testing.T) {
	p := New()

	err := p.Chain(head)
	if err != nil {
		t.Error(err)
	}

	err = p.FanOut(middle, 4)
	if err != nil {
		t.Error(err)
	}

	err = p.ConnectNtoM(4, 7, middle)
	if err != nil {
		t.Error(err)
	}

	err = p.FanIn(7, tail)
	if err != nil {
		t.Error(err)
	}

	for i := 0; i < 10; i++ {
		p.Start(i)
	}

	if err := p.Wait(); err != nil {
		t.Error(err)
	}
}
