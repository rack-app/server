package workers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"sync"
)

var OUT = os.Stdout
var ERR = os.Stderr

type Worker interface {
	Addr() string
	Start() error
	Close() error
	Signal(os.Signal) error
	ServeHTTP(http.ResponseWriter, *http.Request)
}

type worker struct {
	addr string
	cmd  *exec.Cmd
	wg   *sync.WaitGroup
}

func New(port int, out, err io.Writer) Worker {
	return &worker{
		wg:   &sync.WaitGroup{},
		addr: fmt.Sprintf(":%v", port),
		cmd:  createCMD(port, out, err),
	}
}
