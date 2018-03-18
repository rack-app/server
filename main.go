package main

import (
	"os"
	"os/signal"
	"runtime/pprof"
	"syscall"
	"time"

	"github.com/rack-app/server/clusters"
	"github.com/rack-app/server/debug"
	"github.com/rack-app/server/workers"
)

func main() {
	sigs := make(chan os.Signal)
	defer func() { close(sigs) }()
	signal.Notify(sigs)

	c := BuildCluster()
	server := NewServer(c)

	OkOrPanic(c.Start)
	go StartServer(server)
	defer PPROF()()

	HandleSignals(sigs, c, server)
}

func BuildCluster() *clusters.Cluster {

	workerClusterSize := WorkerClusterSize()
	workerThreadCount := WorkerThreadCount()
	debug.Printf("W: %v; T: %v\n", workerClusterSize, workerThreadCount)

	ws := make([]*workers.Worker, 0, workerClusterSize)

	for i := 0; i < workerClusterSize; i++ {
		ws = append(ws, workers.New(GetPort(), os.Stdout, os.Stderr))
	}

	return clusters.New(ws, workerThreadCount)
}

func HandleSignals(sigs chan os.Signal, c *clusters.Cluster, server Server) {
receiveSignals:
	for sig := range sigs {
		errs := c.Signal(sig)

		if len(errs) > 0 {
			panic(errs)
		}

		switch sig.String() {
		case syscall.SIGINT.String():
			ShutdownServer(server, 30*time.Second)
			OkOrPanic(c.Close)
			break receiveSignals

		case syscall.SIGTERM.String():
			ShutdownServer(server, 5*time.Second)
			OkOrPanic(c.Close)
			break receiveSignals

		}
	}
}

func PPROF() func() {
	if _, do := os.LookupEnv("PPROF"); !do {
		return func() {}
	}

	f, err := os.Create("./cpu.prof")

	if err != nil {
		panic(err)
	}

	pprof.StartCPUProfile(f)

	return func() {
		pprof.StopCPUProfile()
	}
}
