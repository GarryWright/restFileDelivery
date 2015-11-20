package main

import "github.com/GarryWright/restFileDelivery/fileDelivery"

/*
Create a new MongoDB session, using a database
named "signatures". Create a new server using
that session, then begin listening for HTTP requests.
*/
func main() {
	session := fileDelivery.NewSession("requestfiles")
	server := fileDelivery.NewServer(session)
	server.Run()
}
