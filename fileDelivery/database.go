package fileDelivery

import (
	"github.com/go-martini/martini"
	"gopkg.in/mgo.v2"
	"os"
	"log"
)

/*
I want to use a different database for my tests,
so I'll embed *mgo.Session and store the database name.
*/
type DatabaseSession struct {
	*mgo.Session
	databaseName string
	logger   *log.Logger
}

/*
Connect to the local MongoDB and set up the database.
*/
func NewSession(name string) *DatabaseSession {
	logger := log.New(os.Stdout, "[restfiledelivery] ", 0)
	db_server_seed := os.Getenv("DATABASE_PORT_27017_TCP_ADDR")

	if db_server_seed == "" {
		db_server_seed = "mongodb://127.0.0.1:27017"
	}

	logger.Printf("Dialing Mongo session at %s \n",db_server_seed)

	// session, err := mgo.Dial("mongodb://root:root@ds057224.mongolab.com:57224/requestfiles")
	session, err := mgo.Dial(db_server_seed)
	if err != nil {
		panic(err)
	} else {
		logger.Printf("Mongo session established with %s %s \n",db_server_seed,name)
	}
	addIndexToRequestedFiles(session.DB(name))
	return &DatabaseSession{session, name, logger}
}

func addIndexToRequestedFiles(db *mgo.Database) {
	index := mgo.Index{
		Key:      []string{"requestid"},
		Unique:   true,
		DropDups: true,
	}
	indexErr := db.C("requestedFiles").EnsureIndex(index)
	if indexErr != nil {
		panic(indexErr)
	}
}

func (session *DatabaseSession) Database() martini.Handler {
	return func(context martini.Context) {
		s := session.Clone()
		context.Map(s.DB(session.databaseName))
		defer s.Close()
		context.Next()
	}
}
