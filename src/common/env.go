package common

import (
	"reflect"
)

func GetEnv() *Environment {
	return p_env
}

func GetConf() *Conf {
	return p_env.getConf()
}

func GetChainConf() *ChainConf {
	return p_env.getChainConf()
}

/*
* internals
 */

var (
	p_env *Environment
)

func init() {
	p_env = new(Environment)
}

type Environment struct {
	reloaderMap map[string]*Reloader
}

func (this *Environment) Init() error {
	err := this.initReloaderMap()
	if err != nil {
		return err
	}
	return nil
}

func (this *Environment) getConf() *Conf {
	r, exists := this.reloaderMap[CONF_RELOADER]
	if exists {
		return r.GetContent().(*Conf)
	}
	return nil
}

func (this *Environment) getChainConf() *ChainConf {
	r, exists := this.reloaderMap[CHAIN_CONF_RELOADER]
	if exists {
		return r.GetContent().(*ChainConf)
	}
	return nil
}

func (this *Environment) initReloaderMap() error {
	this.reloaderMap = make(map[string]*Reloader)

	// register all reloaders
	this.reloaderMap[CONF_RELOADER] = new(Reloader)
	this.reloaderMap[CONF_RELOADER].Init("conf/app.conf.done", reflect.TypeOf(Conf{}))

	this.reloaderMap[CHAIN_CONF_RELOADER] = new(Reloader)
	this.reloaderMap[CHAIN_CONF_RELOADER].Init("conf/chain.conf.done", reflect.TypeOf(ChainConf{}))

	for _, p_reloader := range this.reloaderMap {
		err := p_reloader.DoLoad()
		if err != nil {
			return err
		}
	}

	return nil
}

func (this *Environment) Reload() error {
	for _, p_reloader := range this.reloaderMap {
		err := p_reloader.DoReload()
		if err != nil { // do not return when on reloader failed
			continue
			//return err
		}
	}

	return nil
}
