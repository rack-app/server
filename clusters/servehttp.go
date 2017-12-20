package clusters

import "net/http"

func (c *cluster) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	w := <-c.queue
	defer func() { go func() { c.queue <- w }() }()
	w.ServeHTTP(rw, req)
}
