package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"sync"
	"time"
)

//Book struct
type Book struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	Error string `json:"error"`
}

/*
1. channels buffered/un-buffered (diff b/w chan 1 vs chan)
*/
func main() {
	errCh := make(chan int, 1)
	var wg sync.WaitGroup
	wg.Add(3)
	url := "https://api.mercadolibre.com/items/MLA"
	id := 699899899
	fmt.Printf("\n--- Starting the application ---\n")
	for i := 1; i < 4; i++ {
		go worker(i, url, id+i, errCh, &wg)
	}
	wg.Wait()
	fmt.Printf("\n--- Stopping the program ---\n")
}

func worker(workerNo int, url string, id int, ch chan int, waitgp *sync.WaitGroup) (response string, err string) {
	errCount := 0
	fmt.Printf("\n In Worker no %v\n", workerNo)
	for {
		select {
		case errCount = <-ch:
		default:
			errCount = 0
		}
		if errCount >= 3 {
			fmt.Printf("\nTo may failed attempts to fetch books, terminating the workers, msg from worker: %v", workerNo)
			waitgp.Done()
			return "Terminating", "Error"
		}
		book := Book{}
		rs, err := http.Get(url + strconv.Itoa(id))
		data, _ := ioutil.ReadAll(rs.Body)
		_ = json.Unmarshal(data, &book)

		if book.Error == "not_found" {
			fmt.Printf("\n Unable to fetch Book, worker no: %v\n", workerNo)
			if errCount++; errCount < 3 {
				fmt.Printf("\n Adding error value to channel from worker %v\n", workerNo)
				ch <- errCount
			} else {
				fmt.Printf("\n Terminating the worker no %v \n", workerNo)
				waitgp.Done()
				return "Terminating", "Error"

			}
			fmt.Printf("\n Error while fetching book with: %v \n error: %v\n Worker No: %v", id, err, workerNo)
			if errCount > 0 && errCount < 3 {
				fmt.Printf("\n From Worker %v \n Sleeping for some time since error has occured for id: %v", workerNo, id)
				time.Sleep(1000 * time.Millisecond)
			}
		} else {
			fmt.Printf("\n Worker No: %v \n Response %v 	\n error %v \n", workerNo, book, err)
		}
		id += 3
	}
}
