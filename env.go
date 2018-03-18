package main

import (
	"fmt"
	"os"
	"runtime"
	"strconv"
)

func WorkerClusterSize() int {
	return fetchCount("RACK_APP_WORKER_CLUSTER_COUNT", runtime.NumCPU())
}

func WorkerThreadCount() int {
	return fetchCount("RACK_APP_WORKER_THREAD_COUNT", 1)
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
