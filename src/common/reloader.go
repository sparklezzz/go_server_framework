package common

import (
	"errors"
	"os"
	"reflect"

	l4g "github.com/alecthomas/log4go"
)

type Reloadable interface {
	load(file_name string) error
}

type Reloader struct {
	tm           int64
	using_no     int
	fname        string
	reflect_type reflect.Type
	p_content    [2]Reloadable
}

func (this *Reloader) Init(file_name string, reflect_type reflect.Type) {
	this.tm = 0
	this.fname = file_name
	this.using_no = 1
	this.reflect_type = reflect_type
	this.p_content[0] = nil
	this.p_content[1] = nil
}

func (this *Reloader) GetContent() Reloadable {
	return this.p_content[this.using_no]
}

func (this *Reloader) DoLoad() error {
	if this.using_no != 0 && this.using_no != 1 {
		return errors.New("Reload " + this.fname + " failed!")
	}

	new_tm, err := this.getFileTs(this.fname)
	if err != nil {
		return err
	}

	p_t := reflect.New(this.reflect_type).Interface().(Reloadable)
	if err := p_t.load(this.fname); err != nil {
		return err
	}

	this.tm = new_tm
	this.p_content[this.using_no] = p_t

	l4g.Info("%s load success", this.fname)
	//l4g.Debug(p_t)

	return nil
}

func (this *Reloader) DoReload() error {
	if this.using_no != 0 && this.using_no != 1 {
		return errors.New("Reload " + this.fname + " failed!")
	}

	new_tm, err := this.getFileTs(this.fname)
	if err != nil {
		return err
	}

	if new_tm <= this.tm { // old file do not need to reload
		l4g.Debug("%s: got tm %d not newer than recorded tm %d", this.fname, new_tm, this.tm)
		return nil
	}

	// fill new buffer
	p_t := reflect.New(this.reflect_type).Interface().(Reloadable)
	if err := p_t.load(this.fname); err != nil {
		return err
	}

	reload_no := 1 - this.using_no
	this.tm = new_tm
	this.p_content[reload_no] = p_t

	// switch ptr to new buffer
	this.using_no = reload_no

	l4g.Info("%s reload success", this.fname)

	// release old buffer
	previous_no := 1 - this.using_no
	this.p_content[previous_no] = nil

	return nil
}

func (this *Reloader) NeedReload() (int, error) {
	new_tm, err := this.getFileTs(this.fname)
	if err != nil {
		return 0, err
	}

	if new_tm <= this.tm { // old file do not need to reload
		return 0, nil
	}

	return 1, nil
}

func (this *Reloader) getFileTs(file_name string) (int64, error) {
	fileinfo, err := os.Stat(file_name)
	if err != nil {
		return -1, err
	}
	return fileinfo.ModTime().Unix(), nil
}
