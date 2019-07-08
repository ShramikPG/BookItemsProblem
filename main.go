package main

import (
	"os"
	"sync"

	"./models"
)

const (
	// id initialized as constant (seed id)
	id = 699899899
	// WORKERCOUNT initialized as constant
	workercount = 10
)

func generateID(requestID int, errorCount chan int) int {
	if requestID == 0 {
		return id
	}

	incID := requestID + 1
	errCnt := <-errorCount
	if errCnt >= 3 {
		incID = requestID + 50
		errorCount <- 0
		return incID
	}
	errorCount <- errCnt
	return incID
}

func main() {
	var item models.Item
	var reqTrack models.RequestTrack
	failureCount := make(chan int, 1)
	errorCount := make(chan int)
	failureIDs := make(chan []int)
	var gracefulStop = make(chan os.Signal)
	var wg sync.WaitGroup

	models.readJSONFile

}
