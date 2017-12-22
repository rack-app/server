package clusters

import (
	"fmt"
	"os"

	"github.com/rack-app/server/workers"
)

func (c *cluster) Start() []error {
	return c.each(func(w workers.Worker) error { return w.Start() })
}

func (c *cluster) Close() []error {
	return c.each(func(w workers.Worker) error { return w.Close() })
}

func (c *cluster) Signal(s os.Signal) []error {
	return c.each(func(w workers.Worker) error { return w.Signal(s) })
}

func (c *cluster) each(fn func(workers.Worker) error) []error {
	errs := []error{}

	for _, w := range c.workers {
		err := fn(w)

		if err != nil {
			errs = append(errs, err)
		}
	}

	return errs
}

func (c *cluster) With(do func(workers.Worker) error) {
	w := <-c.queue
	defer func() { go func() { c.queue <- w }() }()

	if err := do(w); err != nil {
		fmt.Println(err)
		// w.Close()
		// w.Inc()
		// return
	}

}
