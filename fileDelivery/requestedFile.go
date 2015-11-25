package fileDelivery

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type RequestedFile struct {
	ClientId   string `json:"clientid"`
	RequestId  string `json:"requestid"`
	RicDays    int    `json:"ricdays"`
	FileURL    string `json:"fileurl"`
	FileBucket string `json:"filebucket"`
	FileKey    string `json:"filekey"`
}

func (requestedFile *RequestedFile) valid() bool {
	return len(requestedFile.ClientId) > 0 &&
		len(requestedFile.RequestId) > 0 &&
		len(requestedFile.FileURL) > 0 &&
		len(requestedFile.FileBucket) > 0 &&
		len(requestedFile.FileKey) > 0 &&
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
