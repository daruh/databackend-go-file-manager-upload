package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/ddosify/go-faker/faker"
	"github.com/dgrijalva/jwt-go"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"path/filepath"
	"time"
)

type SnowClaims struct {
	tenant string `json:"tenant"`
	jwt.StandardClaims
}

func main() {

	env := "dev"
	flag.Parse()

	positionalArgs := flag.Args()
	if len(positionalArgs) == 0 {
		log.Fatalf("This program requires at least 1 positional argument.")
	}

	// New multipart writer.
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	var size int64
	// Media Files.
	for _, mediaFilename := range positionalArgs {
		mediaData, errRead := ioutil.ReadFile(mediaFilename)
		if errRead != nil {
			log.Fatalf("Error reading media file: %v", errRead)
		}
		size = int64(len(mediaData))

		filename := filepath.Base(mediaFilename)

		mediaHeader := textproto.MIMEHeader{}
		mediaHeader.Set("Content-Disposition", fmt.Sprintf("form-data; name=\"file\"; filename=\"%v\"", filename))
		mediaHeader.Set("Content-ID", "media")
		mediaHeader.Set("Content-Type", "application/octet-stream")
		mediaHeader.Set("Content-Filename", filename)

		mediaPart, err := writer.CreatePart(mediaHeader)
		if err != nil {
			log.Fatalf("Error writing media headers: %v", errRead)
		}

		if _, err := io.Copy(mediaPart, bytes.NewReader(mediaData)); err != nil {
			log.Fatalf("Error writing media: %v", errRead)
		}
	}

	bearer := "eyJhbGciOiJSUzI1NiIsImtpZCI6IkRFNDBBMjU5MjRENDk1MERGMTkxNkI3RjlDMkJEOThBNkE4NzE1OEMiLCJ4NXQiOiIza0NpV1NUVWxRM3hrV3RfbkN2WmltcUhGWXciLCJ0eXAiOiJhdCtqd3QifQ.eyJhdXRoX3RpbWUiOjE2NjQ0MzYxOTQsInN1YiI6ImFlODYyNWE0LTRiZDAtNDc4Zi05ZTdiLTA4ZGEzMTg1ODdiNSIsInRlbmFudCI6ImNkYWZlOGRkLTlkZDItNDY0Ny0zNDM1LTA4ZGEzMTg1ODc5MCIsInJlZ2lvbiI6Im5vcnRoZXVyb3BlIiwiYWNjIjoiMDk2YjgxZTUtOGRlMS00NDI2LTI0ZDYtMDhkYTMxODU4N2I4IiwiZW1wbF90eXBlIjoicm5kIiwicm9sZSI6WyJzYWxlcy5tYW5hZ2VyIiwiZmlsZW1hbmFnZXIucm9sZSIsImlkLmlkcC5hZG1pbiJdLCJvaV9wcnN0IjoiMTVmNjQ4OWEtMjc0NS00NWJlLTkwMzctZmE5NzJiZGQzYjAxIiwiY2xpZW50X2lkIjoiMTVmNjQ4OWEtMjc0NS00NWJlLTkwMzctZmE5NzJiZGQzYjAxIiwidG9rX2lkIjoiMDcwMTFjY2UtNjhiMC00NGIyLThkMGEtYzE2NjIyNzU4ZGI5IiwiYXVkIjpbImFwaTovL3Nub3dzb2Z0d2FyZS5pby9hcGkiLCJhcGk6Ly9zbm93c29mdHdhcmUuaW8vaWRwL2FwaSJdLCJzY29wZSI6ImxpY2Vuc2luZy5saWNlbnNlLmNydWQgZWRnZS5pbnYuYXBpLmNydWQgaWQudGVuYW50cy5jcnVkIGZpbGVtYW5hZ2VyLnIgZmlsZW1hbmFnZXIuY3J1ZCBpZC51c2Vycy5jcnVkIGlkLmFwcHMuY3J1ZCBpZC5jbGFpbXMuY3J1ZCBpZC5vcmdhbml6YXRpb25hbHVuaXRzLmNydWQgaWQuc3lzLnRyc2MgaWQuc3lzLmNydWQgaWQubG9naW5wcm92aWRlcnMuY3J1ZCBpZC5yb2xlcy5jcnVkIG9wZW5pZCIsImV4cCI6MTY2NDQzOTc5NCwiaXNzIjoiaHR0cDovL2xvY2FsaG9zdDo4MTE4L2lkcCIsImlhdCI6MTY2NDQzNjE5NH0.ZA5-8fxPy4LdIbVhIq7scFWJHvVvqzdG9LNSezN1eYLUM28jWMXPHuwKBYHNdY-P2wSAiGu-9K9xlPDHBA8TCUpdhgEVLCTaX8xBkvquvctvDk_fwIamOdYpbZj5jkhxyA4ijQjDIu1DN79kLhAoR3cRtgKfSJmKEBugVOi0TWxAZGR9lwrIOf2iJk_mRlaiwoeJ1-2Ok6byvz2pgJF_4AR20dZYOrj_ux8pEvPJVEGD_eiLqlivC28U4sXFS4myda3RZMiYUksn-GgrLlwpgLGX-GbCHc8bwDEwjvyPqARMN3OIfePwUW_ZcQ_bPkqSBpL2pu0zSBcF1kOpJkTzQA"

	var claims SnowClaims
	token, err := jwt.ParseWithClaims(bearer, claims, func(token *jwt.Token) (interface{}, error) {
		return nil, nil
	})

	if err != nil {
		panic(err)
	}

	parsedClaims := token.Claims.(*SnowClaims)

	faker := faker.NewFaker()
	metadata := SaveFile{
		FileName:   faker.RandomFileName() + ".pdf",
		Mime:       "application/pdf",
		DeleteMark: false,
		TenantId:   parsedClaims.tenant,
		User:       parsedClaims.Audience,
		Size:       size,
		Category:   faker.RandomProductAdjective(),
		FileType:   faker.RandomFileType(),
		App:        faker.RandomProductName(),
		Tags: map[string]string{
			"importance":  "Medium",
			"requestedBy": "Some manager",
		},
	}

	// Metadata part.
	metadataHeader := textproto.MIMEHeader{}
	metadataHeader.Set("Content-Type", "application/json")
	metadataHeader.Set("Content-Disposition", "form-data; name=\"metadata\"")
	metadataHeader.Set("Content-ID", "metadata")
	part, err := writer.CreatePart(metadataHeader)
	if err != nil {
		log.Fatalf("Error writing metadata headers: %v", err)
	}

	marshal, _ := json.Marshal(metadata)
	part.Write(marshal)

	// Close multipart writer.
	if err := writer.Close(); err != nil {
		log.Fatalf("Error closing multipart writer: %v", err)
	}

	// Request Content-Type with boundary parameter.
	contentType := fmt.Sprintf("multipart/form-data; boundary=%s", writer.Boundary())

	// Initialize HTTP Request and headers.

	urls := map[string]string{}

	urls["prod"] = "https://westeurope.dev-snowsoftware.io"
	urls["dev"] = "http://localhost:8008"

	uploadURL := urls[env] + "/api/filemanager/tenants/" + parsedClaims.tenant + "/upload"

	r, err := http.NewRequest(http.MethodPost, uploadURL, bytes.NewReader(body.Bytes()))
	if err != nil {
		log.Fatalf("Error initializing a request: %v", err)
	}

	r.Header.Set("Authorization", "Bearer "+bearer)
	r.Header.Set("Content-Type", contentType)
	r.Header.Set("Accept", "*/*")

	// HTTP Client.
	client := &http.Client{Timeout: 180 * time.Second}
	rsp, err := client.Do(r)
	if err != nil {
		log.Fatalf("Error making a request: %v", err)
	}

	// Check response status code.
	if rsp.StatusCode != http.StatusCreated {
		log.Printf("Request failed with response code: %d", rsp.StatusCode)
	} else {
		log.Printf("Request was a success ")
	}
}
