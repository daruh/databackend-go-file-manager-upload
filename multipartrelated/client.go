package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"path/filepath"
	"time"
)

func main() {
	flag.Parse()

	positionalArgs := flag.Args()
	if len(positionalArgs) == 0 {
		log.Fatalf("This program requires at least 1 positional argument.")
	}

	// Metadata content.
	metadata := `{"name": "yearly-report.pdf","app": "ui","category": "finance","type": "report"}`

	// New multipart writer.
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Metadata part.
	metadataHeader := textproto.MIMEHeader{}
	metadataHeader.Set("Content-Type", "application/json")
	metadataHeader.Set("Content-Disposition", "form-data; name=\"metadata\"")
	metadataHeader.Set("Content-ID", "metadata")
	part, err := writer.CreatePart(metadataHeader)
	if err != nil {
		log.Fatalf("Error writing metadata headers: %v", err)
	}
	part.Write([]byte(metadata))

	// Media Files.
	for _, mediaFilename := range positionalArgs {
		mediaData, errRead := ioutil.ReadFile(mediaFilename)
		if errRead != nil {
			log.Fatalf("Error reading media file: %v", errRead)
		}

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

	// Close multipart writer.
	if err := writer.Close(); err != nil {
		log.Fatalf("Error closing multipart writer: %v", err)
	}

	// Request Content-Type with boundary parameter.
	contentType := fmt.Sprintf("multipart/related; boundary=%s", writer.Boundary())

	// Initialize HTTP Request and headers.
	uploadURL := "http://localhost:8008/api/inventory/files"
	//uploadURL := "http://localhost:8089/upload"
	r, err := http.NewRequest(http.MethodPost, uploadURL, bytes.NewReader(body.Bytes()))
	if err != nil {
		log.Fatalf("Error initializing a request: %v", err)
	}
	bearer := "eyJhbGciOiJSUzI1NiIsImtpZCI6IkNENjg3MjAyMzZGQkNCRjBFQ0JCNDY2RkI5OUEyNjYzNUZBMUQ2REYiLCJ4NXQiOiJ6V2h5QWpiN3lfRHN1MFp2dVpvbVkxLWgxdDgiLCJ0eXAiOiJhdCtqd3QifQ.eyJhdXRoX3RpbWUiOjE2NTY5MzQ5MzgsInN1YiI6ImFlODYyNWE0LTRiZDAtNDc4Zi05ZTdiLTA4ZGEzMTg1ODdiNSIsInRlbmFudCI6ImNkYWZlOGRkLTlkZDItNDY0Ny0zNDM1LTA4ZGEzMTg1ODc5MCIsInJlZ2lvbiI6Im5vcnRoZXVyb3BlIiwiYWNjIjoiMDk2YjgxZTUtOGRlMS00NDI2LTI0ZDYtMDhkYTMxODU4N2I4IiwiZW1wbF90eXBlIjoicm5kIiwicm9sZSI6WyJzYWxlcy5tYW5hZ2VyIiwiaWQuaWRwLmFkbWluIl0sIm9pX3Byc3QiOiIxNWY2NDg5YS0yNzQ1LTQ1YmUtOTAzNy1mYTk3MmJkZDNiMDEiLCJjbGllbnRfaWQiOiIxNWY2NDg5YS0yNzQ1LTQ1YmUtOTAzNy1mYTk3MmJkZDNiMDEiLCJ0b2tfaWQiOiI4N2M1Mjg5Ni1hMTIwLTQzYzUtYmJhZS04Zjk3M2EzYjBjMmYiLCJhdWQiOlsiYXBpOi8vc25vd3NvZnR3YXJlLmlvL2FwaSIsImFwaTovL3Nub3dzb2Z0d2FyZS5pby9pZHAvYXBpIl0sInNjb3BlIjoibGljZW5zaW5nLmxpY2Vuc2UuY3J1ZCBpZC50ZW5hbnRzLmNydWQgaWQudXNlcnMuY3J1ZCBpZC5hcHBzLmNydWQgaWQuY2xhaW1zLmNydWQgaWQub3JnYW5pemF0aW9uYWx1bml0cy5jcnVkIGlkLnN5cy50cnNjIGlkLnN5cy5jcnVkIGlkLmxvZ2lucHJvdmlkZXJzLmNydWQgaWQucm9sZXMuY3J1ZCBvcGVuaWQiLCJleHAiOjE2NTY5Mzg1MzgsImlzcyI6Imh0dHA6Ly9sb2NhbGhvc3Q6ODExOC9pZHAiLCJpYXQiOjE2NTY5MzQ5Mzh9.Vk31_My7hW8UxjNU-Yf3CAg21VMeTouCym3CTbGVNED8giqbbLcqMLm2_UfiZphcZFuDrdlxj2n5TqT0JpByheOI4cDqMNIUshU1adSL52vEMrxxMK61v4owFXJVP6gl5ZxmksQXwWOft4Qi4v7talfmrfHyjYKZAnD3UbFnfnrNrqP0vPmDlq-T7guAuiv5507eai0lOiR4UiPN2LQiy-szZB5ttlWSL-uPYI2wB2RaimZ8zFjDc3hLd0Ie3xsXF6-TrxyJ9Kr3RfqEk3p8QKeJuSyH5clnW-QKR-nlfLvGMFoUSgm1jW0jDfBvCCuuqvZ0wJQdcEEcOIlzA9r8_Q"
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
	if rsp.StatusCode != http.StatusOK {
		log.Printf("Request failed with response code: %d", rsp.StatusCode)
	} else {
		log.Print("Request was a success")
	}
}
