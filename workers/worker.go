package workers

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"os"
	"os/exec"
	"sync"
)

var OUT = os.Stdout
var ERR = os.Stderr

type Worker struct {
	addr string
	cmd  *exec.Cmd
	wg   *sync.WaitGroup
	rp   *httputil.ReverseProxy
}

func New(port int, out, err io.Writer) *Worker {
	return &Worker{
		wg:   &sync.WaitGroup{},
		rp:   newReverseProxy(port),
		addr: fmt.Sprintf(":%v", port),
		cmd:  createCMD(port, out, err),
	}
}

func newReverseProxy(port int) *httputil.ReverseProxy {
	director := func(req *http.Request) {
		req.URL.Scheme = "http"
		req.URL.Host = fmt.Sprintf(":%v", port)

		if _, ok := req.Header["User-Agent"]; !ok {
			// explicitly disable User-Agent so it's not set to default value
			req.Header.Set("User-Agent", "")
		}
	}

	return &httputil.ReverseProxy{Director: director}
}
