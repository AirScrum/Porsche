package models

type UserStory struct{
	UserStoryID string `json:"userStoryID" bson:"userStoryID"`
	UserStoryTitle string `json:"userStoryTitle" bson:"userStoryTitle"`
	UserStoryDescription string `json:"userStoryDescription" bson:"userStoryDescription"`
}

type Meeting struct{
	MeetingID string `json:"meetingID" bson:"meetingID"`
	MeetingTitle string `json:"meetingTitle" bson:"meetingTitle"`
	MeetingUserStories []UserStory `json:"meetingUserStories" bson:"meetingUserStories"`
}