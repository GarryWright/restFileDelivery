package fileDelivery

import (
	// "github.com/aws/aws-sdk-go/aws"
	// "github.com/aws/aws-sdk-go/aws/session"
	// "github.com/aws/aws-sdk-go/service/s3"
	"github.com/go-martini/martini"
	"github.com/kr/s3/s3util"
	"github.com/martini-contrib/binding"
	"github.com/martini-contrib/render"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"io"
	"net/http"
)

/*
Wrap the Martini server struct.
*/
//type Server *martini.ClassicMartini

/*
Create a new *martini.ClassicMartini server.
We'll use a JSON renderer and our MongoDB
database handler. We define two routes:
"GET /signatures" and "POST /signatures".
*/
// func NewServer() *martini.ClassicMartini {
func NewServer(msession *DatabaseSession) *martini.ClassicMartini {
	// Create the server and set up middleware.
	m := martini.Classic()
	m.Use(render.Renderer(render.Options{
		IndentJSON: true,
	}))
	m.Use(msession.Database())

	m.Get("/requestedFiles", func(r render.Render, db *mgo.Database, w http.ResponseWriter, t *http.Request) {
		requestedFiles := []RequestedFile{}
		err := db.C("requestedFiles").Find(nil).All(&requestedFiles)
		if err != nil {
			panic(err)
		}
		r.JSON(200, requestedFiles)
	})

	m.Get("/file", func(r render.Render, db *mgo.Database, w http.ResponseWriter, t *http.Request) {
		//r.JSON(200, sendRequestedFile(db))
		//fp := path.Join("images", "foo.png")
		//http.ServeFile(w, r, fp)
		if err := t.ParseForm(); err != nil {
			r.JSON(400, map[string]string{
				"error": "Failed to get parms",
			})
		}
		token := t.Form.Get("req")
		// r.JSON(400, map[string]string{
		// 	"value": token,
		// })
		f := RequestedFile{}
		err := db.C("requestedFiles").Find(bson.M{"requestid": token}).One(&f)
		if err == nil {
			// insert successful, 201 Created
			// r.JSON(201, f)
		} else {
			// get failed, 400 Bad Request
			r.JSON(401, map[string]string{
				"error": err.Error(),
			})
		}
		url := f.FileURL
		// url := "https://s3-us-west-2.amazonaws.com/garrysbucket/rics.txt"
		s3util.DefaultConfig.AccessKey = "AKIAJFFIJIAWUKP3NQMA"
		s3util.DefaultConfig.SecretKey = "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"

		result, err := s3util.Open(url, nil)

		//svc := s3.New(session.New())
		// result, err := svc.GetObject(&s3.GetObjectInput{
		// 	Bucket: aws.String("garysbucket"),
		// 	Key:    aws.String("rics.txt"),
		// })
		//result, err := svc.GetObject(url)
		if err != nil {
			r.JSON(402, map[string]string{
				"error": "Failed to get object",
			})
		}

		if _, err := io.Copy(w, result); err != nil {
			r.JSON(403, map[string]string{
				"error": "Failed to download object",
			})
		}
		result.Close()
	})

	// Define the "POST /signatures" route.
	m.Post("/files", binding.Json(RequestedFile{}),
		func(requestedFile RequestedFile,
			r render.Render,
			db *mgo.Database) {

			if requestedFile.valid() {
				// signature is valid, insert into database
				err := db.C("requestedFiles").Insert(requestedFile)
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
