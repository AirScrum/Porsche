package models
/*
This struct is for sending to the text message queue, for the NLP model to take this text and convert it to user stories
The TextID and UserID here are string format instead of ObjectID to decrease the overhead on the NLP Model in terms of
Marshalling.
*/
type Message struct {
	TextID string `json:"textID"`
	Text   string `json:"text"`
	UserID string `json:"userID"`
}
