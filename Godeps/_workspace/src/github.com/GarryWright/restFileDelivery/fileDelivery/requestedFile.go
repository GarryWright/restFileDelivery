package fileDelivery

import (
	"github.com/GarryWright/restfiledelivery/Godeps/_workspace/src/gopkg.in/mgo.v2"
	"github.com/GarryWright/restfiledelivery/Godeps/_workspace/src/gopkg.in/mgo.v2/bson"
)

type RequestedFile struct {
	ClientId  string `json:"clientid"`
	RequestId string `json:"requestid"`
	RicDays   int    `json:"ricdays"`
	FileURL   string `json:"fileurl"`
}

func (requestedFile *RequestedFile) valid() bool {
	return len(requestedFile.ClientId) > 0 &&
		len(requestedFile.RequestId) > 0 &&
		len(requestedFile.FileURL) > 0 &&
		requestedFile.RicDays >= 1
}
func fetchAllRequestedFiles(db *mgo.Database) ([]RequestedFile, error) {
	requestedFiles := []RequestedFile{}
	err := db.C("requestedFiles").Find(nil).All(&requestedFiles)
	return requestedFiles, err
}

func fetchRequestedFile(db *mgo.Database, requestid string) (RequestedFile, error) {
	requestedFile := RequestedFile{}
	err := db.C("requestedFiles").Find(bson.M{"requestid": requestid}).One(&requestedFile)
	return requestedFile, err
}

func addRequestedFile(db *mgo.Database, requestedFile RequestedFile) error {
	err := db.C("requestedFiles").Insert(requestedFile)
	return err
}
