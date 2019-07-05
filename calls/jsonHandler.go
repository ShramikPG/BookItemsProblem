package calls

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"../structures"
)

//CallInfo is of type callerinfo imported from structures package
type CallInfo structures.CallerInfo

// CallInf to store last requestId and Failure IDsetc
var CallInf CallInfo

//WriteJSONFile method of callInfo type to write last requestID and FalureIDs in a JSON file
func (callInfo *CallInfo) WriteJSONFile() {
	file, err := json.MarshalIndent(callInfo, "", " ")
	if err != nil {
		fmt.Printf("Error: %v", err)
		return
	}
	ioerr := ioutil.WriteFile("JSONtrackfile.json", file, 0644)
	if ioerr != nil {
		fmt.Printf("Error From WriteFile: %v", ioerr)
		return
	}
}

//ReadJSONFile method of callInfo type to read last requestID from JSON File
func (callInfo *CallInfo) ReadJSONFile() {
	data, ioerr := ioutil.ReadFile("JSONtrackfile.json")
	if ioerr != nil {
		fmt.Printf("Error: %v", ioerr)
		return
	}
	err := json.Unmarshal(data, &callInfo)
	if err != nil {
		fmt.Printf("Error from ReadFile: %v", err)
	}
}
