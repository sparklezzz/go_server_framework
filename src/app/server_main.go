package main

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/fcgi"
	"os"
	"reflect"
	"time"
	//"strconv"
	"common"
	"processors"
	"service"

	l4g "github.com/alecthomas/log4go"
)

func testParseConf() {
	conf_file := "conf/app.conf.done"
	var p_conf *common.Conf
	p_conf, err := common.ParseConf(conf_file)
	if err != nil {
		fmt.Println("parse conf error!")
		return
	}
	l4g.Debug("%s", p_conf.GetDebugString())
}

func testOsStat() {
	fileinfo, err := os.Stat("conf/app.conf")
	if err != nil {
		panic(err)
	}
	l4g.Debug(fileinfo.Name())    //获取文件名
	l4g.Debug(fileinfo.IsDir())   //判断是否是目录，返回bool类型
	l4g.Debug(fileinfo.ModTime()) //获取文件修改时间
	l4g.Debug(fileinfo.Mode())
	l4g.Debug(fileinfo.Size()) //获取文件大小
	l4g.Debug(fileinfo.Sys())

	t := fileinfo.ModTime().Unix()

	l4g.Debug(t)
}

func testReflect() {
	//var p_conf *common.Conf = nil
	//var reloadBase common.Reloadable = p_conf
	//t := reflect.TypeOf(reloadBase)
	conf := common.Conf{}
	t := reflect.TypeOf(conf)
	p_reflect_obj := reflect.New(t).Interface().(common.Reloadable)
	//p_reflect_obj.load("conf/app.conf.done")
	//p_reflect_obj.
	l4g.Debug("%s, %s", t, p_reflect_obj.(*common.Conf))
	//var conf common.Conf
	//common.Conf = common.Conf{}
}

func testSubStr() {
	str := "3345667"
	l4g.Debug(str[0 : len(str)-5])
}

func testReloader() {
	reloader := common.Reloader{}
	reloader.Init("conf/app.conf.done", reflect.TypeOf(common.Conf{}))
	reloader.DoLoad()
	for {
		time.Sleep(10 * time.Second)
		reloader.DoReload()
	}
}

func testEnv() {
	for {
		time.Sleep(3 * time.Second)
		common.GetEnv().Reload()
		l4g.Debug(common.GetConf().GetDebugString())
	}
}

func initEnv() {
	env := common.GetEnv()
	err := env.Init()
	if err != nil {
		panic(err)
	}
}

func initLogger() {
	//l4g.AddFilter("stdout", l4g.DEBUG, l4g.NewConsoleLogWriter())             //输出到控制台,级别为DEBUG
	//l4g.AddFilter("file", l4g.DEBUG, l4g.NewFileLogWriter("test.log", false)) //输出到文件,级别为DEBUG,文件名为test.log,每次追加该原文件
	l4g.LoadConfiguration("log4go.xml") //使用加载配置文件,类似与java的log4j.propertites
	l4g.Debug("the time is now :%s -- %s", "213", "sad")
}

func closeLogger() {
	defer l4g.Close() //注:如果不是一直运行的程序,请加上这句话,否则主线程结束后,也不会输出和log到日志文件
}

func sayhello(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "hello world!")
}

func startServer() {
	mux := http.NewServeMux()
	mux.HandleFunc("/hello", sayhello)
	http.ListenAndServe(":54320", mux)
}

func startFCGIServer() {
	//l4g.Debug(common.GetConf().GetDebugString())
	host := common.GetConf().GetStr("server_host", "")
	port := common.GetConf().GetStr("server_port", "")
	listener, _ := net.Listen("tcp", host+":"+port)
	l4g.Info("fcgi server will start at %s", host+":"+port)

	srv := new(service.FastCGIServer)
	fcgi.Serve(listener, srv)
	//l4g.Info("fcgi server started at %s", host + ":" + port)
}

func main() {
	pwd, _ := os.Getwd()
	fmt.Println(pwd)

	// init logger
	initLogger()

	processors.RegisterAllProcessors()

	// init env
	initEnv()

	//testParseConf()
	//testReflect()
	//testOsStat()
	//testSubStr()
	//testReloader()
	go func() {
		testEnv()
	}()

	//startServer()
	startFCGIServer()

	closeLogger()

}
