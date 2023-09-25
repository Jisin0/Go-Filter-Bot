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
	//Stopped Global Filters
	Stopped []string
}

type Filter struct {
	//A filter object stored in the database

	//Unique id of the filter
	Id string `bson:"_id"`

	//Chat where the filter is in effect
	ChatId int64 `bson:"group_id"`

	//The key/text which is filtered
	Text string `bson:"text"`

	//The text content/caption saved
	Content string `bson:"content"`

	//The id of a media saved for the filter if any
	FileID string `bson:"file"`

	//Buttons/markup saved for a filter if any
	Markup [][]map[string]string `bson:"button"`

	//Alerts saved for a filter if any
	Alerts []string `bson:"alert"`

	//Length of the text according to which filters are sorted
	Length int `bson:"length"`

	//Type of media saved if any
	MediaType string `bson:"mediaType"`
}

func NewDatabase() Database {

	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		log.Fatal("You must set your 'MONGODB_URI' environmental variable :(")
	}

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}

	db := client.Database("Adv_Auto_Filter")

	return Database{
		Client: client,
		Db:     db,
		Ucol:   db.Collection("Users"),
		Col:    db.Collection("Main"),
		Mcol:   db.Collection("Manual_Filters"),
	}
}

type Database struct {

	//Mongo Client
	Client *mongo.Client

	//Database
	Db *mongo.Database

	//Users Collection
	Ucol *mongo.Collection

	//Main Collection
	Col *mongo.Collection

	//Manual Filters Collection
	Mcol *mongo.Collection
}

func (db Database) AddUser(userid int64) error {
	var result bson.M
	filter := bson.D{{Key: "_id", Value: userid}}
	err := db.Ucol.FindOne(context.TODO(), filter).Decode(&result)

	if err == mongo.ErrNoDocuments {
		db.Ucol.InsertOne(context.TODO(), filter)
	}

	return nil
}

func (db Database) Stats() string {
	users, _ := db.Ucol.CountDocuments(context.TODO(), bson.M{})
	chats, _ := db.Col.CountDocuments(context.TODO(), bson.M{})
	mfilters, _ := db.Mcol.CountDocuments(context.TODO(), bson.M{})

	return fmt.Sprintf("<u>Cᥙɾɾҽɳ𝜏 Ɗα𝜏αßαടҽ S𝜏α𝜏ട </u>:\n\nUsᴇʀs: %v\nMᴀɴᴜᴀʟ Fɪʟᴛᴇʀs: %v\nCᴜsᴛᴏᴍɪᴢᴇᴅ Cʜᴀᴛs: %v", users, mfilters, chats)
}

func (db Database) GetConnection(user_id int64) (int64, bool) {
	res := db.Ucol.FindOne(context.TODO(), bson.D{{Key: "_id", Value: user_id}})
	var doc bson.M
	res.Decode(&doc)
	val, ok := doc["connected"]
	if val == nil {
		return 0, false
	}
	return val.(int64), ok
}

func (db Database) ConnectUser(user_id int64, chat_id int64) {
	var tf bool = true
	db.Ucol.UpdateOne(context.TODO(), bson.D{{Key: "_id", Value: user_id}}, bson.D{{Key: "$set", Value: bson.D{{Key: "connected", Value: chat_id}}}}, &options.UpdateOptions{Upsert: &tf})
}

func (db Database) SaveMfilter(data Filter) {
	_, err := db.Mcol.InsertOne(context.TODO(), data)
	if err != nil {
		fmt.Println(err)
	}
}

func (db Database) DeleteConnection(user_id int64) {
	db.Ucol.UpdateOne(context.TODO(), bson.D{{Key: "_id", Value: user_id}}, bson.D{{Key: "$unset", Value: "connected"}})
}

