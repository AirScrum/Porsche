package models
import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)
/*
Represents the UserStory type coming from the NLP Model
*/
type UserStory struct{
	UserStoryTitle string `json:"userStoryTitle" bson:"userStoryTitle"`
	UserStoryDescription string `json:"userStoryDescription" bson:"userStoryDescription"`
}

/*
Represents the UserStory schema in the database.
*/
type UserStoryModel struct{
	UserStoryTitle string `json:"userStoryTitle" bson:"userStoryTitle"`
	UserStoryDescription string `json:"userStoryDescription" bson:"userStoryDescription"`
	UserID primitive.ObjectID `json:"userID" bson:"userID"`
	TextID primitive.ObjectID `json:"textID" bson:"textID"`
}