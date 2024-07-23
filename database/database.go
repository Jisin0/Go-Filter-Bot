// (c) Jisin0

package database

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/Jisin0/Go-Filter-Bot/utils/config"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var cachedSettings map[int64]*ChatSettings = make(map[int64]*ChatSettings)
var connectionCache map[int64]int64 = make(map[int64]int64) // cache users' connections.
var defaultChatSettings ChatSettings = ChatSettings{
	Stopped: []string{},
}

type ChatSettings struct {
	// Stopped Global Filters
	Stopped []string
}

// A filter object stored in the database
type Filter struct {
	// Unique id of the filter
	ID string `bson:"_id"`
	// Chat where the filter is in effect
	ChatID int64 `bson:"group_id"`
	// The key/text which is filtered
	Text string `bson:"text"`
	// The text content/caption saved
	Content string `bson:"content,omitempty"`
	// The id of a media saved for the filter if any
	FileID string `bson:"file,omitempty"`
	// Buttons/markup saved for a filter if any
	Markup [][]map[string]string `bson:"button,omitempty"`
	// Alerts saved for a filter if any
	Alerts []string `bson:"alert,omitempty"`
	// Length of the text according to which filters are sorted
	Length int `bson:"length"`
	// Type of media saved if any
	MediaType string `bson:"mediaType,omitempty"`
}

// A User Saved in the database.
type User struct {
	// Unique telegram id of the user.
	ID int64 `bson:"_id"`
	// ID of the chat to which the user is connected.
	ConnectedChat int64 `bson:"connected,omitempty"`
}

func NewDatabase() *Database {
	if config.MongodbURI == "" {
		log.Fatal("You must set your 'MONGODB_URI' environmental variable :(")
	}

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(config.MongodbURI))
	if err != nil {
		panic(err)
	}

	db := client.Database("Adv_Auto_Filter")

	return &Database{
		Client: client,
		DB:     db,
		Ucol:   db.Collection("Users"),
		Col:    db.Collection("Main"),
		Mcol:   db.Collection("Manual_Filters"),
	}
}

type Database struct {
	// Mongo Client
	Client *mongo.Client
	// Database
	DB *mongo.Database
	// Users Collection
	Ucol *mongo.Collection
	// Main Collection
	Col *mongo.Collection
	// Manual Filters Collection
	Mcol *mongo.Collection
}

func (db *Database) AddUser(userid int64) error {
	var (
		filter = User{ID: userid}
	)

	res := db.Ucol.FindOne(context.TODO(), filter)
	if res.Err() == mongo.ErrNoDocuments {
		_, err := db.Ucol.InsertOne(context.TODO(), filter)
		if err != nil {
			fmt.Printf("db.adduser: %v", err)
		}
	}

	return nil
}

var statText string = `
╭ ▸ <b>Users</b> : <code>%v</code> 
├ ▸ <b>Filters</b> : <code>%v</code>
╰ ▸ <b>Groups</b> : <code>%v</code>
`

//nolint:errcheck // meh
func (db *Database) Stats() string {
	users, _ := db.Ucol.CountDocuments(context.TODO(), bson.M{})
	chats, _ := db.Col.CountDocuments(context.TODO(), bson.M{})
	mfilters, _ := db.Mcol.CountDocuments(context.TODO(), bson.M{})

	return fmt.Sprintf(statText, users, mfilters, chats)
}

func (db *Database) GetConnection(userID int64) (int64, bool) {
	if c, ok := connectionCache[userID]; ok {
		if c == 0 {
			ok = false
		}

		return c, ok
	}

	res := db.Ucol.FindOne(context.TODO(), bson.D{{Key: "_id", Value: userID}})
	if res.Err() != nil { // this shouldn't happen. the user should be saved in the db.
		connectionCache[userID] = 0
		return 0, false
	}

	var doc User

	res.Decode(&doc)

	connectionCache[userID] = doc.ConnectedChat

	if doc.ConnectedChat != 0 {
		return doc.ConnectedChat, true
	}

	return doc.ConnectedChat, false
}

