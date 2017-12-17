package clusters

import (
	"net/http"
	"os"
)

type Cluster interface {
	Start() []error
	Close() []error
	Signal(os.Signal) []error
	ServeHTTP(http.ResponseWriter, *http.Request)
}

type cluster struct {
	queue   chan Worker
	workers []Worker
	size    int
}

func New(ws []Worker) Cluster {
	workers := make([]Worker, 0, len(ws))
	workers = append(workers, ws...)
	queue := make(chan Worker, len(ws))
	for _, w := range workers {
		queue <- w
	}
	return &cluster{queue: queue, workers: workers}
}
