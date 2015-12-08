package fileDelivery_test

import (
	. "fileDelivery"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"fmt"
	"github.com/modocache/gory"
	"testing"
)

func TestRestFileDelivery(t *testing.T) {
	defineFactories()
	RegisterFailHandler(Fail)
	RunSpecs(t, "RestFileDelivery Suite")
}

func defineFactories() {
	gory.Define("requestedFile", RequestedFile{},
		func(factory gory.Factory) {
			factory["ClientId"] = "HSBC"
			factory["RequestId"] = gory.Sequence(
				func(n int) interface{} {
					return fmt.Sprintf("%d", n)
				})
			factory["RicDays"] = 27
			factory["FileURL"] = "http://s3-us-west-2.amazonaws.com/garrysbucket/rics.txt"
			factory["FileBucket"] = "garrysbucket"
			factory["FileKey"] = "rics.txt"
		})

	gory.Define("requestedFileNoRequest", RequestedFile{},
		func(factory gory.Factory) {
			factory["ClientId"] = "HSBC"
			factory["RicDays"] = 27
			factory["FileURL"] = "http://s3-us-west-2.amazonaws.com/garrysbucket/rics.txt"
			factory["FileBucket"] = "garrysbucket"
			factory["FileKey"] = "rics.txt"
		})

	gory.Define("requestedFile0", RequestedFile{},
		func(factory gory.Factory) {
			factory["ClientId"] = "HSBC"
			factory["RicDays"] = 27
			factory["RequestId"] = "0"
			factory["FileURL"] = "http://s3-us-west-2.amazonaws.com/garrysbucket/rics.txt"
			factory["FileBucket"] = "garrysbucket"
			factory["FileKey"] = "rics.txt"
		})
	gory.Define("requestedFile1", RequestedFile{},
		func(factory gory.Factory) {
			factory["ClientId"] = "HSBC"
			factory["RicDays"] = 27
			factory["RequestId"] = "1"
			factory["FileURL"] = "http://s3-us-west-2.amazonaws.com/garrysbucket/reeeecs.txt"
			factory["FileBucket"] = "garrysbucket"
			factory["FileKey"] = "reeeecs.txt"
		})
	gory.Define("requestedFile2", RequestedFile{},
		func(factory gory.Factory) {
			factory["ClientId"] = "HSBC"
			factory["RicDays"] = 27
			factory["RequestId"] = "1"
			factory["FileURL"] = "http://s3-us-west-2.amazonaws.com/garrysbucket/reeeecs.txt"
			factory["FileKey"] = "reeeecs.txt"
		})
	gory.Define("requestedFile3", RequestedFile{},
		func(factory gory.Factory) {
			factory["ClientId"] = "HSBC"
			factory["RicDays"] = 27
			factory["RequestId"] = "1"
			factory["FileURL"] = "http://s3-us-west-2.amazonaws.com/garrysbucket/reeeecs.txt"
			factory["FileBucket"] = "garrysbucket"
		})
}