func (db *Database) ConnectUser(userID, chatID int64) {
	var tf = true

	_, err := db.Ucol.UpdateOne(context.TODO(), bson.D{{Key: "_id", Value: userID}}, bson.D{{Key: "$set", Value: bson.D{{Key: "connected", Value: chatID}}}}, &options.UpdateOptions{Upsert: &tf})
	if err != nil {
		fmt.Printf("db.connectuser: %v\n", err)
	}

	// clear any cache.
	delete(connectionCache, userID)
}

func (db *Database) SaveMfilter(data *Filter) {
	_, err := db.Mcol.InsertOne(context.TODO(), data)
	if err != nil {
		fmt.Printf("db.savemfilter: %v\n", err)
	}
}

func (db *Database) DeleteConnection(userID int64) {
	_, err := db.Ucol.UpdateOne(context.TODO(), bson.D{{Key: "_id", Value: userID}}, bson.D{{Key: "$unset", Value: "connected"}})
	if err != nil {
		fmt.Printf("db.connectuser: %v\n", err)
	}

	// clear any cache.
	delete(connectionCache, userID)
}

// Returns all manual filters for a chat.
func (db *Database) GetMfilters(chatID int64) (*mongo.Cursor, error) {
	return db.Mcol.Aggregate(
		context.TODO(),
		[]bson.D{
			{{Key: "$match", Value: bson.D{{Key: "group_id", Value: chatID}}}},
			{{Key: "$sort", Value: bson.D{{Key: "length", Value: -1}}}},
		},
	)
}

// SearchMfilterClassic uses the traditional way of fetching mfilters by fetching all mfilters of a chat and doing regex queries individually.
// This method is more resource intensive but could be faster for large scale bots with several hundred/thousand groups.
func (db *Database) SearchMfilterClassic(chatID int64, input string) (results []*Filter) {
	res, e := db.GetMfilters(chatID)
	if e != nil {
		fmt.Printf("db.searchmfiltersclassic: %v\n", e)
		return results
	}

	for res.Next(context.TODO()) {
		var f Filter

		err := res.Decode(&f)
		if err != nil {
			fmt.Printf("db.searchmfilterclassic: %v\n", err)
			continue
		}

		text := `(?i)( |^|[^\w])` + f.Text + `( |$|[^\w])`

		pattern := regexp.MustCompile(text)

		m := pattern.FindStringSubmatch(input)
		if len(m) > 0 {
			results = append(results, &f)
		}
	}

	return results
}

// SearchMfilterNew does a regex query on the database shifting some load to mongodb.
func (db *Database) SearchMfilterNew(chatID int64, fields []string, multiFilter bool) (results []*Filter) {
	pattern := "(?i).*\\b(" + strings.Join(fields, "|") + ")\\b.*"
	filter := bson.D{
		{Key: "group_id", Value: chatID},
		{Key: "text", Value: bson.M{"$regex": pattern}},
	}

	if !multiFilter {
		res := db.Mcol.FindOne(context.Background(), filter)
		switch res.Err() {
		case mongo.ErrNoDocuments:
			return results
		case nil:
			var f Filter

			if err := res.Decode(&f); err != nil {
				fmt.Printf("db.searchmfilternew: %v\n", err)
				return results
			}

			return append(results, &f)
		}
	}

	res, err := db.Mcol.Find(context.Background(), filter, options.Find().SetSort(bson.D{{Key: "length", Value: -1}}))
	if err != nil {
		fmt.Printf("db.searchmfilternew: %v\n", err)
		return results
	}

	for res.Next(context.Background()) {
		var r Filter

		err := res.Decode(&r)
		if err != nil {
			fmt.Printf("db.searchmfilternew: %v\n", err)
			continue
		}

		results = append(results, &r)
	}

	return results
}

func (db *Database) GetMfilter(chatID int64, key string) (bson.M, bool) {
	res := db.Mcol.FindOne(context.TODO(), bson.D{{Key: "group_id", Value: chatID}, {Key: "text", Value: key}})

	var b bson.M

	if res.Err() != nil {
		return b, false
	} else {
		err := res.Decode(&b)
		if err != nil {
			fmt.Printf("db.getmfilter: %v\n", err)
			return b, false
		}

		return b, true
	}
}

