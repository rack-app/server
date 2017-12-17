package clusters

import "net/http"
import "fmt"

func (c *cluster) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	w := <-c.queue
	defer func() { c.queue <- w }()
	fmt.Println(w)
	w.ServeHTTP(rw, req)
}
