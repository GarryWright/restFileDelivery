package main_test

import (
	. "github.com/GarryWright/restFileDelivery/fileDelivery"
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
			factory["FileURL"] = "https://s3-us-west-2.amazonaws.com/garrysbucket/rics.txt"
		})

	gory.Define("requestedFileNoRequest", RequestedFile{},
		func(factory gory.Factory) {
			factory["ClientId"] = "HSBC"
			factory["RicDays"] = 27
			factory["FileURL"] = "https://s3-us-west-2.amazonaws.com/garrysbucket/rics.txt"
		})

	gory.Define("requestedFile0", RequestedFile{},
		func(factory gory.Factory) {
			factory["ClientId"] = "HSBC"
			factory["RicDays"] = 27
			factory["RequestId"] = "0"
			factory["FileURL"] = "https://s3-us-west-2.amazonaws.com/garrysbucket/rics.txt"
		})
	gory.Define("requestedFile1", RequestedFile{},
		func(factory gory.Factory) {
			factory["ClientId"] = "HSBC"
			factory["RicDays"] = 27
			factory["RequestId"] = "1"
			factory["FileURL"] = "https://s3-us-west-2.amazonaws.com/garrysbucket/reeeecs.txt"
		})
}
