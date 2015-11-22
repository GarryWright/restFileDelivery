package fileDelivery

import (
	"github.com/go-martini/martini"
	// "github.com/kr/s3/s3util"
	"bytes"
	"github.com/martini-contrib/binding"
	"github.com/martini-contrib/render"
	"gopkg.in/mgo.v2"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
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

		requestedFile, err := fetchRequestedFile(db, requestid)

		if err != nil {
			r.JSON(401, map[string]string{
				" DB Error": err.Error(),
			})
		}
		url := requestedFile.FileURL
		// url := "https://s3-us-west-2.amazonaws.com/garrysbucket/rics.txt"
		// s3Key := os.Getenv("S3KEY")
		// s3SecretKey := os.Getenv("S3SECRET")

		response, err := http.Get(url)

		defer response.Body.Close()
		contents, err := ioutil.ReadAll(response.Body)

		if err != nil || strings.Contains(string(contents[:]), ">The specified key does not exist.<") {
			r.JSON(402, map[string]string{
				"error": "Failed to get object",
			})
		}
		rr := bytes.NewReader(contents)

		// result, err := s3util.Open(url, nil)

		if destinationFile == "" {
			// Copy the s3 file contents to the http response writer]
			if _, err := io.Copy(w, rr); err != nil {
				r.JSON(403, map[string]string{
					"error": "Failed to download object",
				})
			}
		} else {
			filewriter, err := os.Create(destinationFile)
			if err != nil {
				r.JSON(404, map[string]error{
					"File error": err,
				})
			} else {
				if _, err := io.Copy(filewriter, rr); err != nil {
					r.JSON(405, map[string]string{
						"error": "Failed to download to file",
					})
				}
			}

			filewriter.Close()
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