func (db *Database) StringMfilter(chatID int64) string {
	r, err := db.GetMfilters(chatID)
	if err != nil {
		return fmt.Sprintf("failed to get filters: %v", err)
	}

	var text string

	for r.Next(context.TODO()) {
		var d bson.M

		r.Decode(&d)
		text += fmt.Sprintf("\n• <code>%v</code>", d["text"])
	}

	return text
}

func (db *Database) StopGfilter(chatID int64, key string) {
	var t = true

	_, err := db.Col.UpdateOne(context.TODO(), bson.D{{Key: "_id", Value: chatID}}, bson.D{{Key: "$append", Value: bson.D{{Key: "stopped", Value: key}}}}, &options.UpdateOptions{Upsert: &t})
	if err != nil {
		fmt.Println(err)
	}

	go db.RecacheSettings(chatID)
}

func (db *Database) DeleteMfilter(chatID int64, key string) {
	_, err := db.Mcol.DeleteOne(context.TODO(), bson.D{{Key: "group_id", Value: chatID}, {Key: "text", Value: key}})
	if err != nil {
		fmt.Printf("db.deletemfilter: %v", err)
	}
}

func (db *Database) SetChatSetting(chatID int64, key string, value any) {
	go db.SetDefaultSettings(chatID)

	_, err := db.Col.UpdateOne(context.TODO(), bson.D{{Key: "_id", Value: chatID}}, bson.D{{Key: "$set", Value: bson.D{{Key: key, Value: value}}}})
	if err != nil {
		fmt.Printf("db.connectuser: %v", err)
		return
	}

	go db.RecacheSettings(chatID)
}

func (db *Database) SetDefaultSettings(chatID int64) {
	r := db.Col.FindOne(context.TODO(), bson.D{{Key: "_id", Value: chatID}})
	if r.Err() == mongo.ErrNoDocuments {
		db.Col.InsertOne(context.TODO(), defaultSettings(chatID))
	}
}

func defaultSettings(chatID int64) bson.D {
	return bson.D{
		{Key: "_id", Value: chatID},
	}
}

// A Function To Update Cached Settings With Latest From DB
func (db *Database) RecacheSettings(chatID int64) {
	res := db.Col.FindOne(context.TODO(), bson.D{{Key: "_id", Value: chatID}})
	if res.Err() == mongo.ErrNoDocuments {
		cachedSettings[chatID] = &defaultChatSettings
	} else {
		var r ChatSettings

		err := res.Decode(&r)
		if err != nil {
			fmt.Printf("db.recachesettings: %v", err)
			return
		}

		cachedSettings[chatID] = &r
	}
}

func (db *Database) GetCachedSetting(chatID int64) *ChatSettings {
	s, e := cachedSettings[chatID]
	if !e {
		go db.RecacheSettings(chatID)
		return &defaultChatSettings
	} else {
		return s
	}
}

func (db *Database) StartGfilter(chatID int64, key string) {
	keys := db.GetCachedSetting(chatID).Stopped
	for i, k := range keys {
		if k == key {
			keys[i] = keys[len(keys)-1]
			keys = keys[:len(keys)-1]
		}
	}

	_, err := db.Col.UpdateOne(context.TODO(), bson.D{{Key: "_id", Value: chatID}}, bson.D{{Key: "$set", Value: bson.D{{Key: "stopped", Value: keys}}}})
	if err != nil {
		fmt.Printf("db.startgfilter: %v", err)
	}

	go db.RecacheSettings(chatID)
}

func (db *Database) GetAlert(uniqueID string, index int) string {
	defaultString := "404: Content Not Found !"
	res := db.Mcol.FindOne(context.TODO(), bson.D{{Key: "_id", Value: uniqueID}})

	var f Filter

	res.Decode(&f)

	if len(f.Alerts) < index+1 {
		return defaultString
	} else {
		return f.Alerts[index]
	}
}
