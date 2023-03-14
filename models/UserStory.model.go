package models

type UserStory struct{
	UserStoryID string `json:"userStoryID"`
	UserStoryTitle string `json:"userStoryTitle"`
	UserStoryDescription string `json:"userStoryDescription"`
}

type Meeting struct{
	MeetingID string `json:"meetingID"`
	MeetingTitle string `json:"meetingTitle"`
	MeetingUserStories []UserStory `json:"meetingUserStories"`
}