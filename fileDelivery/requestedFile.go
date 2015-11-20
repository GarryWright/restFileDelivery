package fileDelivery

import "gopkg.in/mgo.v2"

/*
Each signature is composed of a first name, last name,
email, age, and short message. When represented in
JSON, ditch TitleCase for snake_case.
*/
type RequestedFile struct {
	ClientId  string `json:"clientid"`
	RequestId string `json:"requestid"`
	RicDays   int    `json:"ricdays"`
	FileURL   string `json:"fileurl"`
}

/*
I want to make sure all these fields are present. The message
is optional, but if it's present it has to be less than
140 characters--it's a short blurb, not your life story.
*/
func (requestedFile *RequestedFile) valid() bool {
	return len(requestedFile.ClientId) > 0 &&
		len(requestedFile.RequestId) > 0 &&
		len(requestedFile.FileURL) > 0 &&
		requestedFile.RicDays >= 1
}

/*
I'll use this method when displaying all signatures for
"GET /signatures". Consult the mgo docs for more info:
http://godoc.org/labix.org/v2/mgo
*/
func fetchAllRequestedFiles(db *mgo.Database) []RequestedFile {
	requestedFiles := []RequestedFile{}
	err := db.C("requestedFiles").Find(nil).All(&requestedFiles)
	if err != nil {
		panic(err)
	}

	return requestedFiles
}
