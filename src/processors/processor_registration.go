package processors

import (
	"common"
)

func RegisterAllProcessors() error {
	common.RegisterProcessor(common.DEFAULT_PROCESSOR, new(DefaultProcessor))
	return nil
}
