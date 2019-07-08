package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"sync"
	"time"

	"../models"
)

const (
	url = "https://api.mercadolibre.com/items/MLA"
)

// GetItem is used to request item with specific from mercadolibre api.
func GetItem(requestID int, wg *sync.WaitGroup, reqTrack *models.RequestTrack,
	failureCount chan int, errorCount chan int, lastRequestID chan int, failureIDs chan []int) {
	defer wg.Done()
	lastReqID := <-lastRequestID
	if lastReqID < requestID {
		lastReqID = requestID
	}
	lastRequestID <- requestID
	var item models.Item
	rs, err := http.Get(url + strconv.Itoa(requestID))
	if err != nil {
		fmt.Println("Error: ", err)
		checkErrors(requestID, failureCount, errorCount, failureIDs)
		return
	}

	data, err := ioutil.ReadAll(rs.Body)
	if err != nil {
		fmt.Println("Errors: ", err)
		return
	}

	err = json.Unmarshal(data, &item)
	if err != nil {
		fmt.Println("Error: ", err)
		checkErrors(requestID, failureCount, errorCount, failureIDs)
		return
	}

	if item.ID == "" {
		fmt.Printf("\nData not Found ReqID: %v, Data: %v\n", requestID, item)
		checkErrors(requestID, failureCount, errorCount, failureIDs)
		return
	}

	fmt.Printf("\nBookID: %v, Book-Title: %v\n", item.ID, item.Title)
}

func checkErrors(requestID int, failureCount chan int, errorCount chan int, failureIDs chan []int) {
	failedIDs := <-failureIDs
	failedIDs = append(failedIDs, requestID)
	failCount, ok := <-failureCount
	if !ok {
		fmt.Println("FailureCount channel is closed")
		return
	}

	errCount, ok := <-errorCount
	if !ok {
		fmt.Println("ErrorCount channel is closed")
		return
	}

	failCount++
	if failCount > 10 {
		errCount++
		failCount = 0
		fmt.Println("Sleeping for 3 second due to ten failed requests")
		time.Sleep(3 * time.Second)
	}

	failureCount <- failCount
	errorCount <- errCount
}
