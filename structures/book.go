package structures

// Book Structure
type Book struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}

//CallerInfo struct
type CallerInfo struct {
	LastRequestID int   `json:"lastrequestid"`
	FailureCount  int   `json:"failurecount"`
	ErrorCount    int   `json:"errcount"`
	FailureID     []int `json:"failureid"`
}
