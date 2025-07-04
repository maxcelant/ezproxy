package dispatch

import "net"

type Dispatcher struct {
	upstreams    []string
	dispatchFunc func(c net.Conn)
}

func NewDispatcher() *Dispatcher {
	dispatchFunc := func(c net.Conn) {}
	return &Dispatcher{
		upstreams:    make([]string, 0),
		dispatchFunc: dispatchFunc,
	}
}

func (d *Dispatcher) AddUpstreams(upstreams []string) {
	d.upstreams = append(d.upstreams, upstreams...)
}

func (d *Dispatcher) Mount(f func(c net.Conn)) {
	d.dispatchFunc = f
}

func (d *Dispatcher) Dispatch(c net.Conn) {
	d.dispatchFunc(c)
}
