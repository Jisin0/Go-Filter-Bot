// (c) Jisin0

package autodelete

import (
	"fmt"
	"time"

	"github.com/Jisin0/Go-Filter-Bot/utils/config"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

var db *sqlx.DB

// Data to be entered into the autodelete db.
type AutodelData struct {
	ChatID    int64     `db:"chat_id"`
	MessageID int64     `db:"message_id"`
	ExpTime   time.Time `db:"exp_time"`
}

func init() {
	if config.AutoDelete == 0 {
		return
	}

	var err error

	db, err = sqlx.Open("sqlite3", "./cache.sqlite")
	if err != nil {
		fmt.Printf("failed to open autodelete db: %v\n", err)
		return
	}

	// Create a table to store cache data (if it doesn't exist)
	createTableSQL := `CREATE TABLE IF NOT EXISTS cache (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        chat_id INTEGER,
        message_id INTEGER,
        exp_time DATETIME,
        UNIQUE(chat_id, message_id)
    );`
	_, err = db.Exec(createTableSQL)
	if err != nil {
		fmt.Printf("failed to create autodelete db: %v\n", err)
	}
}

// Insert query.
var insertSQL = `INSERT INTO cache (chat_id, message_id, exp_time) VALUES (:chat_id, :message_id, :exp_time)
	ON CONFLICT(chat_id, message_id) DO UPDATE SET exp_time=excluded.exp_time;`

// Insert a message to be deleted into the database:
//
// - data : Data to be entered.
// - seconds : Expiration time in seconds.
func InsertAutodel(data AutodelData, seconds int64) error {
	data.ExpTime = time.Now().Add(time.Duration(seconds) * time.Second)

	_, err := db.NamedExec(insertSQL, data)
	if err != nil {
		return err
	}

	return nil
}

// Select query.
var selectSQL = `SELECT chat_id, message_id, exp_time FROM cache WHERE exp_time <= ?`

// Delete query.
var deleteSQL = `DELETE FROM cache WHERE chat_id = ? AND message_id = ?`

// TODO:
// Delete all entries by bot_token when fetching bot fails. This may cause false positives though.

// Runs an autodelete ticker job that runs every minute and queries the db for expired messages and deletes them.
func RunAutodel(bot *gotgbot.Bot) {
	ticker := time.NewTicker(time.Minute)
	quit := make(chan struct{})

	for {
		select {
		case <-ticker.C:
			var result []AutodelData

			err := db.Select(&result, selectSQL, time.Now())
			if err != nil {
				fmt.Printf("autodelete db query failed: %v", err)
				break
			}

			for _, r := range result {
				bot.DeleteMessage(r.ChatID, r.MessageID, &gotgbot.DeleteMessageOpts{}) //nolint:errcheck // I don't care for an error here.

				_, err := db.Exec(deleteSQL, r.ChatID, r.MessageID)
				if err != nil {
					fmt.Printf("failed to delete autodel db entry: %v", err)
					continue
				}
			}

		case <-quit:
			ticker.Stop()
			return
		}

	}
}
