package main

import (
	"bytes"
	"encoding/json"
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

func upload(env, bearer, mediaFilename string) {

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
