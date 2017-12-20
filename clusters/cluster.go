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

func New(ws []Worker, threadCount int) Cluster {
	workers := make([]Worker, 0, len(ws))
	workers = append(workers, ws...)
	queue := createQueue(workers, threadCount)
	return &cluster{queue: queue, workers: workers}
}

func createQueue(workers []Worker, threadCount int) chan Worker {
	queue := make(chan Worker, len(workers)*threadCount)

	for i := 0; i < threadCount; i++ {
		for _, w := range workers {
			queue <- w
		}
	}

	return queue
}
