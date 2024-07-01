// (c) Jisin0

package database

import (
	"context"
	"fmt"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var cachedSettings map[int64]*ChatSettings = make(map[int64]*ChatSettings)
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
	Id string `bson:"_id"`
	// Chat where the filter is in effect
	ChatId int64 `bson:"group_id"`
	// The key/text which is filtered
	Text string `bson:"text"`
	// The text content/caption saved
	Content string `bson:"content"`
	// The id of a media saved for the filter if any
	FileID string `bson:"file"`
	// Buttons/markup saved for a filter if any
	Markup [][]map[string]string `bson:"button"`
	// Alerts saved for a filter if any
	Alerts []string `bson:"alert"`
	// Length of the text according to which filters are sorted
	Length int `bson:"length"`
	// Type of media saved if any
	MediaType string `bson:"mediaType"`
}

func NewDatabase() *Database {
	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		log.Fatal("You must set your 'MONGODB_URI' environmental variable :(")
	}

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}

	db := client.Database("Adv_Auto_Filter")

	return &Database{
		Client: client,
		Db:     db,
		Ucol:   db.Collection("Users"),
		Col:    db.Collection("Main"),
		Mcol:   db.Collection("Manual_Filters"),
	}
}

type Database struct {
	// Mongo Client
	Client *mongo.Client
	// Database
	Db *mongo.Database
	// Users Collection
	Ucol *mongo.Collection
	// Main Collection
	Col *mongo.Collection
	// Manual Filters Collection
	Mcol *mongo.Collection
}

func (db *Database) AddUser(userid int64) error {
	var (
		result bson.M
		filter = bson.D{{Key: "_id", Value: userid}}
	)

	err := db.Ucol.FindOne(context.TODO(), filter).Decode(&result)
	if err == mongo.ErrNoDocuments {
		db.Ucol.InsertOne(context.TODO(), filter)
	}

	return nil
}

//nolint:errcheck // meh
func (db *Database) Stats() string {
	users, _ := db.Ucol.CountDocuments(context.TODO(), bson.M{})
	chats, _ := db.Col.CountDocuments(context.TODO(), bson.M{})
	mfilters, _ := db.Mcol.CountDocuments(context.TODO(), bson.M{})

	return fmt.Sprintf("<u>Cᥙɾɾҽɳ𝜏 Ɗα𝜏αßαടҽ S𝜏α𝜏ട </u>:\n\nUsᴇʀs: %v\nMᴀɴᴜᴀʟ Fɪʟᴛᴇʀs: %v\nCᴜsᴛᴏᴍɪᴢᴇᴅ Cʜᴀᴛs: %v", users, mfilters, chats)
}

func (db *Database) GetConnection(userID int64) (int64, bool) {
	res := db.Ucol.FindOne(context.TODO(), bson.D{{Key: "_id", Value: userID}})

	var doc bson.M
	res.Decode(&doc)

	val, ok := doc["connected"]
	if val == nil {
		return 0, false
	}

	return val.(int64), ok
}

func (db *Database) ConnectUser(userID int64, chatID int64) {
	var tf bool = true
	_, err := db.Ucol.UpdateOne(context.TODO(), bson.D{{Key: "_id", Value: userID}}, bson.D{{Key: "$set", Value: bson.D{{Key: "connected", Value: chatID}}}}, &options.UpdateOptions{Upsert: &tf})
	if err != nil {
		fmt.Printf("db.connectuser: %v\n", err)
	}
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
}

func (db *Database) GetMfilters(chatID int64) (*mongo.Cursor, error) {
	return db.Mcol.Aggregate(
		context.TODO(),
		[]bson.D{
			{{Key: "$match", Value: bson.D{{Key: "group_id", Value: chatID}}}},
			{{Key: "$sort", Value: bson.D{{Key: "length", Value: -1}}}},
		},
	)
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
	var t bool = true

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
	// go db.RecacheSettings(chatID)
}

func (db *Database) GetAlert(uniqueID string, index int) string {
	defaultString := "Button Does Not Exist :("
	res := db.Mcol.FindOne(context.TODO(), bson.D{{Key: "unique_id", Value: uniqueID}})

	var f Filter

	res.Decode(&f)

	if len(f.Alerts) < index+1 {
		return defaultString
	} else {
		return f.Alerts[index]
	}
}
