package service

import (
	//"net"
	"net/http"
	//"net/http/fcgi"
	"common"

	l4g "github.com/alecthomas/log4go"
)

type FastCGIServer struct{}

func (s FastCGIServer) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	chain_id := req.URL.Query().Get("chain_id")
	if chain_id != "" {
		l4g.Info("got chain id %s", chain_id)
		p_chain := common.GetChainConf().GetChain(chain_id)
		if nil != p_chain {
			p_context := new(common.Context)
			p_result := new(common.Result)
			p_chain.Run(p_context, p_result)
			resp.Write([]byte("found chain_id: " + chain_id + " and execute chain successfully"))
		} else {
			resp.Write([]byte("cannot found chain_id: " + chain_id))
		}
	} else {
		resp.Write([]byte("<h1>Hello, 世界</h1>\n<p>Behold my Go web app.</p>"))
	}
}
