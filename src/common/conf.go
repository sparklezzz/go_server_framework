package common

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"utils"

	l4g "github.com/alecthomas/log4go"
)

type Conf struct {
	int_map   map[string]int
	str_map   map[string]string
	float_map map[string]float64
}

func (this *Conf) GetInt(key string, default_val int) int {
	value, exists := this.int_map[key]
	if exists {
		return value
	}
	return default_val
}

func (this *Conf) GetStr(key string, default_val string) string {
	value, exists := this.str_map[key]
	if exists {
		return value
	}
	return default_val
}

func (this *Conf) GetFloat(key string, default_val float64) float64 {
	value, exists := this.float_map[key]
	if exists {
		return value
	}
	return default_val
}

func (this *Conf) GetDebugString() string {
	s := ""
	s += "[float map]:"
	for key, val := range this.float_map {
		s += key + "#" + strconv.FormatFloat(val, 'g', 5, 64) + " "
	}
	s += "[int map]:"
	for key, val := range this.int_map {
		s += key + "#" + strconv.Itoa(val) + " "
	}
	s += "[str map]:"
	for key, val := range this.str_map {
		s += key + "#" + val + " "
	}
	return s
}

func (this *Conf) load(file_name string) error {
	// remove .done suffix
	if len(file_name) < 5 {
		return errors.New("conf file name too short")
	}

	file_name = file_name[0 : len(file_name)-5]

	buffer, err := utils.ReadFromJson(file_name)
	if err != nil {
		l4g.Error("load %s error: %T", file_name, err)
		return err
	}

	var dat map[string]interface{}
	if err := json.Unmarshal(buffer, &dat); err != nil {
		l4g.Error("parse json error for file %s: %T", file_name, err)
		return err
	}

	this.str_map = make(map[string]string)
	this.int_map = make(map[string]int)
	this.float_map = make(map[string]float64)
	fmt.Println(this)
	for key, val := range dat {
		l4g.Debug("Key: %s, Value: %s", key, val)
		switch val.(type) {
		case string: /*fmt.Println(val.(string))*/ this.str_map[key] = val.(string)
		case int:
			this.int_map[key] = val.(int)
		case float64:
			this.float_map[key] = val.(float64) // warning: std json lib will parse
			// all digit value into float 64, even if the
			// value is an int
		default:
		}
	}

	return nil
}

func ParseConf(conf_file string) (*Conf, error) {
	var p_conf *Conf
	p_conf = new(Conf)
	err := p_conf.load(conf_file)
	if err != nil {
		return nil, err
	}
	return p_conf, nil
}

/*
func ParseConf(conf_file string) (*Conf, error) {
    buffer, err := utils.ReadFromJson(conf_file)
    if (err != nil) {
        l4g.Error("load %s error: %T", conf_file, err)
        return nil, err
    }

    var dat map[string]interface{}
    if err := json.Unmarshal(buffer, &dat); err != nil {
        l4g.Error("parse json error for file %s: %T", conf_file, err)
        return nil, err
    }

    var p_conf *Conf
    p_conf = new(Conf)
    p_conf.str_map = make(map[string]string)
    p_conf.int_map = make(map[string]int)
    p_conf.float_map = make(map[string]float64)
    fmt.Println(p_conf)
    for key, val := range dat {
        l4g.Debug("Key: %s, Value: %s", key, val)
        switch val. (type) {
                case string: p_conf.str_map[key] = val.(string)
                case int: p_conf.int_map[key] = val.(int)
                case float64: p_conf.float_map[key] = val.(float64)
                default:
        }
    }

    return p_conf, nil
}
*/
