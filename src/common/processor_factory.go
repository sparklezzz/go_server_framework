package common

var (
	processorFactorySingleton ProcessorFactory
)

type ProcessorFactory struct {
	processorMap map[string]Processor
}

func init() {
	processorFactorySingleton = ProcessorFactory{}
	processorFactorySingleton.processorMap = make(map[string]Processor)
}

func (this *ProcessorFactory) getProcessor(name string) Processor {
	val, exists := this.processorMap[name]
	if exists {
		return val
	}
	return nil
}

func (this *ProcessorFactory) registerProcessor(name string, processor Processor) {
	this.processorMap[name] = processor
}

func GetProcessor(name string) Processor {
	return processorFactorySingleton.getProcessor(name)
}

func RegisterProcessor(name string, processor Processor) {
	processorFactorySingleton.registerProcessor(name, processor)
}
