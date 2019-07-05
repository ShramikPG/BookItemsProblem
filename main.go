package main

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"./calls"
)

var requestID int
var gracefulStop = make(chan os.Signal)
var wg sync.WaitGroup

const (
	// ID Initialized as constant
	ID = 699899899
	// WORKERCOUNT initialized as constant
	WORKERCOUNT = 10
)

func generateID() int {
	if requestID == 0 {
		return ID
	}
	incID := requestID + 1
	if calls.CallInf.ErrorCount >= 3 {
		incID = requestID + 50
		calls.CallInf.ErrorCount = 0
	}
	return incID
}

func worker() {
	var callbook calls.CallBook
	wg.Add(WORKERCOUNT)
	for i := 0; i < WORKERCOUNT; i++ {
		requestID = generateID()
		go callbook.Apicall(requestID, &wg)
	}
	wg.Wait()
}

func main() {
	fmt.Printf("\n--- Starting the application ---\n")
	calls.CallInf.ReadJSONFile()
	requestID = calls.CallInf.LastRequestID
	signal.Notify(gracefulStop, syscall.SIGTERM)
	signal.Notify(gracefulStop, syscall.SIGINT)
	go func() {
		sig := <-gracefulStop
		wg.Wait()
		fmt.Printf("Caught signal to terminate: %v", sig)
		fmt.Println("Gracefully Stopping all workers")
		calls.CallInf.WriteJSONFile()
		os.Exit(0)
	}()
	for {
		worker()
	}
}
