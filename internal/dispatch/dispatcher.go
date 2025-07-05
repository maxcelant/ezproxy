package dispatch

import "net"

type DispatchContext struct {
	Upstreams []string
	Conn      net.Conn
}

type Dispatcher struct {
	ctx          DispatchContext
	dispatchFunc func(DispatchContext)
}

func NewDispatcher() *Dispatcher {
	return &Dispatcher{
		dispatchFunc: func(_ DispatchContext) {},
		ctx: DispatchContext{
			Upstreams: make([]string, 0),
		},
	}
}

func (d *Dispatcher) Mount(f func(DispatchContext)) *Dispatcher {
	d.dispatchFunc = f
	return d
}

func (d *Dispatcher) AddUpstreams(upstreams []string) {
	d.ctx.Upstreams = append(d.ctx.Upstreams, upstreams...)
}

func (d *Dispatcher) Dispatch(conn net.Conn) {
	d.ctx.Conn = conn
	d.dispatchFunc(d.ctx)
}
