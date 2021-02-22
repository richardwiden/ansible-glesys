package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"time"
)

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const numbers = "0123456789"
const symbols = `!@#$%^&*()\/><}{[]`

type ModuleArgs struct {
	Length int
	Type   string
}

type Response struct {
	Password string `json:"password"`
	Changed  bool   `json:"changed"`
	Failed   bool   `json:"failed"`
}

func ExitJson(responseBody Response) {
	returnResponse(responseBody)
}

func FailJson(responseBody Response) {
	responseBody.Failed = true
	returnResponse(responseBody)
}

func returnResponse(responseBody Response) {
	var response []byte
	var err error
	response, err = json.Marshal(responseBody)
	if err != nil {
		response, _ = json.Marshal(Response{Password: "Invalid response object"})
	}
	fmt.Println(string(response))
	if responseBody.Failed {
		os.Exit(1)
	} else {
		os.Exit(0)
	}
}

func GeneratePassword(l int, s string) string {
	data := make([]byte, l)
	var characters string
	switch s {
	case "full":
		characters = letters + numbers + symbols
	case "alpha":
		characters = letters
	case "numeric":
		characters = numbers
	case "alphanumeric":
		characters = letters + numbers
	}

	n := len(characters)

	for i := range data {
		data[i] = characters[rand.Intn(n)]
	}

	return string(data)
}

func main() {
	var response Response

	rand.Seed(time.Now().UTC().UnixNano())

	if len(os.Args) != 2 {
		response.Password = "No argument file provided"
		FailJson(response)
	}

	argsFile := os.Args[1]

	text, err := ioutil.ReadFile(argsFile)
	if err != nil {
		response.Password = "Could not read configuration file: " + argsFile
		FailJson(response)
	}

	var moduleArgs ModuleArgs
	err = json.Unmarshal(text, &moduleArgs)
	if err != nil {
		response.Password = "Configuration file not valid JSON: " + argsFile
		FailJson(response)
	}

	var length = 20
	if moduleArgs.Length != 0 {
		length = moduleArgs.Length
	}

	var style = "full"
	if moduleArgs.Type != "" {
		style = moduleArgs.Type
	}

	password := GeneratePassword(length, style)

	response.Password = password
	ExitJson(response)
}
