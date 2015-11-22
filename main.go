package main

import (
	"github.com/GarryWright/restfiledelivery/Godeps/_workspace/src/github.com/GarryWright/restFileDelivery/fileDelivery"
	"os"
)

/*
Create a new MongoDB session, using a database
Create a new server using
that session, then begin listening for HTTP requests.
*/
func main() {

	db_name := os.Getenv("DB_NAME")

	if db_name == "" {
		db_name = "requestfiles"
	}
	session := fileDelivery.NewSession(db_name)
	server := fileDelivery.NewServer(session)
	server.Run()
}
