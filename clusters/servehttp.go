package clusters

import (
	"net/http"

	"github.com/rack-app/server/workers"
)

func (c *Cluster) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	c.With(func(w *workers.Worker) error {

		w.ServeHTTP(rw, req)

		return nil
	})
}
