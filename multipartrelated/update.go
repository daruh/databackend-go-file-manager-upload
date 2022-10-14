package main

import (
	"bytes"
	"fmt"
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

func update(url, bearer, mediaFilename, fileId string, parsedClaims jwt.MapClaims) {

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

	// Request Content-Type with boundary parameter.
	contentType := fmt.Sprintf("multipart/form-data; boundary=%s", writer.Boundary())

	filename := filepath.Base(mediaFilename)

	mediaHeader := textproto.MIMEHeader{}
	mediaHeader.Set("Content-Disposition", fmt.Sprintf("form-data; name=\"file\"; filename=\"%v\"", filename))
	mediaHeader.Set("Content-Length", fmt.Sprint(size))
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

	uploadURL := url + "/api/filemanager/tenants/" + parsedClaims["tenant"].(string) + "/users/" + parsedClaims["sub"].(string) + "/files/" + fileId + "/upload"

	r, err := http.NewRequest(http.MethodPut, uploadURL, bytes.NewReader(body.Bytes()))
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
