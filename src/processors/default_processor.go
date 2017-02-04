package processors

import (
	"common"
)

func init() {

}

type DefaultProcessor struct {
}

func (this *DefaultProcessor) GetName() string {
	return "DEFAULT_PROCESSOR"
}

func (this *DefaultProcessor) ForwardProcess(context *common.Context, result *common.Result) error { // do nothing
	return nil
}

func (this *DefaultProcessor) BackwardProcess(context *common.Context, result *common.Result) error { // do nothing
	return nil
}
