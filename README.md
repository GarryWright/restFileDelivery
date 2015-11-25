GoLang S3 File API
======================
Travis
[![Build Status](https://travis-ci.org/GarryWright/restfiledelivery.svg?branch=master)](https://travis-ci.org/GarryWright/restfiledelivery)
Drone
[![Build Status](https://drone.io/github.com/GarryWright/restFileDelivery/status.png)](https://drone.io/github.com/GarryWright/restFileDelivery/latest)
Coverall
[![Coverage Status](https://coveralls.io/repos/GarryWright/restFileDelivery/badge.svg?branch=master&service=github)](https://coveralls.io/github/GarryWright/restFileDelivery?branch=master)


The example is coded in golang and use mongod as its store. It is hosted on git and is continously built on 
Travis, Drone
Its uses "coverall" to montior test coverage
and "dockerhub" to store docker containers
and shippable to auto deploy the docker containers to aws

```
Add a document to the DB
curl -i -X POST -H "Content-Type: application/json" -d '{"clientid": "HSBC", "requestid": "00005", "ricdays": 21, "fileurl": "http://s3-us-west-2.amazonaws.com/garrysbucket/rics2.txt"}' localhost:3000/files


HTTP/1.1 201 Created
Content-Type: application/json; charset=UTF-8
Date: Tue, 24 Nov 2015 14:33:59 GMT
Content-Length: 138

{
  "clientid": "HSBC",
  "requestid": "00005",
  "ricdays": 21,
  "fileurl": "http://s3-us-west-2.amazonaws.com/garrysbucket/rics2.txt"
```
http://localhost:3000/requestedFiles :- returns all documents
```
curl http://localhost:3000/requestedFiles
[
  {
    "clientid": "HSBC",
    "requestid": "00002",
    "ricdays": 21,
    "fileurl": "https://s3-us-west-2.amazonaws.com/garrysbucket/rics.txt"
  },
  {
    "clientid": "HSBC",
    "requestid": "00003",
    "ricdays": 21,
    "fileurl": "http://s3-us-west-2.amazonaws.com/garrysbucket/rics2.txt"
  },
  {
    "clientid": "HSBC",
    "requestid": "00006",
    "ricdays": 21,
    "fileurl": "http://s3-us-west-2.amazonaws.com/garrysbucket/rics2.txt",
    "filebucket": "garrysbucket",
    "filekey": "rics.txt"
  }
```

curl http://localhost:3000/file?requestid=00002
returns the contents of the s3 file found in the requestedFile document to the session
```
IBM.N 1/1/1960 12.375
IBM.N 1/2/1960 12.375
IBM.N 1/3/1960 12.375
IBM.N 1/4/1960 12.375
IBM.N 1/5/1960 12.375
IBM.N 1/6/1960 12.375

```

localhost:3000/file?requestid=00002&destination=copiedfile.txt
returns the contents of the s3 file found in the requestedFile document to the file specified
```
{
  "done": "http://s3-us-west-2.amazonaws.com/garrysbucket/rics.txt downloaded to xxx.txt [132] bytes"
}
```
