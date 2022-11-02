package main

import (
	"flag"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-cmd/cmd"
	"log"
	"net/url"
)

const (
	updateCmd    = "update"
	uploadCmd    = "upload"
	uploadBusCmd = "uploadbus"
)

var (
	snowcli = "C:\\SNOW\\SOURCES\\databackend-go-file-manager-upload\\multipartrelated\\snowcli.exe"
)

func RunCMD(name string, args ...string) (err error, stdout, stderr []string) {
	c := cmd.NewCmd(name, args...)
	s := <-c.Start()
	stdout = s.Stdout
	stderr = s.Stderr
	return
}

func main() {
	err, stdout, stderr := RunCMD(snowcli, "accounts", "get-access-token")

	if err != nil {
		panic(stderr)
	}

	bearer := stdout[0]

	command := flag.String("cmd", "", "string")
	file := flag.String("file", "", "string")
	fileId := flag.String("fileId", "", "string")
	flag.Parse()

	token, _ := jwt.Parse(bearer, func(token *jwt.Token) (interface{}, error) {
		return nil, nil
	})

	parsedClaims := token.Claims.(jwt.MapClaims)

	urlParsed, err := url.Parse(parsedClaims["iss"].(string))
	if err != nil {
		log.Fatal(err)
	}
	port := ""
	if urlParsed.Hostname() == "localhost" {
		port = ":8008"
	}
	baseUrl := urlParsed.Scheme + `://` + urlParsed.Hostname() + port

	if *command == updateCmd {
		update(baseUrl, bearer, *file, *fileId, parsedClaims)
	} else if *command == uploadCmd {
		upload(baseUrl, bearer, *file, parsedClaims)
	} else if *command == uploadBusCmd {
		uploadbus(baseUrl, bearer, *file, parsedClaims)
	} else {
		fmt.Println("Not supported command")
	}

}
