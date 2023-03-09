package dbPackage

import (
	"context"
	"fmt"
	queuepackage "goserver/queuePackage"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var mongoClient *mongo.Client
var mongoContext context.Context
var mongoCancel context.CancelFunc
var mongoError error

type Message struct {
	ID          string `bson:"_id,omitempty" json:"_id,omitempty"`
	UserID      string `bson:"userID,omitempty" json:"userID,omitempty"`
	TextContent string `bson:"textContent,omitempty" json:"textContent,omitempty"`
}

type UserStory struct{
	userStoryID string
	userStoryTitle string
	userStoryDescription string
}

type Meeting struct{
	meetingID string
	meetingTitle string
	meetingUserStories []UserStory
}

// This is a user defined method to close resources.
// This method closes mongoDB connection and cancel context.
func Close(client *mongo.Client, ctx context.Context,
	cancel context.CancelFunc) {

	// CancelFunc to cancel to context
	defer cancel()

	// client provides a method to close
	// a mongoDB connection.
	defer func() {

		// client.Disconnect method also has deadline.
		// returns error if any,
		if err := client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()
}

// This is a user defined method that returns
// a mongo.Client, context.Context,
// context.CancelFunc and error.
// mongo.Client will be used for further database
// operation. context.Context will be used set
// deadlines for process. context.CancelFunc will
// be used to cancel context and resource
// associated with it.
func Connect(uri string) (*mongo.Client, context.Context,
	context.CancelFunc, error) {

	// ctx will be used to set deadline for process, here
	// deadline will of 30 seconds.
	mongoContext, mongoCancel = context.WithTimeout(context.Background(),
		30*time.Second)

	// mongo.Connect return mongo.Client method
	mongoClient, mongoError = mongo.Connect(mongoContext, options.Client().ApplyURI(uri))

	return mongoClient, mongoContext, mongoCancel, mongoError
}

// query is user defined method used to query MongoDB,
// that accepts mongo.client,context, database name,
// collection name, a query and field.

//  database name and collection name is of type
// string. query is of type interface.
// field is of type interface, which limits
// the field being returned.

// query method returns a cursor and error.
func Query(client *mongo.Client, ctx context.Context,
	dataBase, col string, query, field interface{}) (result *mongo.Cursor, err error) {

	// select database and collection.
	collection := client.Database(dataBase).Collection(col)

	// collection has an method Find,
	// that returns a mongo.cursor
	// based on query and field.
	result, err = collection.Find(ctx, query,
		options.Find().SetProjection(field))
	return
}

// Function that takes the text id and returns the Message (text) Struct
func GetMessageFromTextId(textid string) queuepackage.Message {
	// create a filter an option of type interface,
	// that stores bjson objects.
	var filter, option interface{}

	// filter  gets all document,
	// with maths field greater that 70

	objID, objErr := primitive.ObjectIDFromHex(textid)
	if objErr != nil {
		panic(objErr)
	}
	filter = bson.D{{
		Key:   "_id",
		Value: objID,
	}}
	//option remove id field from all documents
	//option = bson.D{{"_id", 0}}
	option = bson.D{}
	// call the query method with client, context,
	// database name, collection  name, filter and option
	// This method returns mongo.cursor and error if any.
	cursor, queryErr := Query(mongoClient, mongoContext, "test",
		"text", filter, option)
	// handle the errors.
	if queryErr != nil {
		panic(queryErr)
	}
	var results []bson.D
	// to get bson object  from cursor,
	// returns error if any.
	if err := cursor.All(mongoContext, &results); err != nil {
		// handle the error
		panic(err)
	}
	var text, err = convertToStruct(results)
	if err != nil {
		panic(err)
	}
	msg := queuepackage.Message{}
	msg.TextID = text[0].ID
	msg.Text = text[0].TextContent
	msg.UserID = text[0].UserID
	return msg

}

// Convert from []primitive.D to array of struct Message
func convertToStruct(docs []primitive.D) ([]Message, error) {
	var myStructs []Message
	for _, doc := range docs {
		id, ok := doc.Map()["_id"].(primitive.ObjectID)
		if !ok {
			return nil, fmt.Errorf("unable to convert _id field to ObjectID")
		}

		myStruct := Message{
			ID:          id.Hex(),
			UserID:      doc.Map()["userID"].(string),
			TextContent: doc.Map()["textContent"].(string),
		}
		myStructs = append(myStructs, myStruct)
	}
	return myStructs, nil
}
