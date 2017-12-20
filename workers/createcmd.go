package workers

import (
	"io"
	"os/exec"
	"strconv"
)

func createCMD(port int, out, err io.Writer) *exec.Cmd {
	cmd := exec.Command(
		"rackup",
		"--server", "rack-app-receiver",
		"--port", strconv.Itoa(port),
	)

	cmd.Stdout = out
	cmd.Stderr = err
	return cmd
}
