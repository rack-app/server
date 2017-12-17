package main

import (
	"io"
	"net"
	"os"
	"path/filepath"
)

func main() {

	socktest.Sockets

	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	tmpDir := filepath.Join(wd, "tmp")

	if err := os.MkdirAll(tmpDir, os.ModePerm); err != nil {
		panic(err)
	}

	socketPath := filepath.Join(tmpDir, "test.socket")
	defer os.Remove(socketPath)

	go Client(socketPath)
	Server(socketPath)
}

func Server(socketPath string) {

}

func Client(socketPath string) {
	for {
		if _, err := os.Stat(socketPath); err == nil {
			break
		}
	}

	l, err := net.ListenUnix("unix", socketPath)
	if err != nil {
		panic(err)
	}
	defer l.Close()

	for {
		conn, err := l.Accept()

		if err != nil {
			panic(err)
		}

		if _, err := io.Copy(os.Stdout, conn); err != nil {
			panic(err)
		}

		conn.Close()
	}

}
