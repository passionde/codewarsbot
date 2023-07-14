package initApp

import (
	"database/sql"
	"github.com/Yarik-xxx/CodeWarsRestApi/cmd/tools/helper"
	"log"
)

func CreateTable() {
	config := helper.InitConfig()

	queryChallenges := `CREATE TABLE challenge (
    id    text PRIMARY KEY NOT NULL,
    info json NOT NULL,
    lastUpdate timestamp with time zone DEFAULT now()
);
`
	db, err := sql.Open("postgres", config.Store.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.Query(queryChallenges)
	if err != nil {
		log.Fatal(err)
	}
}
