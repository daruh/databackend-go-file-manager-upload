package main

import (
	"encoding/base64"
	"io"
	"io/ioutil"
	"mime"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/upload", func(w http.ResponseWriter, r *http.Request) {
		contentType, params, parseErr := mime.ParseMediaType(r.Header.Get("Content-Type"))
		if parseErr != nil || !strings.HasPrefix(contentType, "multipart/") {
			http.Error(w, "expecting a multipart message", http.StatusBadRequest)
			return
		}

		multipartReader := multipart.NewReader(r.Body, params["boundary"])
		defer r.Body.Close()

		for {
			part, err := multipartReader.NextPart()
			if err == io.EOF {
				break
			}
			if err != nil {
				http.Error(w, "unexpected error when retrieving a part of the message", http.StatusInternalServerError)
				return
			}
			defer part.Close()

			fileBytes, err := ioutil.ReadAll(part)
			if err != nil {
				http.Error(w, "failed to read content of the part", http.StatusInternalServerError)
				return
			}

			switch part.FormName() {
			case "metadata":
				meta := string(fileBytes)
				if meta == "" {
					panic(meta)
				}
			case "file":
				var dataDecoded []byte
				//part.Read(fileData)
				base64.StdEncoding.
					log.Printf("filesize = %d", len(fileBytes))
				f, _ := os.Create(part.Header.Get("Content-Filename"))
				f.Write(fileBytes)
				f.Close()
			}
		}
	})

	http.ListenAndServe(":8080", mux)
}
