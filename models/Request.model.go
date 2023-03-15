
package models
/*
This struct is used, when we have a request from the gateway with the text id, to get the corresponding text from database, then send it to text queue
*/
type Request struct {
	TextID string `json:"textID"`
}
