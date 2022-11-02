package main

import (
	"github.com/dgrijalva/jwt-go"
	eventbus "github.com/snowsoftwareglobal/platform-go-eventbus"
	"io/ioutil"
	"log"
	"time"
)

const (
	handleFileUploadSubject = "command.filemanager.file.upload"
	handleFileUploadType    = eventbus.DefaultTypePrefix + "filemanager.file.upload"
)

func uploadbus(url, bearer, mediaFilename string, parsedClaims jwt.MapClaims) {

	//read file
	body, err := ioutil.ReadFile(mediaFilename)

	if err != nil {
		log.Fatalf("unable to read file: %v", err)
	}

	//construct event
	ebus, err := eventbus.New(
		eventbus.WithServer("localhost", 4222),
		eventbus.WithClusterID("test-cluster"),
		eventbus.WithClientID("client_1"),
	)

	if err != nil {
		log.Fatal(err)
	}

	conn, err := ebus.NewConnection()

	if err != nil {
		log.Fatal(err)
	}
	defer conn.Disconnect()
	event := eventbus.NewEvent("source", handleFileUploadType, body)

	metadata := SaveFile{
		FileName:   "report.pdf",
		Mime:       "application/pdf",
		DeleteMark: false,
		TenantId:   "fe8c651b-18d6-4912-6982-08dab7e9b136",
		User:       "15f6489a-2745-45be-9037-fa972bdd3b01",
		Size:       int64(len(body)),
		Category:   "report",
		FileType:   "pdf",
		App:        "saas",
		Tags: map[string]string{
			"importance": "Medium",
		},
	}

	//File metadata
	event.Extensions["metadata"] = metadata
	// Make a request with a timeout
	_, err = conn.Request(handleFileUploadSubject, event, 5*time.Second)
	if err != nil {
		log.Fatal(err)
	}
}
