package clusters

import "os"

func (wc *cluster) Start() []error {
	return wc.each(func(w Worker) error { return w.Start() })
}

func (wc *cluster) Close() []error {
	return wc.each(func(w Worker) error { return w.Close() })
}

func (wc *cluster) Signal(s os.Signal) []error {
	return wc.each(func(w Worker) error { return w.Signal(s) })
}

func (wc *cluster) each(fn func(Worker) error) []error {
	errs := []error{}

	for _, w := range wc.workers {
		err := fn(w)

		if err != nil {
			errs = append(errs, err)
		}
	}

	return errs
}
