package main

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
)

func main() {
	body, err := ioutil.ReadFile("C:\\SNOW_FILES\\test_pdf.pdf")
	if err != nil {
		log.Fatalf("unable to read file: %v", err)
	}
	str := base64.StdEncoding.EncodeToString(body)
	fmt.Println(str)
}
