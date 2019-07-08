package models

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

//RequestTrack struct is used to store last request ID , error Counts,
type RequestTrack struct {
	LastRequestID int   `json:"last_request_id"`
	FailureCount  int   `json:"failure_count"`
	ErrorCount    int   `json:"err_count"`
	FailureIDs    []int `json:"failure_id"`
}

// WriteJSONFile function is used to write data from reqTrack object
func WriteJSONFile(reqTrack *RequestTrack) {
	file, err := json.MarshalIndent(reqTrack, "", " ")
	if err != nil {
		fmt.Printf("Error: %v", err)
		return
	}
	err = ioutil.WriteFile("JSONtrackfile.json", file, 0644)
	if err != nil {
		fmt.Printf("Error From WriteFile: %v", err)
		return
	}
}

//ReadJSONFile method is used to read last requestID from JSON File
func ReadJSONFile(reqTrack *RequestTrack) {
	data, err := ioutil.ReadFile("JSONtrackfile.json")
	if err != nil {
		fmt.Printf("Error: %v", err)
		return
	}
	err = json.Unmarshal(data, &reqTrack)
	if err != nil {
		fmt.Printf("Error from ReadFile: %v", err)
	}
}
