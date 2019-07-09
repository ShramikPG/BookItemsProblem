package main

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/ShramikPG/BookItemsProblem/handlers"
	"github.com/ShramikPG/BookItemsProblem/models"
)

const (
	id          = 699899899
	workercount = 10
)

func generateID(requestID *int, reqTrack *models.RequestTrack) int {
	if *requestID == 0 {
		return id
	}
	errCnt := reqTrack.ErrorCount
	if errCnt >= 3 {
		*requestID = *requestID + 50
		reqTrack.ErrorCount = 0
		return *requestID
	}

	*requestID++
	reqTrack.ErrorCount = errCnt
	return *requestID
}

func itemFetcher(wg *sync.WaitGroup, requestID *int, reqTrack *models.RequestTrack) {
	wg.Add(workercount)
	respChannel := make(chan models.RequestTrackChannel)
	for i := 0; i < workercount; i++ {
		*requestID = generateID(requestID, reqTrack)
		go handlers.GetItem(*requestID, wg, respChannel)
	}
	for j := 0; j < workercount; j++ {
		reqTrChObj := <-respChannel
		if reqTrChObj.ErrorMessage != "" {
			checkErrors(reqTrChObj, reqTrack)
			continue
		}
		data := reqTrChObj.Data
		fmt.Println("Item Fetched From Server ", *data)
		if reqTrack.LastSuccessfulID < reqTrChObj.RequestID {
			reqTrack.LastSuccessfulID = reqTrChObj.RequestID
		}
		if reqTrack.LastRequestID < reqTrChObj.RequestID {
			reqTrack.LastRequestID = reqTrChObj.RequestID
		}
	}
	wg.Wait()
}

func checkErrors(reqTrChObj models.RequestTrackChannel, reqTrack *models.RequestTrack) {
	if reqTrack.LastRequestID < reqTrChObj.RequestID {
		reqTrack.LastRequestID = reqTrChObj.RequestID
	}
	reqTrack.FailureIDs = append(reqTrack.FailureIDs, reqTrChObj.RequestID)
	fmt.Println(reqTrChObj.ErrorMessage)
	reqTrack.FailureCount++
	if reqTrack.FailureCount > 10 {
		reqTrack.ErrorCount++
		fmt.Println("Sleeping for 3 seconds")
		time.Sleep(3 * time.Second)
		reqTrack.FailureCount = 0
	}
}

func main() {
	//var reqTrackFile models.RequestTrackFile
	var gracefulStop = make(chan os.Signal)
	var reqTrack models.RequestTrack
	var wg sync.WaitGroup
	fmt.Printf("\n--- Starting the application ---\n")
	models.ReadJSONFile(&reqTrack)
	requestID := reqTrack.LastRequestID
	signal.Notify(gracefulStop, syscall.SIGTERM)
	signal.Notify(gracefulStop, syscall.SIGINT)
	go func() {
		sig := <-gracefulStop
		fmt.Printf("Caught signal to terminate: %v", sig)
		fmt.Println("Gracefully Stopping all workers")
		models.WriteJSONFile(&reqTrack)
		wg.Wait()
		os.Exit(0)
	}()
	for {
		itemFetcher(&wg, &requestID, &reqTrack)
	}

}
