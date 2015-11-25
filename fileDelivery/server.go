package fileDelivery

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/binding"
	"github.com/martini-contrib/render"
	"gopkg.in/mgo.v2"
	"net/http"
	"os"
	"strconv"
)

/*
Wrap the Martini server struct.
*/
//type Server *martini.ClassicMartini

/*
Create a new *martini.ClassicMartini server.
We'll use a JSON renderer and our MongoDB
database handler. We define two routes:
"GET /file" and "POST /file".
*/

func NewServer(msession *DatabaseSession) *martini.ClassicMartini {
	// Create the server and set up middleware.
	m := martini.Classic()
	m.Use(render.Renderer(render.Options{
		IndentJSON: true,
	}))
	m.Use(msession.Database())

	// Get the first 20 filerequest documnents
	m.Get("/requestedFiles", func(r render.Render, db *mgo.Database, w http.ResponseWriter, t *http.Request) {
		requestedFiles, err := fetchAllRequestedFiles(db)
		if err != nil {
			r.JSON(401, map[string]string{
				"DB Error": err.Error(),
			})
		} else {
			r.JSON(200, requestedFiles)
		}
	})

	//get the s3 file from the url stored in the rqeuested file document
	m.Get("/file", func(r render.Render, db *mgo.Database, w http.ResponseWriter, t *http.Request) {

		if err := t.ParseForm(); err != nil {
			r.JSON(400, map[string]string{
				"error": "Failed to get parms",
			})
		}
		requestid := t.URL.Query().Get("requestid")
		destinationFile := t.URL.Query().Get("destination")
		if len(destinationFile) < 1 || len(requestid) < 1 {
			r.JSON(406, map[string]string{
				"Error": "destination file (1), requestid (2)  must be specified: [(1):" + destinationFile + " (2):" + requestid + "]",
			})
		} else {
			requestedFile, err := fetchRequestedFile(db, requestid)
			if err != nil {
				r.JSON(401, map[string]string{
					" DB Error": err.Error(),
				})
			}
			url := requestedFile.FileURL
			fbucket := requestedFile.FileBucket
			fkey := requestedFile.FileKey
			if len(fbucket) < 1 || len(fkey) < 1 {
				r.JSON(408, map[string]string{
					"Error": "source bucket (1) and key (2) must be specified in requestfile document [(1):" + fbucket + " (2):" + fkey + "]",
				})
			} else {

				filewriter, err := os.Create(destinationFile)
				defer filewriter.Close()
				if err != nil {
					r.JSON(404, map[string]error{
						"File error": err,
					})
				} else {
					var s3sess = session.New(aws.NewConfig().WithCredentials(credentials.AnonymousCredentials).WithRegion("us-west-2"))
					downloader := s3manager.NewDownloader(s3sess)
					numBytes, err := downloader.Download(filewriter, &s3.GetObjectInput{
						Bucket: &fbucket,
						Key:    &fkey,
					})
					//if _, err := io.Copy(filewriter, rr); err != nil {
					if err != nil {
						r.JSON(405, map[string]string{
							"error": "Failed to download to file",
						})
					} else {

						r.JSON(200, map[string]string{
							"done": url + " downloaded to " + destinationFile + " [" + strconv.FormatInt(numBytes, 10) + "] bytes",
						})
					}
				}
			}

		}
		//result.Close()
	})

	// Define the "POST /files" route. i.e pst requestFile documents
	m.Post("/files", binding.Json(RequestedFile{}),
		func(requestedFile RequestedFile,
			r render.Render,
			db *mgo.Database) {

			if requestedFile.valid() {
				// signature is valid, insert into database
				err := addRequestedFile(db, requestedFile)
				if err == nil {
					// insert successful, 201 Created
					r.JSON(201, requestedFile)
				} else {
					// insert failed, 400 Bad Request
					r.JSON(400, map[string]string{
						"error": err.Error(),
					})
				}
			} else {
				// signature is invalid, 400 Bad Request
				r.JSON(400, map[string]string{
					"error": "Not a valid requestedFile",
				})
			}
		})

	// Return the server. Call Run() on the server to
	// begin listening for HTTP requests.
	return m
}
