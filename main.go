package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"runtime/pprof"
	"strconv"
	"syscall"
	"time"

	"github.com/rack-app/server/clusters"
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

func BuildCluster() clusters.Cluster {
	ws := make([]clusters.Worker, 0, runtime.NumCPU())

	for i := 0; i < WorkerClusterSize(); i++ {
		ws = append(ws, workers.New(GetPort(), os.Stdout, os.Stderr))
	}

	return clusters.New(ws, WorkerThreadCount())
}

func WorkerClusterSize() int {
	return fetchCount("RACK_APP_WORKER_CLUSTER_COUNT", runtime.NumCPU())
}

func WorkerThreadCount() int {
	return fetchCount("RACK_APP_WORKER_THREAD_COUNT", 16)
}

func fetchCount(envKey string, defaultValue int) int {
	rawCount, isSet := os.LookupEnv(envKey)

	if !isSet {
		return defaultValue
	}

	count, err := strconv.Atoi(rawCount)

	if err != nil {
		fmt.Println(fmt.Sprintf("%s must be a valid number"))
		os.Exit(1)
	}

	return count
}

func HandleSignals(sigs chan os.Signal, c clusters.Cluster, server *http.Server) {
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

// GetFreePort asks the kernel for a free open port that is ready to use.
func GetFreePort() (int, error) {
	addr, err := net.ResolveTCPAddr("tcp", ":0")
	if err != nil {
		return 0, err
	}

	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return 0, err
	}
	defer l.Close()
	return l.Addr().(*net.TCPAddr).Port, nil
}

// GetPort is deprecated, use GetFreePort instead
// Ask the kernel for a free open port that is ready to use
func GetPort() int {
	port, err := GetFreePort()
	if err != nil {
		panic(err)
	}
	return port
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
