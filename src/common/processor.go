package common

type Processor interface {
	GetName() string
	ForwardProcess(context *Context, result *Result) error
	BackwardProcess(context *Context, result *Result) error
}
