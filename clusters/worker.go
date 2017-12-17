package clusters

import (
	"net/http"
	"os"
)

type Worker interface {
	Start() error
	Close() error
	Signal(os.Signal) error
	ServeHTTP(http.ResponseWriter, *http.Request)
}
