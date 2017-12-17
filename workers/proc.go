package workers

import "os"

func (w *worker) Start() error {
	return w.cmd.Start()
}

func (w *worker) Close() error {
	if w.Exited() {
		return nil
	}

	if err := w.Signal(os.Interrupt); err != nil {
		return err
	}

	return w.cmd.Wait()
}

func (w *worker) Exited() bool {
	return w.cmd.ProcessState != nil && w.cmd.ProcessState.Exited()
}

func (w *worker) Signal(s os.Signal) error {
	return w.cmd.Process.Signal(s)
}
