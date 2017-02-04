package common

import (
	"container/list"
	"encoding/json"
	"errors"
	"fmt"
	"utils"

	l4g "github.com/alecthomas/log4go"
)

type ChainConf struct {
	chainMap map[string]*Chain
}

func (this *ChainConf) GetChain(name string) *Chain {
	val, exists := this.chainMap[name]
	if exists {
		return val
	}
	return nil
}

func (this *ChainConf) load(file_name string) error {

	// remove .done suffix
	if len(file_name) < 5 {
		return errors.New("chain conf file name too short")
	}

	file_name = file_name[0 : len(file_name)-5]

	buffer, err := utils.ReadFromJson(file_name)
	if err != nil {
		l4g.Error("load %s error: %T", file_name, err)
		return err
	}

	this.chainMap = make(map[string]*Chain)

	var dat map[string]interface{}
	if err := json.Unmarshal(buffer, &dat); err != nil {
		l4g.Error("parse json error for file %s: %T", file_name, err)
		return err
	}

	for key, val := range dat {
		l4g.Debug("Key: %s, Value: %s", key, val)
		chain_name := key

		switch val.(type) {
		case []interface{}:
			{
			} //pass
		default:
			{
				error_msg := fmt.Sprintf("val is not an array for chain name: %s", chain_name)
				return errors.New(error_msg)
			}
		}

		var process_name_arr []interface{} = val.([]interface{})

		p_processors := list.New()
		for _, elem := range process_name_arr {
			l4g.Debug("elem: %s", elem)
			switch elem.(type) {
			case string:
				{
					process_name := elem.(string)
					one_processor := GetProcessor(process_name)
					if nil == one_processor {
						return errors.New("cannot find process name: " + process_name)
					}
					p_processors.PushBack(one_processor)
				}
			default:
				{
					error_msg := fmt.Sprintf("array element is not a string for chain name: %s", chain_name)
					return errors.New(error_msg)
				}
			}
		}

		p_chain := new(Chain)
		if err := p_chain.Init(p_processors); err != nil {
			l4g.Error("chain init failed to chain name %s, %T", chain_name, err)
			return err
		}

		this.chainMap[key] = p_chain
	}

	return nil
}
