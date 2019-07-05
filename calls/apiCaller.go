package calls

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"sync"
	"time"

	"../structures"
)

//CallBook type is of type Book imported from Structures package
type CallBook structures.Book

const (
	//URL initialization
	URL = "https://api.mercadolibre.com/items/MLA"
)

func checkErrors(requestID int) {
	CallInf.FailureCount++
	CallInf.FailureID = append(CallInf.FailureID, requestID)
	fmt.Println("Error at ID : ", requestID)
	if CallInf.FailureCount > 10 {
		CallInf.ErrorCount++
		CallInf.FailureCount = 0
		fmt.Println("Sleeping for 30 second due to ten consecutive failed requests")
		time.Sleep(30 * time.Second)
	}
}

//Apicall method to request Item form api
func (bk *CallBook) Apicall(requestID int, wg *sync.WaitGroup) {
	if CallInf.LastRequestID < requestID {
		CallInf.LastRequestID = requestID
	}
	var book CallBook
	defer wg.Done()
	rs, err := http.Get(URL + strconv.Itoa(requestID))
	if err != nil {
		fmt.Println("Error: ", err)
		checkErrors(requestID)
		return
	}
	data, _ := ioutil.ReadAll(rs.Body)
	unmErr := json.Unmarshal(data, &book)
	if unmErr != nil {
		fmt.Println("Error: ", unmErr)
		checkErrors(requestID)
		return
	}

	if book.ID == "" {
		fmt.Printf("\nData not Found ReqID: %v, Data: %v\n", requestID, book)
		checkErrors(requestID)
		return
	}
	fmt.Printf("\nBookID: %v, Book-Title: %v\n", book.ID, book.Title)
}
