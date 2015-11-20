package fileDelivery

import (
	"github.com/go-martini/martini"
	"gopkg.in/mgo.v2"
)

/*
I want to use a different database for my tests,
so I'll embed *mgo.Session and store the database name.
*/
type DatabaseSession struct {
	*mgo.Session
	databaseName string
}

/*
Connect to the local MongoDB and set up the database.
*/
func NewSession(name string) *DatabaseSession {
	// session, err := mgo.Dial("mongodb://root:root@ds057224.mongolab.com:57224/requestfiles")
	session, err := mgo.Dial("mongodb://localhost:27017/requestfiles")
	if err != nil {
		panic(err)
	}

	addIndexToRequestedFiles(session.DB(name))
	return &DatabaseSession{session, name}
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
