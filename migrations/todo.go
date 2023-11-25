package migrations

import (
	// "log"

	"github.com/ycchuang99/todo-list/models"
)

func Migrate() {
	models.ConnectDatabase()

	_, err := models.DB.Exec(`
		CREATE TABLE todo_list (
			id int NOT NULL,
			title varchar(128) NOT NULL,
			description varchar(1024) DEFAULT NULL,
			done_at timestamp NULL DEFAULT NULL
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
	`)
	if err != nil {
		panic(err.Error())
	}

	_, err = models.DB.Exec(`ALTER TABLE todo_list ADD PRIMARY KEY (id);`)
	if err != nil {
		panic(err.Error())
	}

	_, err = models.DB.Exec(`ALTER TABLE todo_list MODIFY id int NOT NULL AUTO_INCREMENT;`)
	if err != nil {
		panic(err.Error())
	}
}
