package migrations

import (
	"github.com/ycchuang99/todo-list/models"
)

func Migrate() {
	models.ConnectDatabase()

	models.DB.Query(`
		CREATE TABLE IF NOT EXISTS todo_list (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			title VARCHAR(255) NOT NULL,
			description VARCHAR(255) NOT NULL,
			done_at TEXT
		)
	`)
}
