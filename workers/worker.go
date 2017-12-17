package workers

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"os/exec"
	"sync"
)

var OUT = os.Stdout
var ERR = os.Stderr

type Worker interface {
	Start() error
	Close() error
	Signal(os.Signal) error
	ServeHTTP(http.ResponseWriter, *http.Request)
}

type worker struct {
	cmd   *exec.Cmd
	wg    *sync.WaitGroup
	proxy *httputil.ReverseProxy
}

func New(port int, out, err io.Writer) Worker {
	c := createCMD(port, out, err)
	w := &sync.WaitGroup{}
	u, _ := url.Parse(fmt.Sprintf("http://localhost:%v", port))
	rp := httputil.NewSingleHostReverseProxy(u)
	return &worker{wg: w, cmd: c, proxy: rp}
}
