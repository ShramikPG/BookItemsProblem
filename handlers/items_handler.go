package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"sync"

	"github.com/ShramikPG/BookItemsProblem/models"
)

const (
	url = "https://api.mercadolibre.com/items/MLA"
)

// GetItem is used to request item with specific from mercadolibre api.
func GetItem(requestID int, wg *sync.WaitGroup, reqTrackChannel chan models.RequestTrackChannel) {
	defer wg.Done()
	var item models.Item
	reqTrChObj := models.RequestTrackChannel{}
	rs, err := http.Get(url + strconv.Itoa(requestID))
	if err != nil {
		fmt.Println("Error: ", err)
		reqTrChObj.RequestID = requestID
		reqTrChObj.ErrorMessage = err.Error()
		reqTrackChannel <- reqTrChObj
		return
	}

	data, err := ioutil.ReadAll(rs.Body)
	if err != nil {
		reqTrChObj.RequestID = requestID
		reqTrChObj.ErrorMessage = err.Error()
		reqTrackChannel <- reqTrChObj
		return
	}

	err = json.Unmarshal(data, &item)
	if err != nil {
		reqTrChObj.RequestID = requestID
		reqTrChObj.ErrorMessage = err.Error()
		reqTrackChannel <- reqTrChObj
		return
	}

	if item.ID == "" {
		reqTrChObj.RequestID = requestID
		reqTrChObj.ErrorMessage = fmt.Sprintf("\nData not Found ReqID: %v, Data: %v\n", requestID, item)
		reqTrackChannel <- reqTrChObj
		return
	}
	reqTrChObj.RequestID = requestID
	reqTrChObj.Data = &item
	reqTrackChannel <- reqTrChObj
}
