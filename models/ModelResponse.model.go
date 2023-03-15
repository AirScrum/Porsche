package models
/*
Represents the response from the machine learning model.
*/
type ModelResponse struct{
	UserID string `json:"userID" bson:"userID"`
	TextID string `json:"textID" bson:"textID"`
	UserStories []UserStory `json:"userStories" bson:"userStories"`
}