func (db Database) GetMfilters(chat_id int64) (*mongo.Cursor, error) {
	return db.Mcol.Aggregate(
		context.TODO(),
		[]bson.D{
			{{Key: "$match", Value: bson.D{{Key: "group_id", Value: chat_id}}}},
			{{Key: "$sort", Value: bson.D{{Key: "length", Value: -1}}}},
		},
	)
}

func (db Database) GetMfilter(chat_id int64, key string) (bson.M, bool) {
	res := db.Mcol.FindOne(context.TODO(), bson.D{{Key: "group_id", Value: chat_id}, {Key: "text", Value: key}})
	var b bson.M

	if res.Err() != nil {
		return b, false
	} else {
		res.Decode(&b)
		return b, true
	}
}

func (db Database) StringMfilter(chat_id int64) string {
	r, _ := db.GetMfilters(chat_id)
	var text string
	for r.Next(context.TODO()) {
		var d bson.M
		r.Decode(&d)
		text += fmt.Sprintf("\n• <code>%v</code>", d["text"])
	}

	return text
}

func (db Database) StopGfilter(chat_id int64, key string) {
	var t bool = true
	_, err := db.Col.UpdateOne(context.TODO(), bson.D{{Key: "_id", Value: chat_id}}, bson.D{{Key: "$append", Value: bson.D{{Key: "stopped", Value: key}}}}, &options.UpdateOptions{Upsert: &t})
	if err != nil {
		fmt.Println(err)
	}
	go db.RecacheSettings(chat_id)
}

func (db Database) DeleteMfilter(chat_id int64, key string) {
	db.Mcol.DeleteOne(context.TODO(), bson.D{{Key: "group_id", Value: chat_id}, {Key: "text", Value: key}})
}

func (db Database) SetChatSetting(chat_id int64, key string, value any) {
	go db.SetDefaultSettings(chat_id)
	db.Col.UpdateOne(context.TODO(), bson.D{{Key: "_id", Value: chat_id}}, bson.D{{Key: "$set", Value: bson.D{{Key: key, Value: value}}}})
	go db.RecacheSettings(chat_id)
}

func (db Database) SetDefaultSettings(chat_id int64) {
	r := db.Col.FindOne(context.TODO(), bson.D{{Key: "_id", Value: chat_id}})
	if r.Err() == mongo.ErrNoDocuments {
		db.Col.InsertOne(context.TODO(), defaultSettings(chat_id))
	}
}

func defaultSettings(chat_id int64) bson.D {
	return bson.D{
		{Key: "_id", Value: chat_id},
	}
}

func (db Database) RecacheSettings(chat_id int64) {
	//A Function To Update Cached Settings With Latest From DB

	res := db.Col.FindOne(context.TODO(), bson.D{{Key: "_id", Value: chat_id}})
	if res.Err() == mongo.ErrNoDocuments {
		cachedSettings[chat_id] = &defaultChatSettings
	} else {
		var dict bson.M
		res.Decode(&dict)

		var stopped []string

		if a, b := dict["stopped"]; b {
			stopped = a.([]string)
		}

		settings := ChatSettings{
			Stopped: stopped,
		}

		cachedSettings[chat_id] = &settings
	}

}

func (db Database) GetCachedSetting(chat_id int64) *ChatSettings {
	s, e := cachedSettings[chat_id]
	if !e {
		go db.RecacheSettings(chat_id)
		return &defaultChatSettings
	} else {
		return s
	}

}

func (db Database) StartGfilter(chat_id int64, key string) {
	keys := db.GetCachedSetting(chat_id).Stopped
	for i, k := range keys {
		if k == key {
			keys[i] = keys[len(keys)-1]
			keys = keys[:len(keys)-1]
		}
	}

	db.Col.UpdateOne(context.TODO(), bson.D{{Key: "_id", Value: chat_id}}, bson.D{{Key: "$set", Value: bson.D{{Key: "stopped", Value: keys}}}})
	//go db.RecacheSettings(chat_id)
}

func (db Database) GetAlert(uniqueID string, index int) string {
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
