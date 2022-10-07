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

type SaveFile struct {
	FileName   string            `json:"fileName"`
	TenantId   string            `json:"tenantId"`
	Mime       string            `json:"mime"`
	Category   string            `json:"category"`
	FileType   string            `json:"fileType"`
	User       string            `json:"user"`
	DeleteMark bool              `bson:"deleteMark"`
	ExpiryTs   int64             `bson:"expiryTs"`
	App        string            `json:"app"`
	Size       int64             `json:"size"`
	Tags       map[string]string `json:"tags"`
}

func main() {

	env := "dev"
	flag.Parse()

	 positionalArgs := flag.Args()
	if len(positionalArgs) == 0 {
		log.Fatalf("This program requires at least 1 positional argument.")
	}
for _, mediaFilename := range positionalArgs {
	// New multipart writer.
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	var size int64
	// Media Files.
	
		mediaData, errRead := ioutil.ReadFile(mediaFilename)
		if errRead != nil {
			log.Fatalf("Error reading media file: %v", errRead)
		}
		size = int64(len(mediaData))

		
	

	bearer := "eyJhbGciOiJSUzI1NiIsImtpZCI6IkI3RTExQTVFNjA0MEU0M0QyNzhGRUYxQkJDOEEwNUNGNkJDRjQ4QUIiLCJ4NXQiOiJ0LUVhWG1CQTVEMG5qLThidklvRnoydlBTS3MiLCJ0eXAiOiJhdCtqd3QifQ.eyJhdXRoX3RpbWUiOjE2NjUxMzc3NTUsInN1YiI6IjEzNDE3OWY0LTJjM2ItNDRkOS0xZWFkLTA4ZGFhNzliYjRjYiIsInRlbmFudCI6ImI0MWUzNDU0LWExOTMtNGM0NC0xMGFlLTA4ZGFhNzliYjRhNSIsInJlZ2lvbiI6Im5vcnRoZXVyb3BlIiwiYWNjIjoiNjk0YTU0OTItMGRlOS00YWMxLWEwOTQtMDhkYWE3OWJiNGNkIiwiZW1wbF90eXBlIjoicm5kIiwicm9sZSI6WyJzYWxlcy5tYW5hZ2VyIiwiaWQuaWRwLmFkbWluIl0sIm9pX3Byc3QiOiIxNWY2NDg5YS0yNzQ1LTQ1YmUtOTAzNy1mYTk3MmJkZDNiMDEiLCJjbGllbnRfaWQiOiIxNWY2NDg5YS0yNzQ1LTQ1YmUtOTAzNy1mYTk3MmJkZDNiMDEiLCJ0b2tfaWQiOiIzZWFmZDRhMi00OWQ1LTQ2ZWUtYjYwZS03MWIxMzVjM2ZiNTIiLCJhdWQiOlsiYXBpOi8vc25vd3NvZnR3YXJlLmlvL2FwaSIsImFwaTovL3Nub3dzb2Z0d2FyZS5pby9pZHAvYXBpIl0sInNjb3BlIjoibGljZW5zaW5nLmxpY2Vuc2UuY3J1ZCBpZC50ZW5hbnRzLmNydWQgaWQudXNlcnMuY3J1ZCBpZC5hcHBzLmNydWQgaWQuY2xhaW1zLmNydWQgaWQub3JnYW5pemF0aW9uYWx1bml0cy5jcnVkIGlkLnN5cy50cnNjIGlkLnN5cy5jcnVkIGlkLmxvZ2lucHJvdmlkZXJzLmNydWQgaWQucm9sZXMuY3J1ZCBvcGVuaWQgb2ZmbGluZV9hY2Nlc3MiLCJleHAiOjE2NjUxNDEzNTUsImlzcyI6Imh0dHA6Ly9sb2NhbGhvc3Q6ODExOC9pZHAiLCJpYXQiOjE2NjUxMzc3NTV9.asSzK6Rs6Y8KUenRVim9bfXZ-18nqDE_3TaCzA_H31I0zM-osPeuY3ENMLnl_mj9nYtS3zG7gJpueaNnwmA8boC4NgIo4cEJbh7OdeqUk2e0CiUW3RtivaZv_YXb9xo3lPMaPwG_-745LgusLbD93RogPdp4u5iNwTTnZIItcjc275rK9Xk0iS1_QMHVp6FjQSZ18QdsD1dycSUgYFE-OLcnGBVSyY1oGSyHp7y_TvGzz9d7o6WA64qCkBAARLIl3zGCiVYGNXLNa-MwlZwXT02ARQZhGSizNpGJ5lKrLBAthPxVwSYjZX9URe-Y2K18U_pUMCORvCSl_zdf_RUuCg"

	token, _ := jwt.Parse(bearer, func(token *jwt.Token) (interface{}, error) {
		return nil, nil
	})

	parsedClaims := token.Claims.(jwt.MapClaims)

	faker := faker.NewFaker()
	metadata := SaveFile{
		FileName:   faker.RandomFileName() + ".pdf",
		Mime:       "application/pdf",
		DeleteMark: false,
		TenantId:   parsedClaims["tenant"].(string),
		User:       parsedClaims["sub"].(string),
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

	

	// Request Content-Type with boundary parameter.
	contentType := fmt.Sprintf("multipart/form-data; boundary=%s", writer.Boundary())


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


// Close multipart writer.
if err := writer.Close(); err != nil {
	log.Fatalf("Error closing multipart writer: %v", err)
}

	// Initialize HTTP Request and headers.
	urls := map[string]string{}

	urls["prod"] = "https://westeurope.dev-snowsoftware.io"
	urls["dev"] = "http://localhost:8008"

	uploadURL := urls[env] + "/api/filemanager/tenants/" + parsedClaims["tenant"].(string) + "/upload"

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
}