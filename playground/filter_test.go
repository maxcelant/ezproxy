package playground

type Request struct {
	Body    map[string]interface{}
	Headers map[string]string
}

// We can step through each filter state by using a recursive definition

type StateFn func(*Request) StateFn

type Filter func(Request) Request

func Process(r *Request) {
	foo := Filter(func(r *Request) { r.Headers["foo"] = "bar" })
	baz := Filter(func(r *Request) { r.Headers["baz"] = "wux" })

}
