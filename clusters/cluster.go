package clusters

import (
	"net/http"
	"os"

	"github.com/rack-app/server/workers"
)

type Cluster interface {
	Start() []error
	Close() []error
	Signal(os.Signal) []error
	With(func(workers.Worker) error)
	ServeHTTP(http.ResponseWriter, *http.Request)
}

type cluster struct {
	queue   chan workers.Worker
	workers []workers.Worker
	size    int
}

func New(ws []workers.Worker, threadCount int) Cluster {
	workers := make([]workers.Worker, 0, len(ws))
	workers = append(workers, ws...)
	queue := createQueue(workers, threadCount)
	return &cluster{queue: queue, workers: workers}
}

func createQueue(ws []workers.Worker, threadCount int) chan workers.Worker {
	queue := make(chan workers.Worker, len(ws)*threadCount)

	for i := 0; i < threadCount; i++ {
		for _, w := range ws {
			queue <- w
		}
	}

	return queue
}
