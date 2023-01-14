package framework

type IGroup interface {
	Get(string, ControllerHandler)
	Post(string, ControllerHandler)
	Put(string, ControllerHandler)
	Delete(string, ControllerHandler)
}

type Group struct {
	core   *Core
	prefix string
}

func NewGroup(core *Core, prefix string) *Group {
	return &Group{
		core:   core,
		prefix: prefix,
	}
}

func (this *Group) Get(url string, handler ControllerHandler) {
	this.core.Get(this.prefix+url, handler)
}

func (this *Group) Post(url string, handler ControllerHandler) {
	this.core.Post(this.prefix+url, handler)
}

func (this *Group) Put(url string, handler ControllerHandler) {
	this.core.Put(this.prefix+url, handler)
}

func (this *Group) Delete(url string, handler ControllerHandler) {
	this.core.Delete(this.prefix+url, handler)
}
