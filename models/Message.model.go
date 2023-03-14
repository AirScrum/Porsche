package models
/*
This struct is for sending to the text message queue, for the NLP model to take this text and convert it to user stories
*/
type Message struct {
	TextID string `json:"textID"`
	Text   string `json:"text"`
	UserID string `json:"userID"`
}
