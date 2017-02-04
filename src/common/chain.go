package common

import (
	"container/list"

	l4g "github.com/alecthomas/log4go"
)

type Chain struct {
	p_processors *list.List
}

func (this *Chain) Init(processors *list.List) error {
	this.SetProcessors(processors)
	return nil
}

func (this *Chain) SetProcessors(processors *list.List) {
	this.p_processors = processors
}

func (this *Chain) Run(context *Context, result *Result) error {
	if this.p_processors == nil || this.p_processors.Len() == 0 {
		return nil
	}

	for e := this.p_processors.Front(); e != nil; e = e.Next() {
		processor := e.Value.(Processor)
		l4g.Debug("forward in process: " + processor.GetName())
		if err := processor.ForwardProcess(context, result); err != nil {
			l4g.Error("%s forward processing error: %T", processor.GetName(), err)
			return err
		}
	}

	for e := this.p_processors.Back(); e != nil; e = e.Prev() {
		processor := e.Value.(Processor)
		l4g.Debug("backward in process: " + processor.GetName())

		if err := processor.BackwardProcess(context, result); err != nil {
			l4g.Error("%s backward processing error: %T", processor.GetName(), err)
			return err
		}
	}

	return nil
}
