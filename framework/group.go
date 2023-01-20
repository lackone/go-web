package framework

type IGroup interface {
	Get(string, ...ControllerHandler)
	Post(string, ...ControllerHandler)
	Put(string, ...ControllerHandler)
	Delete(string, ...ControllerHandler)
	Use(...ControllerHandler)
	Group(string) IGroup
}

type Group struct {
	core        *Core
	parent      *Group
	prefix      string
	middlewares []ControllerHandler
}

func NewGroup(core *Core, prefix string) *Group {
	return &Group{
		core:        core,
		parent:      nil,
		prefix:      prefix,
		middlewares: []ControllerHandler{},
	}
}

func (this *Group) GetAbsPrefix() string {
	if this.parent == nil {
		return this.prefix
	}
	return this.parent.GetAbsPrefix() + this.prefix
}

func (this *Group) GetMiddlewares() []ControllerHandler {
	if this.parent == nil {
		return this.middlewares
	}
	return append(this.parent.GetMiddlewares(), this.middlewares...)
}

func (this *Group) Group(uri string) IGroup {
	group := NewGroup(this.core, uri)
	group.parent = this
	return group
}

func (this *Group) Get(uri string, handlers ...ControllerHandler) {
	allHandlers := append(this.GetMiddlewares(), handlers...)
	this.core.Get(this.GetAbsPrefix()+uri, allHandlers...)
}

func (this *Group) Post(uri string, handlers ...ControllerHandler) {
	allHandlers := append(this.GetMiddlewares(), handlers...)
	this.core.Post(this.GetAbsPrefix()+uri, allHandlers...)
}

func (this *Group) Put(uri string, handlers ...ControllerHandler) {
	allHandlers := append(this.GetMiddlewares(), handlers...)
	this.core.Put(this.GetAbsPrefix()+uri, allHandlers...)
}

func (this *Group) Delete(uri string, handlers ...ControllerHandler) {
	allHandlers := append(this.GetMiddlewares(), handlers...)
	this.core.Delete(this.GetAbsPrefix()+uri, allHandlers...)
}

func (this *Group) Use(middlewares ...ControllerHandler) {
	this.middlewares = append(this.middlewares, middlewares...)
}
