package common

import (
	"app/config"
	db "github.com/dancannon/gorethink"
	"log"
)

func InitDB(cfg *config.Config) *db.Session {

	options := db.ConnectOpts{
		Address:  cfg.DB.Address,
		Database: cfg.DB.Database,
	}

	sess, err := db.Connect(options)

	if err != nil {
		log.Fatalln(err.Error())
	}
	err = db.DbCreate(cfg.DB.Database).Exec(sess)
	if err != nil {
		log.Println(err)
	}

	for _, table := range cfg.DB.Tables {
		_, err = db.Db(cfg.DB.Database).TableCreate(table).RunWrite(sess)
		if err != nil {
			log.Println(err)
		}
	}

	return sess
}