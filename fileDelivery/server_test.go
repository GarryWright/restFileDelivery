package fileDelivery_test

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	. "github.com/GarryWright/restFileDelivery/fileDelivery"
	"github.com/go-martini/martini"
	"github.com/modocache/gory"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
)

/*
Convert JSON data into a slice.
*/
func sliceFromJSON(data []byte) []interface{} {
	var result interface{}
	json.Unmarshal(data, &result)
	return result.([]interface{})
}

/*
Convert JSON data into a map.
*/
func mapFromJSON(data []byte) map[string]interface{} {
	var result interface{}
	json.Unmarshal(data, &result)
	return result.(map[string]interface{})
}

/*
Server unit tests.
*/
var _ = Describe("Server", func() {
	var dbName string
	var session *DatabaseSession
	var server *martini.ClassicMartini
	var request *http.Request
	var recorder *httptest.ResponseRecorder

	BeforeEach(func() {
		// Set up a new server, connected to a test database,
		// before each test.
		dbName = "requestedFiles"
		session = NewSession(dbName)
		server = NewServer(session)

		// Record HTTP responses.
		recorder = httptest.NewRecorder()
	})

	AfterEach(func() {
		// Clear the database after each test.
		session.DB(dbName).DropDatabase()
	})

	Describe("GET /files", func() {

		// Set up a new GET request before every test
		// in this describe block.
		BeforeEach(func() {
			request, _ = http.NewRequest("GET", "/requestedFiles", nil)
		})

		Context("when no requestedFiles exist", func() {
			It("returns a status code of 200", func() {
				server.ServeHTTP(recorder, request)
				Expect(recorder.Code).To(Equal(200))
			})

			It("returns a empty body", func() {
				server.ServeHTTP(recorder, request)
				Expect(recorder.Body.String()).To(Equal("[]"))
			})
		})

		Context("when requestedFiles exist", func() {

			// Insert two valid requestedFiles into the database
			// before each test in this context.
			BeforeEach(func() {
				collection := session.DB(dbName).C("requestedFiles")
				collection.Insert(gory.Build("requestedFile"))
				collection.Insert(gory.Build("requestedFile"))

			})

			It("returns a status code of 200", func() {
				server.ServeHTTP(recorder, request)
				Expect(recorder.Code).To(Equal(200))
			})

			It("returns those requestedFiles in the body", func() {
				server.ServeHTTP(recorder, request)
				peopleJSON := sliceFromJSON(recorder.Body.Bytes())
				fmt.Println(recorder)
				Expect(len(peopleJSON)).To(Equal(2))

				personJSON := peopleJSON[0].(map[string]interface{})
				Expect(personJSON["clientid"]).To(Equal("HSBC"))
				Expect(personJSON["requestid"]).ShouldNot(Equal("A"))
				Expect(personJSON["ricdays"]).To(Equal(float64(27)))
				Expect(personJSON["fileurl"]).To(Equal("http://s3-us-west-2.amazonaws.com/garrysbucket/rics.txt"))

			})
		})
	})

	Describe("POST /files", func() {

		Context("with invalid JSON", func() {

			// Create a POST request using JSON from our invalid
			// factory object before each test in this context.
			BeforeEach(func() {
				body, _ := json.Marshal(
					gory.Build("requestedFileNoRequest"))
				request, _ = http.NewRequest(
					"POST", "/files", bytes.NewReader(body))
			})

			It("returns a status code of 400", func() {
				server.ServeHTTP(recorder, request)
				Expect(recorder.Code).To(Equal(400))
			})
		})

		Context("with valid JSON", func() {

			// Create a POST request with valid JSON from
			// our factory before each test in this context.
			BeforeEach(func() {
				body, _ := json.Marshal(
					gory.Build("requestedFile0"))
				request, _ = http.NewRequest(
					"POST", "/files", bytes.NewReader(body))
			})

			It("returns a status code of 201", func() {
				server.ServeHTTP(recorder, request)
				Expect(recorder.Code).To(Equal(201))
			})

			It("returns the inserted requestedfile", func() {
				server.ServeHTTP(recorder, request)

				personJSON := mapFromJSON(recorder.Body.Bytes())
				Expect(personJSON["clientid"]).To(Equal("HSBC"))
				Expect(personJSON["requestid"]).To(Equal("0"))
				Expect(personJSON["ricdays"]).To(Equal(float64(27)))
				Expect(personJSON["fileurl"]).To(Equal("http://s3-us-west-2.amazonaws.com/garrysbucket/rics.txt"))
			})
		})

		Context("with JSON containing a duplicate request", func() {
			BeforeEach(func() {
				requestedFile := gory.Build("requestedFile")
				session.DB(dbName).C("requestedFiles").Insert(requestedFile)

				body, _ := json.Marshal(requestedFile)
				request, _ = http.NewRequest(
					"POST", "/files", bytes.NewReader(body))
			})

			It("returns a status code of 400", func() {
				server.ServeHTTP(recorder, request)
				Expect(recorder.Code).To(Equal(400))
			})
		})
	})

	Describe("POST /files", func() {

		Context("with invalid JSON", func() {

			// Create a POST request using JSON from our invalid
			// factory object before each test in this context.
			BeforeEach(func() {
				body, _ := json.Marshal(
					gory.Build("requestedFileNoRequest"))
				request, _ = http.NewRequest(
					"POST", "/files", bytes.NewReader(body))
			})

			It("returns a status code of 400", func() {
				server.ServeHTTP(recorder, request)
				Expect(recorder.Code).To(Equal(400))
			})
		})

		Context("with valid JSON", func() {

			// Create a POST request with valid JSON from
			// our factory before each test in this context.
			BeforeEach(func() {
				body, _ := json.Marshal(
					gory.Build("requestedFile"))
				request, _ = http.NewRequest(
					"POST", "/files", bytes.NewReader(body))
			})

			It("returns a status code of 201", func() {
				server.ServeHTTP(recorder, request)
				Expect(recorder.Code).To(Equal(201))
			})

			It("returns the inserted requestedfile", func() {
				server.ServeHTTP(recorder, request)

				personJSON := mapFromJSON(recorder.Body.Bytes())
				Expect(personJSON["clientid"]).To(Equal("HSBC"))
				Expect(personJSON["requestid"]).ShouldNot(Equal("A"))
				Expect(personJSON["ricdays"]).To(Equal(float64(27)))
				Expect(personJSON["fileurl"]).To(Equal("http://s3-us-west-2.amazonaws.com/garrysbucket/rics.txt"))
			})
		})

		Context("with JSON containing a duplicate request", func() {
			BeforeEach(func() {
				requestedFile := gory.Build("requestedFile")
				session.DB(dbName).C("requestedFiles").Insert(requestedFile)

				body, _ := json.Marshal(requestedFile)
				request, _ = http.NewRequest(
					"POST", "/files", bytes.NewReader(body))
			})

			It("returns a status code of 400", func() {
				server.ServeHTTP(recorder, request)
				Expect(recorder.Code).To(Equal(400))
			})
		})
	})

	Describe("GET /file?requestid=0&destination=copiedfile.txt", func() {

		Context("with invalid request id", func() {

			BeforeEach(func() {
				collection := session.DB(dbName).C("requestedFiles")
				collection.Insert(gory.Build("requestedFile1"))
				request, _ = http.NewRequest("GET", "/file?requestid=0&destination=copiedfile.txt", nil)

			})

			Context("when a requestedFiles does not exists with request id 0", func() {
				It("returns a status code of 401", func() {
					server.ServeHTTP(recorder, request)
					Expect(recorder.Code).To(Equal(401))
				})

			})

		})
		Context("with valid token and a bad file url", func() {

			BeforeEach(func() {
				collection := session.DB(dbName).C("requestedFiles")
				collection.Insert(gory.Build("requestedFile1"))
				request, _ = http.NewRequest("GET", "/file?requestid=1&destination=copiedfile.txt", nil)

			})

			Context("when a s3 file does not exits", func() {
				It("returns a status code of 402", func() {
					server.ServeHTTP(recorder, request)
					Expect(recorder.Code).To(Equal(405))
				})

			})

		})

	})
	Describe("GET /file?requestid=0&destination=copiedfile.txt", func() {

		Context("with valid request ", func() {

			BeforeEach(func() {
				collection := session.DB(dbName).C("requestedFiles")
				collection.Insert(gory.Build("requestedFile0"))
				request, _ = http.NewRequest("GET", "/file?requestid=0&destination=copiedfile.txt", nil)

			})

			Context("when a requestedFiles exists with request id 0", func() {
				It("returns a status code of 200", func() {
					server.ServeHTTP(recorder, request)
					Expect(recorder.Code).To(Equal(200))
					Expect(strings.Contains(recorder.Body.String(), "[132] bytes")).To(BeTrue())
				})

				It("returns a the s3 file contents into a file", func() {
					os.Remove("copiedfile.txt")
					server.ServeHTTP(recorder, request)
					file, _ := os.Open("copiedfile.txt")
					defer file.Close()

					var lines []string
					scanner := bufio.NewScanner(file)
					for scanner.Scan() {
						lines = append(lines, scanner.Text())
					}
					Expect(lines[0]).To(Equal("IBM.N 1/1/1960 12.375"))
				})
			})

		})
		Context("with invalid parms ", func() {

			BeforeEach(func() {
				collection := session.DB(dbName).C("requestedFiles")
				collection.Insert(gory.Build("requestedFile0"))
				request, _ = http.NewRequest("GET", "/file?requestid=0", nil)

			})

			Context("when destination is missing", func() {
				It("returns a status code of 406", func() {
					server.ServeHTTP(recorder, request)
					Expect(recorder.Code).To(Equal(406))
					Expect(strings.Contains(recorder.Body.String(), "[(1): (2):")).To(BeTrue())
				})

			})

		})
		Context("with invalid parms ", func() {

			BeforeEach(func() {
				collection := session.DB(dbName).C("requestedFiles")
				collection.Insert(gory.Build("requestedFile0"))
				request, _ = http.NewRequest("GET", "/file?destination=copiedfile.txt", nil)

			})

			Context("when requestid is missing", func() {
				It("returns a status code of 406", func() {
					server.ServeHTTP(recorder, request)
					Expect(recorder.Code).To(Equal(406))
					Expect(strings.Contains(recorder.Body.String(), "(2):]")).To(BeTrue())
				})

			})

		})
		Context("with invalid requestfile document ", func() {

			BeforeEach(func() {
				collection := session.DB(dbName).C("requestedFiles")
				collection.Insert(gory.Build("requestedFile3"))
				request, _ = http.NewRequest("GET", "/file?requestid=1&destination=copiedfile.txt", nil)

			})

			Context("when s3 key is missing", func() {
				It("returns a status code of 408", func() {
					server.ServeHTTP(recorder, request)
					Expect(recorder.Code).To(Equal(408))
					Expect(strings.Contains(recorder.Body.String(), "(2):]")).To(BeTrue())
				})

			})

		})
		Context("with invalid requestfile document ", func() {

			BeforeEach(func() {
				collection := session.DB(dbName).C("requestedFiles")
				collection.Insert(gory.Build("requestedFile2"))
				request, _ = http.NewRequest("GET", "/file?requestid=1&destination=copiedfile.txt", nil)

			})

			Context("when s3 bucket is missing", func() {
				It("returns a status code of 408", func() {
					server.ServeHTTP(recorder, request)
					Expect(recorder.Code).To(Equal(408))
					Expect(strings.Contains(recorder.Body.String(), "(1): (2)")).To(BeTrue())
				})

			})

		})

	})
})
